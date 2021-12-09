package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	observerpb "github.com/cilium/cilium/api/v1/observer"
	peerpb "github.com/cilium/cilium/api/v1/peer"
	v1 "github.com/cilium/cilium/pkg/hubble/api/v1"
	"github.com/cilium/cilium/pkg/hubble/relay/defaults"
	"github.com/rueian/gke-hubble-export/observer"
	"github.com/rueian/gke-hubble-export/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	target := env("HUBBLE_TARGET", defaults.HubbleTarget)
	replace := env("GKE_HUBBLE_EXPORT_SOCK", "/var/run/cilium/gke-hubble-export.sock")
	expAddr := env("GKE_HUBBLE_EXPORT_ADDR", "127.0.0.1")
	expPort := env("GKE_HUBBLE_EXPORT_PORT", "42444")

	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	healthSvc := health.NewServer()
	peersvc := &peer.Service{Client: peerpb.NewPeerClient(conn), Port: expPort}
	obsersvc := &observer.Service{Client: observerpb.NewObserverClient(conn)}
	healthSvc.SetServingStatus(v1.ObserverServiceName, healthpb.HealthCheckResponse_SERVING)

	os.Remove(replace)
	sock, err := net.ResolveUnixAddr("unix", replace)
	if err != nil {
		panic(err)
	}
	local, err := net.ListenUnix("unix", sock)
	if err != nil {
		panic(err)
	}
	defer local.Close()
	ln, err := net.Listen("tcp", expAddr+":"+expPort)
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	localsvc := grpc.NewServer()
	peerpb.RegisterPeerServer(localsvc, peersvc)
	observerpb.RegisterObserverServer(localsvc, obsersvc)
	healthpb.RegisterHealthServer(localsvc, healthSvc)
	go localsvc.Serve(local)

	exportsvc := grpc.NewServer()
	observerpb.RegisterObserverServer(exportsvc, obsersvc)
	healthpb.RegisterHealthServer(exportsvc, healthSvc)
	go exportsvc.Serve(ln)

	fmt.Printf("Started at %s and %s:%s\n", sock, expAddr, expPort)

	for client := observerpb.NewObserverClient(conn); true; {
		resp, err := client.ServerStatus(context.Background(), &observerpb.ServerStatusRequest{})
		if err == nil {
			fmt.Printf("ServerStatus: %v\n", resp.String())
			break
		}
		fmt.Printf("Fail to get ServerStatus: %v\n", err)
		time.Sleep(time.Second * 5)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c
	localsvc.GracefulStop()
	exportsvc.GracefulStop()
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
