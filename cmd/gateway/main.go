package main

import (
	// stdlib
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	// external
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/oklog/run"
	"google.golang.org/grpc"

	// project
	pb "github.com/dio/hello/api/hello/v1"
	"github.com/dio/hello/internal/hello"
	"github.com/dio/hello/internal/network"
)

func main() {
	var err error
	var bindIP string
	{
		if bindIP, err = network.HostIP(); err != nil {
			os.Exit(-1)
		}
	}

	listener, _ := net.Listen("tcp", bindIP+":0")
	var g run.Group
	{
		var serverOptions = []grpc.ServerOption{}
		grpcServer := grpc.NewServer(serverOptions...)

		pb.RegisterHelloServiceServer(grpcServer, hello.NewService())

		g.Add(func() error {
			return grpcServer.Serve(listener)
		}, func(error) {
			listener.Close()
		})
	}

	{
		ctx := context.Background()
		mux := runtime.NewServeMux()
		clientOptions := []grpc.DialOption{grpc.WithInsecure()}

		var httpServer *http.Server
		g.Add(func() error {
			if err = pb.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, listener.Addr().String(), clientOptions); err != nil {
				return err
			}

			httpServer = &http.Server{
				Addr:    ":8000",
				Handler: mux,
			}

			return httpServer.ListenAndServe()
		}, func(error) {
			if httpServer != nil {
				httpServer.Shutdown(ctx)
			}
		})
	}

	{
		// set-up our signal handler
		var (
			cancelInterrupt = make(chan struct{})
			c               = make(chan os.Signal, 2)
		)
		defer close(c)

		g.Add(func() error {
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}

	g.Run()
}
