/*
 *
 * Copyright 2019 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package edsbalancer contains EDS balancer implementation.
package edsbalancer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/minio/minio/pkg/grpc/balancer"
	"github.com/minio/minio/pkg/grpc/balancer/roundrobin"
	"github.com/minio/minio/pkg/grpc/connectivity"
	"github.com/minio/minio/pkg/grpc/internal/buffer"
	"github.com/minio/minio/pkg/grpc/internal/grpclog"
	"github.com/minio/minio/pkg/grpc/resolver"
	"github.com/minio/minio/pkg/grpc/serviceconfig"
	"github.com/minio/minio/pkg/grpc/xds/internal/balancer/lrs"
	xdsclient "github.com/minio/minio/pkg/grpc/xds/internal/client"
)

const (
	defaultTimeout = 10 * time.Second
	edsName        = "eds_experimental"
)

var (
	newEDSBalancer = func(cc balancer.ClientConn, enqueueState func(priorityType, balancer.State), loadStore lrs.Store, logger *grpclog.PrefixLogger) edsBalancerImplInterface {
		return newEDSBalancerImpl(cc, enqueueState, loadStore, logger)
	}
)

func init() {
	balancer.Register(&edsBalancerBuilder{})
}

type edsBalancerBuilder struct{}

// Build helps implement the balancer.Builder interface.
func (b *edsBalancerBuilder) Build(cc balancer.ClientConn, opts balancer.BuildOptions) balancer.Balancer {
	ctx, cancel := context.WithCancel(context.Background())
	x := &edsBalancer{
		ctx:               ctx,
		cancel:            cancel,
		cc:                cc,
		buildOpts:         opts,
		grpcUpdate:        make(chan interface{}),
		xdsClientUpdate:   make(chan interface{}),
		childPolicyUpdate: buffer.NewUnbounded(),
	}
	loadStore := lrs.NewStore()
	x.logger = grpclog.NewPrefixLogger(loggingPrefix(x))
	x.edsImpl = newEDSBalancer(x.cc, x.enqueueChildBalancerState, loadStore, x.logger)
	x.client = newXDSClientWrapper(x.handleEDSUpdate, x.loseContact, x.buildOpts, loadStore, x.logger)
	x.logger.Infof("Created")
	go x.run()
	return x
}

func (b *edsBalancerBuilder) Name() string {
	return edsName
}

func (b *edsBalancerBuilder) ParseConfig(c json.RawMessage) (serviceconfig.LoadBalancingConfig, error) {
	var cfg EDSConfig
	if err := json.Unmarshal(c, &cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal balancer config %s into EDSConfig, error: %v", string(c), err)
	}
	return &cfg, nil
}

// edsBalancerImplInterface defines the interface that edsBalancerImpl must
// implement to communicate with edsBalancer.
//
// It's implemented by the real eds balancer and a fake testing eds balancer.
//
// TODO: none of the methods in this interface needs to be exported.
type edsBalancerImplInterface interface {
	// HandleEDSResponse passes the received EDS message from traffic director to eds balancer.
	HandleEDSResponse(edsResp *xdsclient.EDSUpdate)
	// HandleChildPolicy updates the eds balancer the intra-cluster load balancing policy to use.
	HandleChildPolicy(name string, config json.RawMessage)
	// HandleSubConnStateChange handles state change for SubConn.
	HandleSubConnStateChange(sc balancer.SubConn, state connectivity.State)
	// updateState handle a balancer state update from the priority.
	updateState(priority priorityType, s balancer.State)
	// Close closes the eds balancer.
	Close()
}

var _ balancer.V2Balancer = (*edsBalancer)(nil) // Assert that we implement V2Balancer

// edsBalancer manages xdsClient and the actual EDS balancer implementation that
// does load balancing.
//
// It currently has only an edsBalancer. Later, we may add fallback.
type edsBalancer struct {
	cc        balancer.ClientConn // *xdsClientConn
	buildOpts balancer.BuildOptions
	ctx       context.Context
	cancel    context.CancelFunc

	logger *grpclog.PrefixLogger

	// edsBalancer continuously monitor the channels below, and will handle events from them in sync.
	grpcUpdate        chan interface{}
	xdsClientUpdate   chan interface{}
	childPolicyUpdate *buffer.Unbounded

	client  *xdsclientWrapper // may change when passed a different service config
	config  *EDSConfig        // may change when passed a different service config
	edsImpl edsBalancerImplInterface
}

// run gets executed in a goroutine once edsBalancer is created. It monitors updates from grpc,
// xdsClient and load balancer. It synchronizes the operations that happen inside edsBalancer. It
// exits when edsBalancer is closed.
func (x *edsBalancer) run() {
	for {
		select {
		case update := <-x.grpcUpdate:
			x.handleGRPCUpdate(update)
		case update := <-x.xdsClientUpdate:
			x.handleXDSClientUpdate(update)
		case update := <-x.childPolicyUpdate.Get():
			x.childPolicyUpdate.Load()
			u := update.(*balancerStateWithPriority)
			x.edsImpl.updateState(u.priority, u.s)
		case <-x.ctx.Done():
			if x.client != nil {
				x.client.close()
			}
			if x.edsImpl != nil {
				x.edsImpl.Close()
			}
			return
		}
	}
}

func (x *edsBalancer) handleGRPCUpdate(update interface{}) {
	switch u := update.(type) {
	case *subConnStateUpdate:
		if x.edsImpl != nil {
			x.edsImpl.HandleSubConnStateChange(u.sc, u.state.ConnectivityState)
		}
	case *balancer.ClientConnState:
		x.logger.Infof("Receive update from resolver, balancer config: %+v", u.BalancerConfig)
		cfg, _ := u.BalancerConfig.(*EDSConfig)
		if cfg == nil {
			// service config parsing failed. should never happen.
			return
		}

		x.client.handleUpdate(cfg, u.ResolverState.Attributes)

		if x.config == nil {
			x.config = cfg
			return
		}

		// We will update the edsImpl with the new child policy, if we got a
		// different one.
		if x.edsImpl != nil && !cmp.Equal(cfg.ChildPolicy, x.config.ChildPolicy) {
			if cfg.ChildPolicy != nil {
				x.edsImpl.HandleChildPolicy(cfg.ChildPolicy.Name, cfg.ChildPolicy.Config)
			} else {
				x.edsImpl.HandleChildPolicy(roundrobin.Name, nil)
			}
		}

		x.config = cfg
	default:
		// unreachable path
		panic("wrong update type")
	}
}

func (x *edsBalancer) handleXDSClientUpdate(update interface{}) {
	switch u := update.(type) {
	// TODO: this func should accept (*xdsclient.EDSUpdate, error), and process
	// the error, instead of having a separate loseContact signal.
	case *xdsclient.EDSUpdate:
		x.edsImpl.HandleEDSResponse(u)
	case *loseContact:
		// loseContact can be useful for going into fallback.
	default:
		panic("unexpected xds client update type")
	}
}

type subConnStateUpdate struct {
	sc    balancer.SubConn
	state balancer.SubConnState
}

func (x *edsBalancer) HandleSubConnStateChange(sc balancer.SubConn, state connectivity.State) {
	x.logger.Errorf("UpdateSubConnState should be called instead of HandleSubConnStateChange")
}

func (x *edsBalancer) HandleResolvedAddrs(addrs []resolver.Address, err error) {
	x.logger.Errorf("UpdateClientConnState should be called instead of HandleResolvedAddrs")
}

func (x *edsBalancer) UpdateSubConnState(sc balancer.SubConn, state balancer.SubConnState) {
	update := &subConnStateUpdate{
		sc:    sc,
		state: state,
	}
	select {
	case x.grpcUpdate <- update:
	case <-x.ctx.Done():
	}
}

func (x *edsBalancer) ResolverError(error) {
	// TODO: Need to distinguish between connection errors and resource removed
	// errors. For the former, we will need to handle it later on for fallback.
	// For the latter, handle it by stopping the watch, closing sub-balancers
	// and pickers.
}

func (x *edsBalancer) UpdateClientConnState(s balancer.ClientConnState) error {
	select {
	case x.grpcUpdate <- &s:
	case <-x.ctx.Done():
	}
	return nil
}

func (x *edsBalancer) handleEDSUpdate(resp *xdsclient.EDSUpdate) error {
	// TODO: this function should take (resp, error), and send them together on
	// the channel. There doesn't need to be a separate `loseContact` function.
	select {
	case x.xdsClientUpdate <- resp:
	case <-x.ctx.Done():
	}

	return nil
}

type loseContact struct {
}

// TODO: delete loseContact when handleEDSUpdate takes (resp, error).
func (x *edsBalancer) loseContact() {
	select {
	case x.xdsClientUpdate <- &loseContact{}:
	case <-x.ctx.Done():
	}
}

type balancerStateWithPriority struct {
	priority priorityType
	s        balancer.State
}

func (x *edsBalancer) enqueueChildBalancerState(p priorityType, s balancer.State) {
	x.childPolicyUpdate.Put(&balancerStateWithPriority{
		priority: p,
		s:        s,
	})
}

func (x *edsBalancer) Close() {
	x.cancel()
	x.logger.Infof("Shutdown")
}
