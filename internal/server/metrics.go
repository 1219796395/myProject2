package server

// NewMetricServer new a HTTP server.
import (
	"game-config/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type MetricsServer struct {
	Server *http.Server
}

// NewMatricsServer new a HTTP server.
func NewMetricsServer(bc *conf.Bootstrap, logger log.Logger) *MetricsServer {
	var opts = []http.ServerOption{}
	if bc.Server.Http.Network != "" {
		opts = append(opts, http.Network(bc.Server.Http.Network))
	}
	if bc.Server.Http.MetricAddr != "" {
		opts = append(opts, http.Address(bc.Server.Http.MetricAddr))
	}
	if bc.Server.Http.Timeout != nil {
		opts = append(opts, http.Timeout(bc.Server.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)

	// metrics
	srv.Handle("/metrics", promhttp.Handler())

	return &MetricsServer{
		Server: srv,
	}
}
