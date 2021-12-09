package observer

import (
	"context"

	observerpb "github.com/cilium/cilium/api/v1/observer"
	"github.com/rueian/gke-hubble-export/proxy"
	"google.golang.org/protobuf/proto"
)

var _ observerpb.ObserverServer = (*Service)(nil)

type Service struct {
	Client observerpb.ObserverClient
}

func (s *Service) GetFlows(request *observerpb.GetFlowsRequest, stream observerpb.Observer_GetFlowsServer) error {
	client, err := s.Client.GetFlows(stream.Context(), request)
	if err != nil {
		return err
	}
	return proxy.ServerStreaming(client, stream, func() proto.Message {
		return new(observerpb.GetFlowsResponse)
	}, func(msg proto.Message) {})
}

func (s *Service) GetAgentEvents(request *observerpb.GetAgentEventsRequest, stream observerpb.Observer_GetAgentEventsServer) error {
	client, err := s.Client.GetAgentEvents(stream.Context(), request)
	if err != nil {
		return err
	}
	return proxy.ServerStreaming(client, stream, func() proto.Message {
		return new(observerpb.GetAgentEventsResponse)
	}, func(msg proto.Message) {})
}

func (s *Service) GetDebugEvents(request *observerpb.GetDebugEventsRequest, stream observerpb.Observer_GetDebugEventsServer) error {
	client, err := s.Client.GetDebugEvents(stream.Context(), request)
	if err != nil {
		return err
	}
	return proxy.ServerStreaming(client, stream, func() proto.Message {
		return new(observerpb.GetDebugEventsResponse)
	}, func(msg proto.Message) {})
}

func (s *Service) GetNodes(ctx context.Context, request *observerpb.GetNodesRequest) (*observerpb.GetNodesResponse, error) {
	return s.Client.GetNodes(ctx, request)
}

func (s *Service) ServerStatus(ctx context.Context, request *observerpb.ServerStatusRequest) (*observerpb.ServerStatusResponse, error) {
	return s.Client.ServerStatus(ctx, request)
}
