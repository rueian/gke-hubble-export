package peer

import (
	"fmt"
	"net"
	"strings"

	peerpb "github.com/cilium/cilium/api/v1/peer"
	"github.com/rueian/gke-hubble-export/proxy"
	"google.golang.org/protobuf/proto"
)

var _ peerpb.PeerServer = (*Service)(nil)

type Service struct {
	Client peerpb.PeerClient
	Port   string
}

func (s *Service) Notify(req *peerpb.NotifyRequest, stream peerpb.Peer_NotifyServer) error {
	client, err := s.Client.Notify(stream.Context(), req)
	if err != nil {
		return err
	}
	return proxy.ServerStreaming(client, stream, func() proto.Message {
		return new(peerpb.ChangeNotification)
	}, func(msg proto.Message) {
		v := msg.(*peerpb.ChangeNotification)
		v.Tls = nil
		fmt.Printf("Receive From Notify: %v\n", v.String())
		if host, _, err := net.SplitHostPort(v.Address); err == nil {
			v.Address = net.JoinHostPort(host, s.Port)
		} else if strings.Contains(err.Error(), "missing port") {
			v.Address = v.Address + ":" + s.Port
		}
	})
}
