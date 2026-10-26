package main

import (
	"context"
	"flag"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grpc-demo/calculator/config"
	"grpc-demo/calculator/service"
	"grpc-demo/common/healthing"
	"grpc-demo/pbdata/calculatorpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

var serviceName = "Calculator"

func main() {
	configFile := flag.String("c", "calculator/config/config.yml", "configure file path")
	flag.Parse()

	if err := config.Init(*configFile); err != nil {
		slog.Error("init configuration", "error", err)
		os.Exit(1)
	}

	listener, err := net.Listen(config.Conf.RPCServer.Protocol, config.Conf.RPCServer.Address)
	if err != nil {
		slog.Error("listen gRPC", "address", config.Conf.RPCServer.Address, "error", err)
		os.Exit(1)
	}

	grpcServer := grpc.NewServer()
	calculatorpb.RegisterCalculatorServer(grpcServer, &service.Server{Name: serviceName})

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	reflection.Register(grpcServer)

	mux := http.NewServeMux()
	mux.HandleFunc("/status.html", healthing.HttpStatus)
	statusServer := &http.Server{
		Addr:    config.Conf.HealthServer.Address,
		Handler: mux,
	}

	serverErr := make(chan error, 2)
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			serverErr <- err
		}
	}()
	go func() {
		if err := statusServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()
	slog.Info("gRPC server started", "address", config.Conf.RPCServer.Address)
	slog.Info("status server started", "address", config.Conf.HealthServer.Address)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stopChan)

	select {
	case err := <-serverErr:
		slog.Error("server failed", "error", err)
	case sig := <-stopChan:
		slog.Info("shutdown requested", "signal", sig.String())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = statusServer.Shutdown(shutdownCtx)
	grpcServer.GracefulStop()
}
