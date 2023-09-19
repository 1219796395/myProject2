package server

import (
	pb3 "game-config/api/networkconfig"
	pb2 "game-config/api/operationlog/remoteconfiglog"
	pb1 "game-config/api/projectconfig/envmanage"
	pb "game-config/api/remoteconfig"
	"game-config/internal/conf"
	"game-config/internal/service"

	"github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// NewGRPCServer new a gRPC server.
func NewGRPCServer(bc *conf.Bootstrap, s *service.RemoteConfigService, s1 *service.EnvManageService, s2 *service.RemoteConfigLogService,
	s3 *service.NetworkConfigService, log log.Logger) *grpc.Server {
	var c = bc.Server
	var opts = []grpc.ServerOption{
		grpc.Middleware(
			recovery.Recovery(),
			metrics.Server(
				metrics.WithSeconds(prometheus.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prometheus.NewCounter(_metricRequests)),
			),
			tracing.Server(),
			ratelimit.Server(),
			validate.Validator(),
		),
	}
	if c.Grpc.Network != "" {
		opts = append(opts, grpc.Network(c.Grpc.Network))
	}
	if c.Grpc.Addr != "" {
		opts = append(opts, grpc.Address(c.Grpc.Addr))
	}
	if c.Grpc.Timeout != nil {
		opts = append(opts, grpc.Timeout(c.Grpc.Timeout.AsDuration()))
	}
	srv := grpc.NewServer(opts...)

	pb.RegisterRemoteConfigServer(srv, s)
	pb1.RegisterEnvManageServer(srv, s1)
	pb2.RegisterRemoteConfigLogServer(srv, s2)
	pb3.RegisterNetworkConfigServer(srv, s3)
	return srv
}
