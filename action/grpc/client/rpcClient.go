package client

import (
	"hcc/cello/action/grpc/pb/rpcflute"
	"hcc/cello/action/grpc/pb/rpcharp"
)

// RPCClient : Struct type of gRPC clients
type RPCClient struct {
	flute rpcflute.FluteClient
	harp  rpcharp.HarpClient
}

// RC : Exported variable pointed to RPCClient
var RC = &RPCClient{}

// Init : Initialize clients of gRPC
func Init() error {

	return nil
}

// End : Close connections of gRPC clients
func End() {

}
