package server

import (
	"context"
	"fmt"
	"github.com/mklfarha/radarcdmx/backend/rcapi/core"
	pb "github.com/mklfarha/radarcdmx/backend/rcapi/idl/gen"
	"net"

	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedRcapiServer
	core *core.Implementation
}

type Params struct {
	fx.In
	Logger    *zap.Logger `optional:"true"`
	Lifecycle fx.Lifecycle
	Core      *core.Implementation
	Config    config.Provider
}

func NewServer(params Params) pb.RcapiServer {
	return &server{
		core: params.Core,
	}
}

func New(params Params) *grpc.Server {
	log := params.Logger
	if log == nil {
		log = zap.NewNop()
	}
	s := grpc.NewServer()
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// GRPC port from config
			grpcPort := params.Config.Get("ports.grpc").String()

			// proto server
			lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
			if err != nil {
				log.Error("failed to listen: %v", zap.Error(err))
				return err
			}

			// register grpc servers
			grpc_health_v1.RegisterHealthServer(s, health.NewServer())
			pb.RegisterRcapiServer(s, NewServer(params))
			reflection.Register(s)

			log.Info("GRPC Server listening at %v", zap.Any("addr", lis.Addr()))

			go s.Serve(lis)

			return nil

		},
		OnStop: func(ctx context.Context) error {
			s.Stop()
			return nil
		},
	})

	return s
}
