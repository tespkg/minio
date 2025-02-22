/*
 *
 * Copyright 2017 gRPC authors.
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

package grpc

import (
	"errors"

	"github.com/minio/minio/pkg/grpc/balancer"
	"github.com/minio/minio/pkg/grpc/codes"
	"github.com/minio/minio/pkg/grpc/connectivity"
	"github.com/minio/minio/pkg/grpc/grpclog"
	"github.com/minio/minio/pkg/grpc/resolver"
	"github.com/minio/minio/pkg/grpc/status"
)

// PickFirstBalancerName is the name of the pick_first balancer.
const PickFirstBalancerName = "pick_first"

func newPickfirstBuilder() balancer.Builder {
	return &pickfirstBuilder{}
}

type pickfirstBuilder struct{}

func (*pickfirstBuilder) Build(cc balancer.ClientConn, opt balancer.BuildOptions) balancer.Balancer {
	return &pickfirstBalancer{cc: cc}
}

func (*pickfirstBuilder) Name() string {
	return PickFirstBalancerName
}

type pickfirstBalancer struct {
	state connectivity.State
	cc    balancer.ClientConn
	sc    balancer.SubConn
}

var _ balancer.V2Balancer = &pickfirstBalancer{} // Assert we implement v2

func (b *pickfirstBalancer) HandleResolvedAddrs(addrs []resolver.Address, err error) {
	if err != nil {
		b.ResolverError(err)
		return
	}
	b.UpdateClientConnState(balancer.ClientConnState{ResolverState: resolver.State{Addresses: addrs}}) // Ignore error
}

func (b *pickfirstBalancer) HandleSubConnStateChange(sc balancer.SubConn, s connectivity.State) {
	b.UpdateSubConnState(sc, balancer.SubConnState{ConnectivityState: s})
}

func (b *pickfirstBalancer) ResolverError(err error) {
	switch b.state {
	case connectivity.TransientFailure, connectivity.Idle, connectivity.Connecting:
		// Set a failing picker if we don't have a good picker.
		b.cc.UpdateState(balancer.State{ConnectivityState: connectivity.TransientFailure,
			Picker: &picker{err: status.Errorf(codes.Unavailable, "name resolver error: %v", err)}},
		)
	}
	if grpclog.V(2) {
		grpclog.Infof("pickfirstBalancer: ResolverError called with error %v", err)
	}
}

func (b *pickfirstBalancer) UpdateClientConnState(cs balancer.ClientConnState) error {
	if len(cs.ResolverState.Addresses) == 0 {
		b.ResolverError(errors.New("produced zero addresses"))
		return balancer.ErrBadResolverState
	}
	if b.sc == nil {
		var err error
		b.sc, err = b.cc.NewSubConn(cs.ResolverState.Addresses, balancer.NewSubConnOptions{})
		if err != nil {
			if grpclog.V(2) {
				grpclog.Errorf("pickfirstBalancer: failed to NewSubConn: %v", err)
			}
			b.state = connectivity.TransientFailure
			b.cc.UpdateState(balancer.State{ConnectivityState: connectivity.TransientFailure,
				Picker: &picker{err: status.Errorf(codes.Unavailable, "error creating connection: %v", err)}},
			)
			return balancer.ErrBadResolverState
		}
		b.state = connectivity.Idle
		b.cc.UpdateState(balancer.State{ConnectivityState: connectivity.Idle, Picker: &picker{result: balancer.PickResult{SubConn: b.sc}}})
		b.sc.Connect()
	} else {
		b.sc.UpdateAddresses(cs.ResolverState.Addresses)
		b.sc.Connect()
	}
	return nil
}

func (b *pickfirstBalancer) UpdateSubConnState(sc balancer.SubConn, s balancer.SubConnState) {
	if grpclog.V(2) {
		grpclog.Infof("pickfirstBalancer: HandleSubConnStateChange: %p, %v", sc, s)
	}
	if b.sc != sc {
		if grpclog.V(2) {
			grpclog.Infof("pickfirstBalancer: ignored state change because sc is not recognized")
		}
		return
	}
	b.state = s.ConnectivityState
	if s.ConnectivityState == connectivity.Shutdown {
		b.sc = nil
		return
	}

	switch s.ConnectivityState {
	case connectivity.Ready, connectivity.Idle:
		b.cc.UpdateState(balancer.State{ConnectivityState: s.ConnectivityState, Picker: &picker{result: balancer.PickResult{SubConn: sc}}})
	case connectivity.Connecting:
		b.cc.UpdateState(balancer.State{ConnectivityState: s.ConnectivityState, Picker: &picker{err: balancer.ErrNoSubConnAvailable}})
	case connectivity.TransientFailure:
		err := balancer.ErrTransientFailure
		// TODO: this can be unconditional after the V1 API is removed, as
		// SubConnState will always contain a connection error.
		if s.ConnectionError != nil {
			err = balancer.TransientFailureError(s.ConnectionError)
		}
		b.cc.UpdateState(balancer.State{
			ConnectivityState: s.ConnectivityState,
			Picker:            &picker{err: err},
		})
	}
}

func (b *pickfirstBalancer) Close() {
}

type picker struct {
	result balancer.PickResult
	err    error
}

func (p *picker) Pick(info balancer.PickInfo) (balancer.PickResult, error) {
	return p.result, p.err
}

func init() {
	balancer.Register(newPickfirstBuilder())
}
