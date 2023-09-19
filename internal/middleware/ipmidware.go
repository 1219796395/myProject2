package middleware

import (
	"context"

	"github.com/1219796395/myProject2/internal/conf"

	"net"
	"strings"

	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

var clientIpKeys = []string{"X-Envoy-External-Address", "X-Real-IP"}

type IpMiddleware struct {
	log *log.Helper
}

func NewIpMiddleware(logger log.Logger, bc *conf.Bootstrap) *IpMiddleware {
	return &IpMiddleware{
		log: log.NewHelper(logger),
	}
}

// Validator is a validator middleware
func (m *IpMiddleware) GetIp() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			trans, ok := transport.FromServerContext(ctx)
			if !ok {
				m.log.WithContext(ctx).Error("[TimeoutSettingMiddleware] getting context from server failed")
				return handler(ctx, req)
			}

			// set user ip value into context.value
			ctx = context.WithValue(ctx, "clientIP", m.getClientIP(ctx, clientIpKeys, trans))
			return handler(ctx, req)
		}
	}
}

func (m *IpMiddleware) getClientIP(ctx context.Context, keys []string, transport transport.Transporter) (clientIP string) {
	for _, key := range keys {
		clientIP = transport.RequestHeader().Get(key)
		if clientIP != "" {
			return clientIP
		}
	}

	clientIP = transport.RequestHeader().Get("X-Forwarded-For")
	if clientIP != "" {
		items := strings.Split(clientIP, ",")
		if items[0] != "" {
			return items[0]
		}
	}

	clientIP, err := m.remoteIp(ctx, transport)
	if err != nil {
		return ""
	}
	return clientIP
}

func (m *IpMiddleware) remoteIp(ctx context.Context, transport transport.Transporter) (string, error) {
	httpTransport, ok := transport.(http.Transporter)
	if !ok {
		m.log.WithContext(ctx).Error("[IpMiddleware.remoteIp] transport is not a http transport")
		return "", errors.New("not a http trasnport")
	}
	ip, _, err := net.SplitHostPort(strings.TrimSpace(httpTransport.Request().RemoteAddr))
	if err != nil {
		m.log.WithContext(ctx).Error("[IpMiddleware.remoteIp] failed to get remote address: ", err)
		return "", err
	}
	return ip, nil
}
