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

package rpctypes

import (
	"testing"

	"github.com/minio/minio/pkg/grpc/codes"
	"github.com/minio/minio/pkg/grpc/status"
)

func TestConvert(t *testing.T) {
	e1 := status.New(codes.InvalidArgument, "etcdserver: key is not provided").Err()
	e2 := ErrGRPCEmptyKey
	e3 := ErrEmptyKey

	if e1.Error() != e2.Error() {
		t.Fatalf("expected %q == %q", e1.Error(), e2.Error())
	}
	ev1, _ := status.FromError(e1)
	if ev1.Code() != e3.(EtcdError).Code() {
		t.Fatalf("expected them to be equal, got %v / %v", ev1.Code(), e3.(EtcdError).Code())
	}

	if e1.Error() == e3.Error() {
		t.Fatalf("expected %q != %q", e1.Error(), e3.Error())
	}
	ev2, _ := status.FromError(e2)
	if ev2.Code() != e3.(EtcdError).Code() {
		t.Fatalf("expected them to be equal, got %v / %v", ev2.Code(), e3.(EtcdError).Code())
	}
}
