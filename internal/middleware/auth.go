package middleware

import (
	"context"

	authPb "github.com/1219796395/myProject2/api/auth"
	commonPb "github.com/1219796395/myProject2/api/common"
	errorPb "github.com/1219796395/myProject2/api/errorcode"
	"github.com/1219796395/myProject2/internal/biz"
	"github.com/1219796395/myProject2/internal/biz/bo"
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel/trace"
)

const (
	XTraceID            = "x-trace-id"
	LoginToken          = "x-hg-login-token"
	AdminUserContextKey = "adminUser"
)

type AdminCommonReq interface {
	GetCommon() *commonPb.Common
}

type AuthMiddleware struct {
	log       *log.Helper
	adminUser *biz.AdminUserLogic
	secretes  map[uint32]string
}

func NewAuthMiddleware(
	logger log.Logger,
	adminUser *biz.AdminUserLogic,
	bc *conf.Bootstrap,
) *AuthMiddleware {
	return &AuthMiddleware{
		log:       log.NewHelper(logger),
		adminUser: adminUser,
		secretes:  bc.Biz.Admin.Secretes,
	}
}

func (m *AuthMiddleware) Auth() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
			tr, ok := transport.FromServerContext(ctx)
			if !ok {
				return nil, errorPb.ErrorInternalServerError("transport context error")
			}
			// set trace_id in reply header
			traceID := trace.SpanContextFromContext(ctx).TraceID()
			if traceID.IsValid() {
				tr.ReplyHeader().Set(XTraceID, traceID.String())
			}
			// common
			tempReq, _ := req.(AdminCommonReq)
			common := tempReq.GetCommon()
			// 校验并获取admin_user_id
			var token string
			if ts, ok := transport.FromServerContext(ctx); ok && ts.Kind() == transport.KindHTTP {
				token = ts.RequestHeader().Get(LoginToken)
			} else {
				m.log.WithContext(ctx).Errorf("[AuthMiddleware.Auth] can not get token from this request")
				return nil, errorPb.ErrorNotAuthorized("invalid token")
			}

			if len(token) == 0 {
				m.log.WithContext(ctx).Errorf("[AuthMiddleware.Auth] token is empty from this request")
				return nil, errorPb.ErrorBadRequest("empty token")
			}

			// for secretes for each app id
			if secret, ok := m.secretes[common.AppId]; ok && secret == token {
				ctx = context.WithValue(ctx, AdminUserContextKey, authPb.AdminUser{
					SkipAuth: true,
				})
				// let pass
				return handler(ctx, req)
			}

			adminUser, err := m.adminUser.VerifyUser(ctx, token)
			if err != nil {
				m.log.Error("Auth middelware ", err.Error())
				return nil, errorPb.ErrorInternalServerError("transport context error")
			}
			if adminUser == nil {
				// 用户正常状态数据不存在，校验失败
				m.log.WithContext(ctx).Errorf("[AuthMiddleware.Auth] not authorized, user not exist")
				return nil, errorPb.ErrorNotAuthorized("invalid token")
			}

			// store user info into context
			adminPb := ConvertToPBAdminUser(adminUser)
			ctx = context.WithValue(ctx, AdminUserContextKey, adminPb)

			// validate the authentication of this user
			can, err := m.adminUser.Can(ctx, adminUser, tr.Operation(), common.AppId)
			if err != nil {
				m.log.Error("Auth middelware ", err.Error())
				return nil, errorPb.ErrorInternalServerError("transport context error")
			}
			if !can {
				// 用户无权限，校验失败
				m.log.WithContext(ctx).Errorf("[AuthMiddleware.Auth] not authorized, %s, operation %s", adminUser.Name, tr.Operation())
				return nil, errorPb.ErrorInvalidToken("invalid token")
			}

			defer func() {
				// Do something on exiting
			}()

			return handler(ctx, req)
		}
	}
}

func ConvertToPBAdminUser(adminUser *bo.AdminUser) *authPb.AdminUser {
	return &authPb.AdminUser{
		Id:          adminUser.Id,
		HgId:        adminUser.HgId,
		Email:       adminUser.Email,
		NickName:    adminUser.Nickname,
		Name:        adminUser.Name,
		CreatedAtTs: uint32(adminUser.CreatedAt.Unix()),
		UpdatedAtTs: uint32(adminUser.UpdatedAt.Unix()),
	}
}

func GetAdminUser(ctx context.Context) *authPb.AdminUser {
	res, ok := ctx.Value(AdminUserContextKey).(*authPb.AdminUser)
	if !ok {
		return &authPb.AdminUser{SkipAuth: true}
	}
	return res
}
