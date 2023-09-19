package middleware

import (
	"context"
	"game-config/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
)

type ValidateWithLogMiddleware struct {
	log                      *log.Helper
	onlyLogValidationFailure bool
}

func NewValidateWithLogMiddleware(logger log.Logger, bc *conf.Bootstrap) *ValidateWithLogMiddleware {
	return &ValidateWithLogMiddleware{
		log:                      log.NewHelper(logger),
		onlyLogValidationFailure: bc.Server.Log.OnlyLogValidationFailure,
	}
}

type validator interface {
	Validate() error
}

// Validator is a validator middleware.
func (m *ValidateWithLogMiddleware) Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			if v, ok := req.(validator); ok {
				if err := v.Validate(); err != nil {
					// print log
					reqJson, _ := json.MarshalToString(req)
					m.log.WithContext(ctx).Warnf("[ValidateWithLogMiddleware] bad request %s, request content is %s", err.Error(), reqJson)
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}

			// if every request need to be logged
			if !m.onlyLogValidationFailure {
				reqJson, _ := json.MarshalToString(req)
				m.log.WithContext(ctx).Infof("[ValidateWithLogMiddleware] request content is %s", reqJson)
			}
			return handler(ctx, req)
		}
	}
}
