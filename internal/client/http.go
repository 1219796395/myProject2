package client

import (
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var Client *resty.Client

func init() {
	Client = resty.New()
	// set max conns
	Client.GetClient().Transport = &http.Transport{
		// max idle connection nums is unlimited(0)
		// max idle connection num is 100
		MaxIdleConns:        0,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		// default settings
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		ForceAttemptHTTP2:     true,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	// set timeout and header
	Client.SetTimeout(5 * time.Second)
	Client.SetHeader("Accept", "application/json")

	Client.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		ctx, span := otel.Tracer("resty_mannul").Start(r.Context(), "http.req")
		span.SetAttributes(
			attribute.String("http.req.url", r.URL),
			attribute.String("http.req.method", r.Method),
			attribute.String("http.req.query_param", r.QueryParam.Encode()),
		)
		r.SetContext(ctx)

		return nil
	})

	Client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		span := trace.SpanFromContext(r.Request.Context())
		defer span.End()

		span.SetAttributes(
			attribute.Int("http.resp.status_code", r.StatusCode()),
		)

		return nil
	})
}
