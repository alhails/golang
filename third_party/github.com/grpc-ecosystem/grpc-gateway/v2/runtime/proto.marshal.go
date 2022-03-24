// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// ProtoMarshaller []byte -> proto|interface{}
type ProtoMarshaller struct {
	runtime.ProtoMarshaller
}

// ContentType always returns "application/x-protobuf".
func (*ProtoMarshaller) ContentType(_ interface{}) string {
	return binding.MIMEPROTOBUF
}
