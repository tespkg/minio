// Copyright 2016 The etcd Authors
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

package integration

import (
	"io/ioutil"

	"github.com/minio/minio/pkg/coreos/etcd/clientv3"

	"github.com/coreos/pkg/capnslog"
	"github.com/minio/minio/pkg/grpc/grpclog"
)

const defaultLogLevel = capnslog.CRITICAL

func init() {
	capnslog.SetGlobalLogLevel(defaultLogLevel)
	clientv3.SetLogger(grpclog.NewLoggerV2(ioutil.Discard, ioutil.Discard, ioutil.Discard))
}
