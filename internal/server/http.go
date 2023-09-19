package server

import (
	"context"
	pb4 "game-config/api/auth"
	pb3 "game-config/api/networkconfig"
	pb2 "game-config/api/operationlog/remoteconfiglog"
	pb1 "game-config/api/projectconfig/envmanage"
	pb "game-config/api/remoteconfig"
	"game-config/internal/conf"
	"game-config/internal/middleware"
	"game-config/internal/service"
	srcHttp "net/http"
	"strings"

	"github.com/go-kratos/kratos/contrib/metrics/prometheus/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/metrics"
	"github.com/go-kratos/kratos/v2/middleware/ratelimit"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/tracing"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-kratos/swagger-api/openapiv2"
	"github.com/gorilla/handlers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpServer struct {
	Server *http.Server
}

// NewHTTPServer new a HTTP server.
func NewHTTPServer(
	bc *conf.Bootstrap,
	s *service.RemoteConfigService,
	s1 *service.EnvManageService,
	s2 *service.RemoteConfigLogService,
	s3 *service.NetworkConfigService,
	s4 *service.AdminAuthService,
	validator *middleware.ValidateWithLogMiddleware,
	ipMidware *middleware.IpMiddleware,
	authMidware *middleware.AuthMiddleware,
	logger log.Logger) *HttpServer {
	var c = bc.Server
	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			ipMidware.GetIp(),
			metrics.Server(
				metrics.WithSeconds(prometheus.NewHistogram(_metricSeconds)),
				metrics.WithRequests(prometheus.NewCounter(_metricRequests)),
			),
			tracing.Server(),
			ratelimit.Server(),
			validator.Validator(),
			// admin auth
			selector.Server(recovery.Recovery(), authMidware.Auth()).Match(
				func(ctx context.Context, operation string) bool {
					// basic role and login, no need to auth
					if operation == "/api.auth.Auth/Me" {
						return false
					}
					if operation == "/api.auth.Auth/GenerateTokenByCode" {
						return false
					}

					// toc interfaces, no need to auth
					if operation == "/api.remoteconfig.RemoteConfig/GetRemoteConfig" ||
						operation == "/api.remoteconfig.RemoteConfig/GetRemoteConfigV1" ||
						operation == "/api.remoteconfig.RemoteConfig/CreateRemoteConfigV1" ||
						operation == "/api.remoteconfig.RemoteConfig/UpdateRemoteConfigV1" ||
						operation == "/api.remoteconfig.RemoteConfig/PublishRemoteConfigV1" ||
						operation == "/api.networkconfig.NetworkConfig/GetNetworkConfig" {
						return false
					}

					// TODO, modules finished, others still under develipment
					// TODO, remove when finish
					if strings.HasPrefix(operation, "/api.remoteconfig.RemoteConfig") {
						return true
					}
					if strings.HasPrefix(operation, "/api.projectconfig.envmanage.EnvManage") {
						return true
					}

					// auth all tob interfaces
					return true
				},
			).Build(),
		),
		http.Filter(handlers.CORS(
			handlers.AllowedHeaders([]string{"Content-Type", "x-hg-login-token", "x-hg-internal"}),
			handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}),
			handlers.AllowedOrigins(bc.Server.Http.Cors.AllowedOrigins),
		)),
	}

	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	opts = append(opts, http.ResponseEncoder(customResponseEncoder))
	srv := http.NewServer(opts...)

	// health check
	srv.HandleFunc("/checkHealth", func(w srcHttp.ResponseWriter, r *srcHttp.Request) {
		w.WriteHeader(200)
		w.Write([]byte("success"))
	})

	// swagger
	if c.Env == "dev" || c.Env == "qa" {
		srv.HandlePrefix("/q/", openapiv2.NewHandler())
	}

	// metrics
	// TODO: new port
	srv.Handle("/metrics", promhttp.Handler())

	pb.RegisterRemoteConfigHTTPServer(srv, s)
	pb1.RegisterEnvManageHTTPServer(srv, s1)
	pb2.RegisterRemoteConfigLogHTTPServer(srv, s2)
	pb3.RegisterNetworkConfigHTTPServer(srv, s3)
	pb4.RegisterAuthHTTPServer(srv, s4)

	return &HttpServer{
		Server: srv,
	}
}

func customResponseEncoder(w http.ResponseWriter, r *http.Request, v interface{}) error {
	// logic for download
	if res, ok := v.(*pb.GetRemoteConfigRsp); ok {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(res.ConfigData))
		if err != nil {
			return err
		}
		return nil
	}

	// logic for normal interfaces
	codec, _ := http.CodecForRequest(r, "Accept")
	data, err := codec.Marshal(v)
	if err != nil {
		return err
	}

	// write response body
	w.Header().Set("Content-Type", strings.Join([]string{"application", codec.Name()}, "/"))
	w.Write(data)
	return nil
}
