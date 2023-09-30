package internalgrpc

import (
	"fmt"
	"net"

	pb "github.com/grahovsky/go-hw/hw12_13_14_15_calendar/api/apppb"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/api"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/grahovsky/go-hw/hw12_13_14_15_calendar/internal/logger"
	"google.golang.org/grpc"
)

func NewServer(app *app.App, addr string) *Server {
	srv := grpc.NewServer(grpc.UnaryInterceptor(serverUnaryInterceptor))
	pb.RegisterAppServer(srv, api.NewAPI(app))

	return &Server{
		addr: addr,
		srv:  srv,
	}
}

type Server struct {
	addr string
	srv  *grpc.Server
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", s.addr, err)
	}
	logger.Info(fmt.Sprintf("gRPC server is starting on %s", s.addr))
	if err = s.srv.Serve(listener); err != nil {
		return fmt.Errorf("gRPC serve: %w", err)
	}
	return nil
}

func (s *Server) Stop() {
	s.srv.GracefulStop()
}
