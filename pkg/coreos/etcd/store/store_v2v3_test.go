// Copyright 2017 The etcd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build v2v3

package store_test

import (
	"io/ioutil"
	"testing"

	"github.com/minio/minio/pkg/coreos/etcd/clientv3"
	"github.com/minio/minio/pkg/coreos/etcd/etcdserver/api/v2v3"
	"github.com/minio/minio/pkg/coreos/etcd/integration"
	"github.com/minio/minio/pkg/coreos/etcd/store"

	"github.com/coreos/pkg/capnslog"
	"github.com/minio/minio/pkg/grpc/grpclog"
)

func init() {
	capnslog.SetGlobalLogLevel(capnslog.CRITICAL)
	clientv3.SetLogger(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, ioutil.Discard))
}

type v2v3TestStore struct {
	store.Store
	clus *integration.ClusterV3
	t    *testing.T
}

func (s *v2v3TestStore) Close() { s.clus.Terminate(s.t) }

func newTestStore(t *testing.T, ns ...string) StoreCloser {
	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	return &v2v3TestStore{
		v2v3.NewStore(clus.Client(0), "/v2/"),
		clus,
		t,
	}
}
