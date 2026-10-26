package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"grpc-demo/common/healthing"
	"grpc-demo/generalgateway/config"
	"grpc-demo/generalgateway/controller"
	"grpc-demo/generalgateway/grpc_client"

	"github.com/gin-gonic/gin"
)

func main() {
	configFile := flag.String("c", "generalgateway/config/config.yml", "configure file path")
	flag.Parse()

	if err := config.Init(*configFile); err != nil {
		slog.Error("init configuration", "error", err)
		os.Exit(1)
	}
	if err := grpc_client.Init(config.Conf.GRpcClientConfig); err != nil {
		slog.Error("init gRPC client", "error", err)
		os.Exit(1)
	}
	defer grpc_client.Close()

	gin.SetMode(config.Conf.Gin.Mode)
	router := gin.Default()
	router.GET("/status.html", healthing.Status)
	v1 := router.Group("/v1")
	{
		v1.POST("/hello", controller.Hello)
		v1.POST("/calculate", controller.Calculate)
	}

	server := &http.Server{
		Addr:    config.Conf.Server.Address,
		Handler: router,
	}
	serverErr := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
		}
	}()
	slog.Info("HTTP gateway started", "address", config.Conf.Server.Address)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stopChan)

	select {
	case err := <-serverErr:
		slog.Error("gateway failed", "error", err)
	case sig := <-stopChan:
		slog.Info("shutdown requested", "signal", sig.String())
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}
