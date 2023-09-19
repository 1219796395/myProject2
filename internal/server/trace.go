package server

import (
	"game-config/internal/conf"
	"os"

	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var (
	// trace agent endpoint
	MY_HOST_IP = os.Getenv("MY_HOST_IP")
)

func init() {
	bc, err := conf.GetConf()
	if err != nil {
		return
	}

	if bc.Server.Trace == nil || !bc.Server.Trace.OnOff {
		return
	}

	endpoint := bc.Server.Trace.Endpoint
	if len(MY_HOST_IP) > 0 {
		endpoint = MY_HOST_IP
	}

	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(endpoint)))
	if err != nil {
		return
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(bc.Server.Trace.Ratio))),
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String(bc.Server.Name),
			attribute.String("env", bc.Server.Env),
		)),
	)
	otel.SetTracerProvider(tp)

	log.Infof("[trace] init ok")
}
