/*
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
 */

package balancergroup

import (
	"fmt"
	"testing"
	"time"

	orcapb "github.com/cncf/udpa/go/udpa/data/orca/v1"
	"github.com/google/go-cmp/cmp"
	"github.com/minio/minio/pkg/grpc/balancer"
	"github.com/minio/minio/pkg/grpc/balancer/roundrobin"
	"github.com/minio/minio/pkg/grpc/connectivity"
	"github.com/minio/minio/pkg/grpc/resolver"
	"github.com/minio/minio/pkg/grpc/xds/internal"
	"github.com/minio/minio/pkg/grpc/xds/internal/testutils"
)

var (
	rrBuilder        = balancer.Get(roundrobin.Name)
	testBalancerIDs  = []internal.Locality{{Region: "b1"}, {Region: "b2"}, {Region: "b3"}}
	testBackendAddrs []resolver.Address
)

const testBackendAddrsCount = 12

func init() {
	for i := 0; i < testBackendAddrsCount; i++ {
		testBackendAddrs = append(testBackendAddrs, resolver.Address{Addr: fmt.Sprintf("%d.%d.%d.%d:%d", i, i, i, i, i)})
	}

	// Disable caching for all tests. It will be re-enabled in caching specific
	// tests.
	DefaultSubBalancerCloseTimeout = time.Millisecond
}

func subConnFromPicker(p balancer.V2Picker) func() balancer.SubConn {
	return func() balancer.SubConn {
		scst, _ := p.Pick(balancer.PickInfo{})
		return scst.SubConn
	}
}

// 1 balancer, 1 backend -> 2 backends -> 1 backend.
func (s) TestBalancerGroup_OneRR_AddRemoveBackend(t *testing.T) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)
	bg.Start()

	// Add one balancer to group.
	bg.Add(testBalancerIDs[0], 1, rrBuilder)
	// Send one resolved address.
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:1])

	// Send subconn state change.
	sc1 := <-cc.NewSubConnCh
	bg.HandleSubConnStateChange(sc1, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc1, connectivity.Ready)

	// Test pick with one backend.
	p1 := <-cc.NewPickerCh
	for i := 0; i < 5; i++ {
		gotSCSt, _ := p1.Pick(balancer.PickInfo{})
		if !cmp.Equal(gotSCSt.SubConn, sc1, cmp.AllowUnexported(testutils.TestSubConn{})) {
			t.Fatalf("picker.Pick, got %v, want SubConn=%v", gotSCSt, sc1)
		}
	}

	// Send two addresses.
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:2])
	// Expect one new subconn, send state update.
	sc2 := <-cc.NewSubConnCh
	bg.HandleSubConnStateChange(sc2, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc2, connectivity.Ready)

	// Test roundrobin pick.
	p2 := <-cc.NewPickerCh
	want := []balancer.SubConn{sc1, sc2}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p2)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	// Remove the first address.
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[1:2])
	scToRemove := <-cc.RemoveSubConnCh
	if !cmp.Equal(scToRemove, sc1, cmp.AllowUnexported(testutils.TestSubConn{})) {
		t.Fatalf("RemoveSubConn, want %v, got %v", sc1, scToRemove)
	}
	bg.HandleSubConnStateChange(scToRemove, connectivity.Shutdown)

	// Test pick with only the second subconn.
	p3 := <-cc.NewPickerCh
	for i := 0; i < 5; i++ {
		gotSC, _ := p3.Pick(balancer.PickInfo{})
		if !cmp.Equal(gotSC.SubConn, sc2, cmp.AllowUnexported(testutils.TestSubConn{})) {
			t.Fatalf("picker.Pick, got %v, want SubConn=%v", gotSC, sc2)
		}
	}
}

// 2 balancers, each with 1 backend.
func (s) TestBalancerGroup_TwoRR_OneBackend(t *testing.T) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)
	bg.Start()

	// Add two balancers to group and send one resolved address to both
	// balancers.
	bg.Add(testBalancerIDs[0], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:1])
	sc1 := <-cc.NewSubConnCh

	bg.Add(testBalancerIDs[1], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[0:1])
	sc2 := <-cc.NewSubConnCh

	// Send state changes for both subconns.
	bg.HandleSubConnStateChange(sc1, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc1, connectivity.Ready)
	bg.HandleSubConnStateChange(sc2, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc2, connectivity.Ready)

	// Test roundrobin on the last picker.
	p1 := <-cc.NewPickerCh
	want := []balancer.SubConn{sc1, sc2}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p1)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}
}

// 2 balancers, each with more than 1 backends.
func (s) TestBalancerGroup_TwoRR_MoreBackends(t *testing.T) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)
	bg.Start()

	// Add two balancers to group and send one resolved address to both
	// balancers.
	bg.Add(testBalancerIDs[0], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:2])
	sc1 := <-cc.NewSubConnCh
	sc2 := <-cc.NewSubConnCh

	bg.Add(testBalancerIDs[1], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[2:4])
	sc3 := <-cc.NewSubConnCh
	sc4 := <-cc.NewSubConnCh

	// Send state changes for both subconns.
	bg.HandleSubConnStateChange(sc1, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc1, connectivity.Ready)
	bg.HandleSubConnStateChange(sc2, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc2, connectivity.Ready)
	bg.HandleSubConnStateChange(sc3, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc3, connectivity.Ready)
	bg.HandleSubConnStateChange(sc4, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc4, connectivity.Ready)

	// Test roundrobin on the last picker.
	p1 := <-cc.NewPickerCh
	want := []balancer.SubConn{sc1, sc2, sc3, sc4}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p1)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	// Turn sc2's connection down, should be RR between balancers.
	bg.HandleSubConnStateChange(sc2, connectivity.TransientFailure)
	p2 := <-cc.NewPickerCh
	// Expect two sc1's in the result, because balancer1 will be picked twice,
	// but there's only one sc in it.
	want = []balancer.SubConn{sc1, sc1, sc3, sc4}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p2)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	// Remove sc3's addresses.
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[3:4])
	scToRemove := <-cc.RemoveSubConnCh
	if !cmp.Equal(scToRemove, sc3, cmp.AllowUnexported(testutils.TestSubConn{})) {
		t.Fatalf("RemoveSubConn, want %v, got %v", sc3, scToRemove)
	}
	bg.HandleSubConnStateChange(scToRemove, connectivity.Shutdown)
	p3 := <-cc.NewPickerCh
	want = []balancer.SubConn{sc1, sc4}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p3)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	// Turn sc1's connection down.
	bg.HandleSubConnStateChange(sc1, connectivity.TransientFailure)
	p4 := <-cc.NewPickerCh
	want = []balancer.SubConn{sc4}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p4)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	// Turn last connection to connecting.
	bg.HandleSubConnStateChange(sc4, connectivity.Connecting)
	p5 := <-cc.NewPickerCh
	for i := 0; i < 5; i++ {
		if _, err := p5.Pick(balancer.PickInfo{}); err != balancer.ErrNoSubConnAvailable {
			t.Fatalf("want pick error %v, got %v", balancer.ErrNoSubConnAvailable, err)
		}
	}

	// Turn all connections down.
	bg.HandleSubConnStateChange(sc4, connectivity.TransientFailure)
	p6 := <-cc.NewPickerCh
	for i := 0; i < 5; i++ {
		if _, err := p6.Pick(balancer.PickInfo{}); err != balancer.ErrTransientFailure {
			t.Fatalf("want pick error %v, got %v", balancer.ErrTransientFailure, err)
		}
	}
}

// 2 balancers with different weights.
func (s) TestBalancerGroup_TwoRR_DifferentWeight_MoreBackends(t *testing.T) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)
	bg.Start()

	// Add two balancers to group and send two resolved addresses to both
	// balancers.
	bg.Add(testBalancerIDs[0], 2, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:2])
	sc1 := <-cc.NewSubConnCh
	sc2 := <-cc.NewSubConnCh

	bg.Add(testBalancerIDs[1], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[2:4])
	sc3 := <-cc.NewSubConnCh
	sc4 := <-cc.NewSubConnCh

	// Send state changes for both subconns.
	bg.HandleSubConnStateChange(sc1, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc1, connectivity.Ready)
	bg.HandleSubConnStateChange(sc2, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc2, connectivity.Ready)
	bg.HandleSubConnStateChange(sc3, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc3, connectivity.Ready)
	bg.HandleSubConnStateChange(sc4, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc4, connectivity.Ready)

	// Test roundrobin on the last picker.
	p1 := <-cc.NewPickerCh
	want := []balancer.SubConn{sc1, sc1, sc2, sc2, sc3, sc4}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p1)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}
}

// totally 3 balancers, add/remove balancer.
func (s) TestBalancerGroup_ThreeRR_RemoveBalancer(t *testing.T) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)
	bg.Start()

	// Add three balancers to group and send one resolved address to both
	// balancers.
	bg.Add(testBalancerIDs[0], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:1])
	sc1 := <-cc.NewSubConnCh

	bg.Add(testBalancerIDs[1], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[1:2])
	sc2 := <-cc.NewSubConnCh

	bg.Add(testBalancerIDs[2], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[2], testBackendAddrs[1:2])
	sc3 := <-cc.NewSubConnCh

	// Send state changes for both subconns.
	bg.HandleSubConnStateChange(sc1, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc1, connectivity.Ready)
	bg.HandleSubConnStateChange(sc2, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc2, connectivity.Ready)
	bg.HandleSubConnStateChange(sc3, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc3, connectivity.Ready)

	p1 := <-cc.NewPickerCh
	want := []balancer.SubConn{sc1, sc2, sc3}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p1)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	// Remove the second balancer, while the others two are ready.
	bg.Remove(testBalancerIDs[1])
	scToRemove := <-cc.RemoveSubConnCh
	if !cmp.Equal(scToRemove, sc2, cmp.AllowUnexported(testutils.TestSubConn{})) {
		t.Fatalf("RemoveSubConn, want %v, got %v", sc2, scToRemove)
	}
	p2 := <-cc.NewPickerCh
	want = []balancer.SubConn{sc1, sc3}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p2)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	// move balancer 3 into transient failure.
	bg.HandleSubConnStateChange(sc3, connectivity.TransientFailure)
	// Remove the first balancer, while the third is transient failure.
	bg.Remove(testBalancerIDs[0])
	scToRemove = <-cc.RemoveSubConnCh
	if !cmp.Equal(scToRemove, sc1, cmp.AllowUnexported(testutils.TestSubConn{})) {
		t.Fatalf("RemoveSubConn, want %v, got %v", sc1, scToRemove)
	}
	p3 := <-cc.NewPickerCh
	for i := 0; i < 5; i++ {
		if _, err := p3.Pick(balancer.PickInfo{}); err != balancer.ErrTransientFailure {
			t.Fatalf("want pick error %v, got %v", balancer.ErrTransientFailure, err)
		}
	}
}

// 2 balancers, change balancer weight.
func (s) TestBalancerGroup_TwoRR_ChangeWeight_MoreBackends(t *testing.T) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)
	bg.Start()

	// Add two balancers to group and send two resolved addresses to both
	// balancers.
	bg.Add(testBalancerIDs[0], 2, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:2])
	sc1 := <-cc.NewSubConnCh
	sc2 := <-cc.NewSubConnCh

	bg.Add(testBalancerIDs[1], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[2:4])
	sc3 := <-cc.NewSubConnCh
	sc4 := <-cc.NewSubConnCh

	// Send state changes for both subconns.
	bg.HandleSubConnStateChange(sc1, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc1, connectivity.Ready)
	bg.HandleSubConnStateChange(sc2, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc2, connectivity.Ready)
	bg.HandleSubConnStateChange(sc3, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc3, connectivity.Ready)
	bg.HandleSubConnStateChange(sc4, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc4, connectivity.Ready)

	// Test roundrobin on the last picker.
	p1 := <-cc.NewPickerCh
	want := []balancer.SubConn{sc1, sc1, sc2, sc2, sc3, sc4}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p1)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	bg.ChangeWeight(testBalancerIDs[0], 3)

	// Test roundrobin with new weight.
	p2 := <-cc.NewPickerCh
	want = []balancer.SubConn{sc1, sc1, sc1, sc2, sc2, sc2, sc3, sc4}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p2)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}
}

func (s) TestBalancerGroup_LoadReport(t *testing.T) {
	testLoadStore := testutils.NewTestLoadStore()

	cc := testutils.NewTestClientConn(t)
	bg := New(cc, testLoadStore, nil)
	bg.Start()

	backendToBalancerID := make(map[balancer.SubConn]internal.Locality)

	// Add two balancers to group and send two resolved addresses to both
	// balancers.
	bg.Add(testBalancerIDs[0], 2, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:2])
	sc1 := <-cc.NewSubConnCh
	sc2 := <-cc.NewSubConnCh
	backendToBalancerID[sc1] = testBalancerIDs[0]
	backendToBalancerID[sc2] = testBalancerIDs[0]

	bg.Add(testBalancerIDs[1], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[2:4])
	sc3 := <-cc.NewSubConnCh
	sc4 := <-cc.NewSubConnCh
	backendToBalancerID[sc3] = testBalancerIDs[1]
	backendToBalancerID[sc4] = testBalancerIDs[1]

	// Send state changes for both subconns.
	bg.HandleSubConnStateChange(sc1, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc1, connectivity.Ready)
	bg.HandleSubConnStateChange(sc2, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc2, connectivity.Ready)
	bg.HandleSubConnStateChange(sc3, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc3, connectivity.Ready)
	bg.HandleSubConnStateChange(sc4, connectivity.Connecting)
	bg.HandleSubConnStateChange(sc4, connectivity.Ready)

	// Test roundrobin on the last picker.
	p1 := <-cc.NewPickerCh
	var (
		wantStart []internal.Locality
		wantEnd   []internal.Locality
		wantCost  []testutils.TestServerLoad
	)
	for i := 0; i < 10; i++ {
		scst, _ := p1.Pick(balancer.PickInfo{})
		locality := backendToBalancerID[scst.SubConn]
		wantStart = append(wantStart, locality)
		if scst.Done != nil && scst.SubConn != sc1 {
			scst.Done(balancer.DoneInfo{
				ServerLoad: &orcapb.OrcaLoadReport{
					CpuUtilization: 10,
					MemUtilization: 5,
					RequestCost:    map[string]float64{"pic": 3.14},
					Utilization:    map[string]float64{"piu": 3.14},
				},
			})
			wantEnd = append(wantEnd, locality)
			wantCost = append(wantCost,
				testutils.TestServerLoad{Name: serverLoadCPUName, D: 10},
				testutils.TestServerLoad{Name: serverLoadMemoryName, D: 5},
				testutils.TestServerLoad{Name: "pic", D: 3.14},
				testutils.TestServerLoad{Name: "piu", D: 3.14})
		}
	}

	if !cmp.Equal(testLoadStore.CallsStarted, wantStart) {
		t.Fatalf("want started: %v, got: %v", testLoadStore.CallsStarted, wantStart)
	}
	if !cmp.Equal(testLoadStore.CallsEnded, wantEnd) {
		t.Fatalf("want ended: %v, got: %v", testLoadStore.CallsEnded, wantEnd)
	}
	if !cmp.Equal(testLoadStore.CallsCost, wantCost, cmp.AllowUnexported(testutils.TestServerLoad{})) {
		t.Fatalf("want cost: %v, got: %v", testLoadStore.CallsCost, wantCost)
	}
}

// Create a new balancer group, add balancer and backends, but not start.
// - b1, weight 2, backends [0,1]
// - b2, weight 1, backends [2,3]
// Start the balancer group and check behavior.
//
// Close the balancer group, call add/remove/change weight/change address.
// - b2, weight 3, backends [0,3]
// - b3, weight 1, backends [1,2]
// Start the balancer group again and check for behavior.
func (s) TestBalancerGroup_start_close(t *testing.T) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)

	// Add two balancers to group and send two resolved addresses to both
	// balancers.
	bg.Add(testBalancerIDs[0], 2, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:2])
	bg.Add(testBalancerIDs[1], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[2:4])

	bg.Start()

	m1 := make(map[resolver.Address]balancer.SubConn)
	for i := 0; i < 4; i++ {
		addrs := <-cc.NewSubConnAddrsCh
		sc := <-cc.NewSubConnCh
		m1[addrs[0]] = sc
		bg.HandleSubConnStateChange(sc, connectivity.Connecting)
		bg.HandleSubConnStateChange(sc, connectivity.Ready)
	}

	// Test roundrobin on the last picker.
	p1 := <-cc.NewPickerCh
	want := []balancer.SubConn{
		m1[testBackendAddrs[0]], m1[testBackendAddrs[0]],
		m1[testBackendAddrs[1]], m1[testBackendAddrs[1]],
		m1[testBackendAddrs[2]], m1[testBackendAddrs[3]],
	}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p1)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	bg.Close()
	for i := 0; i < 4; i++ {
		bg.HandleSubConnStateChange(<-cc.RemoveSubConnCh, connectivity.Shutdown)
	}

	// Add b3, weight 1, backends [1,2].
	bg.Add(testBalancerIDs[2], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[2], testBackendAddrs[1:3])

	// Remove b1.
	bg.Remove(testBalancerIDs[0])

	// Update b2 to weight 3, backends [0,3].
	bg.ChangeWeight(testBalancerIDs[1], 3)
	bg.HandleResolvedAddrs(testBalancerIDs[1], append([]resolver.Address(nil), testBackendAddrs[0], testBackendAddrs[3]))

	bg.Start()

	m2 := make(map[resolver.Address]balancer.SubConn)
	for i := 0; i < 4; i++ {
		addrs := <-cc.NewSubConnAddrsCh
		sc := <-cc.NewSubConnCh
		m2[addrs[0]] = sc
		bg.HandleSubConnStateChange(sc, connectivity.Connecting)
		bg.HandleSubConnStateChange(sc, connectivity.Ready)
	}

	// Test roundrobin on the last picker.
	p2 := <-cc.NewPickerCh
	want = []balancer.SubConn{
		m2[testBackendAddrs[0]], m2[testBackendAddrs[0]], m2[testBackendAddrs[0]],
		m2[testBackendAddrs[3]], m2[testBackendAddrs[3]], m2[testBackendAddrs[3]],
		m2[testBackendAddrs[1]], m2[testBackendAddrs[2]],
	}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p2)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}
}

// Test that balancer group start() doesn't deadlock if the balancer calls back
// into balancer group inline when it gets an update.
//
// The potential deadlock can happen if we
//   - hold a lock and send updates to balancer (e.g. update resolved addresses)
//   - the balancer calls back (NewSubConn or update picker) in line
//
// The callback will try to hold hte same lock again, which will cause a
// deadlock.
//
// This test starts the balancer group with a test balancer, will updates picker
// whenever it gets an address update. It's expected that start() doesn't block
// because of deadlock.
func (s) TestBalancerGroup_start_close_deadlock(t *testing.T) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)

	bg.Add(testBalancerIDs[0], 2, &testutils.TestConstBalancerBuilder{})
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:2])
	bg.Add(testBalancerIDs[1], 1, &testutils.TestConstBalancerBuilder{})
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[2:4])

	bg.Start()
}

func replaceDefaultSubBalancerCloseTimeout(n time.Duration) func() {
	old := DefaultSubBalancerCloseTimeout
	DefaultSubBalancerCloseTimeout = n
	return func() { DefaultSubBalancerCloseTimeout = old }
}

// initBalancerGroupForCachingTest creates a balancer group, and initialize it
// to be ready for caching tests.
//
// Two rr balancers are added to bg, each with 2 ready subConns. A sub-balancer
// is removed later, so the balancer group returned has one sub-balancer in its
// own map, and one sub-balancer in cache.
func initBalancerGroupForCachingTest(t *testing.T) (*BalancerGroup, *testutils.TestClientConn, map[resolver.Address]balancer.SubConn) {
	cc := testutils.NewTestClientConn(t)
	bg := New(cc, nil, nil)

	// Add two balancers to group and send two resolved addresses to both
	// balancers.
	bg.Add(testBalancerIDs[0], 2, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[0], testBackendAddrs[0:2])
	bg.Add(testBalancerIDs[1], 1, rrBuilder)
	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[2:4])

	bg.Start()

	m1 := make(map[resolver.Address]balancer.SubConn)
	for i := 0; i < 4; i++ {
		addrs := <-cc.NewSubConnAddrsCh
		sc := <-cc.NewSubConnCh
		m1[addrs[0]] = sc
		bg.HandleSubConnStateChange(sc, connectivity.Connecting)
		bg.HandleSubConnStateChange(sc, connectivity.Ready)
	}

	// Test roundrobin on the last picker.
	p1 := <-cc.NewPickerCh
	want := []balancer.SubConn{
		m1[testBackendAddrs[0]], m1[testBackendAddrs[0]],
		m1[testBackendAddrs[1]], m1[testBackendAddrs[1]],
		m1[testBackendAddrs[2]], m1[testBackendAddrs[3]],
	}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p1)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	bg.Remove(testBalancerIDs[1])
	// Don't wait for SubConns to be removed after close, because they are only
	// removed after close timeout.
	for i := 0; i < 10; i++ {
		select {
		case <-cc.RemoveSubConnCh:
			t.Fatalf("Got request to remove subconn, want no remove subconn (because subconns were still in cache)")
		default:
		}
		time.Sleep(time.Millisecond)
	}
	// Test roundrobin on the with only sub-balancer0.
	p2 := <-cc.NewPickerCh
	want = []balancer.SubConn{
		m1[testBackendAddrs[0]], m1[testBackendAddrs[1]],
	}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p2)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	return bg, cc, m1
}

// Test that if a sub-balancer is removed, and re-added within close timeout,
// the subConns won't be re-created.
func (s) TestBalancerGroup_locality_caching(t *testing.T) {
	defer replaceDefaultSubBalancerCloseTimeout(10 * time.Second)()
	bg, cc, addrToSC := initBalancerGroupForCachingTest(t)

	// Turn down subconn for addr2, shouldn't get picker update because
	// sub-balancer1 was removed.
	bg.HandleSubConnStateChange(addrToSC[testBackendAddrs[2]], connectivity.TransientFailure)
	for i := 0; i < 10; i++ {
		select {
		case <-cc.NewPickerCh:
			t.Fatalf("Got new picker, want no new picker (because the sub-balancer was removed)")
		default:
		}
		time.Sleep(time.Millisecond)
	}

	// Sleep, but sleep less then close timeout.
	time.Sleep(time.Millisecond * 100)

	// Re-add sub-balancer-1, because subconns were in cache, no new subconns
	// should be created. But a new picker will still be generated, with subconn
	// states update to date.
	bg.Add(testBalancerIDs[1], 1, rrBuilder)

	p3 := <-cc.NewPickerCh
	want := []balancer.SubConn{
		addrToSC[testBackendAddrs[0]], addrToSC[testBackendAddrs[0]],
		addrToSC[testBackendAddrs[1]], addrToSC[testBackendAddrs[1]],
		// addr2 is down, b2 only has addr3 in READY state.
		addrToSC[testBackendAddrs[3]], addrToSC[testBackendAddrs[3]],
	}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p3)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}

	for i := 0; i < 10; i++ {
		select {
		case <-cc.NewSubConnAddrsCh:
			t.Fatalf("Got new subconn, want no new subconn (because subconns were still in cache)")
		default:
		}
		time.Sleep(time.Millisecond * 10)
	}
}

// Sub-balancers are put in cache when they are removed. If balancer group is
// closed within close timeout, all subconns should still be rmeoved
// immediately.
func (s) TestBalancerGroup_locality_caching_close_group(t *testing.T) {
	defer replaceDefaultSubBalancerCloseTimeout(10 * time.Second)()
	bg, cc, addrToSC := initBalancerGroupForCachingTest(t)

	bg.Close()
	// The balancer group is closed. The subconns should be removed immediately.
	removeTimeout := time.After(time.Millisecond * 500)
	scToRemove := map[balancer.SubConn]int{
		addrToSC[testBackendAddrs[0]]: 1,
		addrToSC[testBackendAddrs[1]]: 1,
		addrToSC[testBackendAddrs[2]]: 1,
		addrToSC[testBackendAddrs[3]]: 1,
	}
	for i := 0; i < len(scToRemove); i++ {
		select {
		case sc := <-cc.RemoveSubConnCh:
			c := scToRemove[sc]
			if c == 0 {
				t.Fatalf("Got removeSubConn for %v when there's %d remove expected", sc, c)
			}
			scToRemove[sc] = c - 1
		case <-removeTimeout:
			t.Fatalf("timeout waiting for subConns (from balancer in cache) to be removed")
		}
	}
}

// Sub-balancers in cache will be closed if not re-added within timeout, and
// subConns will be removed.
func (s) TestBalancerGroup_locality_caching_not_readd_within_timeout(t *testing.T) {
	defer replaceDefaultSubBalancerCloseTimeout(time.Second)()
	_, cc, addrToSC := initBalancerGroupForCachingTest(t)

	// The sub-balancer is not re-added withtin timeout. The subconns should be
	// removed.
	removeTimeout := time.After(DefaultSubBalancerCloseTimeout)
	scToRemove := map[balancer.SubConn]int{
		addrToSC[testBackendAddrs[2]]: 1,
		addrToSC[testBackendAddrs[3]]: 1,
	}
	for i := 0; i < len(scToRemove); i++ {
		select {
		case sc := <-cc.RemoveSubConnCh:
			c := scToRemove[sc]
			if c == 0 {
				t.Fatalf("Got removeSubConn for %v when there's %d remove expected", sc, c)
			}
			scToRemove[sc] = c - 1
		case <-removeTimeout:
			t.Fatalf("timeout waiting for subConns (from balancer in cache) to be removed")
		}
	}
}

// Wrap the rr builder, so it behaves the same, but has a different pointer.
type noopBalancerBuilderWrapper struct {
	balancer.Builder
}

// After removing a sub-balancer, re-add with same ID, but different balancer
// builder. Old subconns should be removed, and new subconns should be created.
func (s) TestBalancerGroup_locality_caching_readd_with_different_builder(t *testing.T) {
	defer replaceDefaultSubBalancerCloseTimeout(10 * time.Second)()
	bg, cc, addrToSC := initBalancerGroupForCachingTest(t)

	// Re-add sub-balancer-1, but with a different balancer builder. The
	// sub-balancer was still in cache, but cann't be reused. This should cause
	// old sub-balancer's subconns to be removed immediately, and new subconns
	// to be created.
	bg.Add(testBalancerIDs[1], 1, &noopBalancerBuilderWrapper{rrBuilder})

	// The cached sub-balancer should be closed, and the subconns should be
	// removed immediately.
	removeTimeout := time.After(time.Millisecond * 500)
	scToRemove := map[balancer.SubConn]int{
		addrToSC[testBackendAddrs[2]]: 1,
		addrToSC[testBackendAddrs[3]]: 1,
	}
	for i := 0; i < len(scToRemove); i++ {
		select {
		case sc := <-cc.RemoveSubConnCh:
			c := scToRemove[sc]
			if c == 0 {
				t.Fatalf("Got removeSubConn for %v when there's %d remove expected", sc, c)
			}
			scToRemove[sc] = c - 1
		case <-removeTimeout:
			t.Fatalf("timeout waiting for subConns (from balancer in cache) to be removed")
		}
	}

	bg.HandleResolvedAddrs(testBalancerIDs[1], testBackendAddrs[4:6])

	newSCTimeout := time.After(time.Millisecond * 500)
	scToAdd := map[resolver.Address]int{
		testBackendAddrs[4]: 1,
		testBackendAddrs[5]: 1,
	}
	for i := 0; i < len(scToAdd); i++ {
		select {
		case addr := <-cc.NewSubConnAddrsCh:
			c := scToAdd[addr[0]]
			if c == 0 {
				t.Fatalf("Got newSubConn for %v when there's %d new expected", addr, c)
			}
			scToAdd[addr[0]] = c - 1
			sc := <-cc.NewSubConnCh
			addrToSC[addr[0]] = sc
			bg.HandleSubConnStateChange(sc, connectivity.Connecting)
			bg.HandleSubConnStateChange(sc, connectivity.Ready)
		case <-newSCTimeout:
			t.Fatalf("timeout waiting for subConns (from new sub-balancer) to be newed")
		}
	}

	// Test roundrobin on the new picker.
	p3 := <-cc.NewPickerCh
	want := []balancer.SubConn{
		addrToSC[testBackendAddrs[0]], addrToSC[testBackendAddrs[0]],
		addrToSC[testBackendAddrs[1]], addrToSC[testBackendAddrs[1]],
		addrToSC[testBackendAddrs[4]], addrToSC[testBackendAddrs[5]],
	}
	if err := testutils.IsRoundRobin(want, subConnFromPicker(p3)); err != nil {
		t.Fatalf("want %v, got %v", want, err)
	}
}
