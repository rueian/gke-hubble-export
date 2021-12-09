package proxy

import (
	"io"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func ServerStreaming(client grpc.ClientStream, server grpc.ServerStream, msgFn func() proto.Message, mutate func(msg proto.Message)) error {
	defer client.CloseSend()
	for {
		msg := msgFn()
		if err := client.RecvMsg(msg); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		mutate(msg)
		if err := server.SendMsg(msg); err != nil {
			return err
		}
	}
}
