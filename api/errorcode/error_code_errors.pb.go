// Code generated by protoc-gen-go-errors. DO NOT EDIT.

package errorcode

import (
	fmt "fmt"
	errors "github.com/go-kratos/kratos/v2/errors"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
const _ = errors.SupportPackageIsVersion1

// 预留值无意义
func IsReserve(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_RESERVE.String() && e.Code == 200
}

// 预留值无意义
func ErrorReserve(format string, args ...interface{}) *errors.Error {
	return errors.New(200, Errors_RESERVE.String(), fmt.Sprintf(format, args...))
}

// common bizcode, [1,99]
func IsBadRequest(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_BAD_REQUEST.String() && e.Code == 400
}

// common bizcode, [1,99]
func ErrorBadRequest(format string, args ...interface{}) *errors.Error {
	return errors.New(400, Errors_BAD_REQUEST.String(), fmt.Sprintf(format, args...))
}

// 服务器内部错误
func IsInternalServerError(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_INTERNAL_SERVER_ERROR.String() && e.Code == 500
}

// 服务器内部错误
func ErrorInternalServerError(format string, args ...interface{}) *errors.Error {
	return errors.New(500, Errors_INTERNAL_SERVER_ERROR.String(), fmt.Sprintf(format, args...))
}

// remote_config bizcode, [1001, 1999]
func IsRemoteConfigAlreadyExist(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_REMOTE_CONFIG_ALREADY_EXIST.String() && e.Code == 400
}

// remote_config bizcode, [1001, 1999]
func ErrorRemoteConfigAlreadyExist(format string, args ...interface{}) *errors.Error {
	return errors.New(400, Errors_REMOTE_CONFIG_ALREADY_EXIST.String(), fmt.Sprintf(format, args...))
}

// 远程配置不存在
func IsRemoteConfigNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_REMOTE_CONFIG_NOT_FOUND.String() && e.Code == 404
}

// 远程配置不存在
func ErrorRemoteConfigNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, Errors_REMOTE_CONFIG_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// request过期
func IsRemoteConfigExpireRequest(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_REMOTE_CONFIG_EXPIRE_REQUEST.String() && e.Code == 400
}

// request过期
func ErrorRemoteConfigExpireRequest(format string, args ...interface{}) *errors.Error {
	return errors.New(400, Errors_REMOTE_CONFIG_EXPIRE_REQUEST.String(), fmt.Sprintf(format, args...))
}

// env_manage bizcode, [2001, 2999]
func IsEnvManageNameAlreadyExist(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_ENV_MANAGE_NAME_ALREADY_EXIST.String() && e.Code == 400
}

// env_manage bizcode, [2001, 2999]
func ErrorEnvManageNameAlreadyExist(format string, args ...interface{}) *errors.Error {
	return errors.New(400, Errors_ENV_MANAGE_NAME_ALREADY_EXIST.String(), fmt.Sprintf(format, args...))
}

// 环境不存在
func IsEnvManageNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_ENV_MANAGE_NOT_FOUND.String() && e.Code == 404
}

// 环境不存在
func ErrorEnvManageNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, Errors_ENV_MANAGE_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// request过期
func IsEnvManageExpireRequest(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_ENV_MANAGE_EXPIRE_REQUEST.String() && e.Code == 400
}

// request过期
func ErrorEnvManageExpireRequest(format string, args ...interface{}) *errors.Error {
	return errors.New(400, Errors_ENV_MANAGE_EXPIRE_REQUEST.String(), fmt.Sprintf(format, args...))
}

// 禁止修改拥有配置的环境
func IsEnvManageNotModifyEnvHasConf(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_ENV_MANAGE_NOT_MODIFY_ENV_HAS_CONF.String() && e.Code == 403
}

// 禁止修改拥有配置的环境
func ErrorEnvManageNotModifyEnvHasConf(format string, args ...interface{}) *errors.Error {
	return errors.New(403, Errors_ENV_MANAGE_NOT_MODIFY_ENV_HAS_CONF.String(), fmt.Sprintf(format, args...))
}

// network_config bizcode, [3001, 3999]
func IsNetworkConfigAlreadyExist(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_NETWORK_CONFIG_ALREADY_EXIST.String() && e.Code == 400
}

// network_config bizcode, [3001, 3999]
func ErrorNetworkConfigAlreadyExist(format string, args ...interface{}) *errors.Error {
	return errors.New(400, Errors_NETWORK_CONFIG_ALREADY_EXIST.String(), fmt.Sprintf(format, args...))
}

// network_config不存在
func IsNetworkConfigNotFound(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_NETWORK_CONFIG_NOT_FOUND.String() && e.Code == 404
}

// network_config不存在
func ErrorNetworkConfigNotFound(format string, args ...interface{}) *errors.Error {
	return errors.New(404, Errors_NETWORK_CONFIG_NOT_FOUND.String(), fmt.Sprintf(format, args...))
}

// request过期
func IsNetworkConfigExpireRequest(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_NETWORK_CONFIG_EXPIRE_REQUEST.String() && e.Code == 400
}

// request过期
func ErrorNetworkConfigExpireRequest(format string, args ...interface{}) *errors.Error {
	return errors.New(400, Errors_NETWORK_CONFIG_EXPIRE_REQUEST.String(), fmt.Sprintf(format, args...))
}

// auth bizcode, [4001, 4999]
func IsInvalidToken(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_INVALID_TOKEN.String() && e.Code == 403
}

// auth bizcode, [4001, 4999]
func ErrorInvalidToken(format string, args ...interface{}) *errors.Error {
	return errors.New(403, Errors_INVALID_TOKEN.String(), fmt.Sprintf(format, args...))
}

// 用户缺少权限
func IsNotAuthorized(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_NOT_AUTHORIZED.String() && e.Code == 401
}

// 用户缺少权限
func ErrorNotAuthorized(format string, args ...interface{}) *errors.Error {
	return errors.New(401, Errors_NOT_AUTHORIZED.String(), fmt.Sprintf(format, args...))
}

// 用户在该应用下无权限
func IsAppidMismatch(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_APPID_MISMATCH.String() && e.Code == 401
}

// 用户在该应用下无权限
func ErrorAppidMismatch(format string, args ...interface{}) *errors.Error {
	return errors.New(401, Errors_APPID_MISMATCH.String(), fmt.Sprintf(format, args...))
}

// 用户不存在
func IsNoSuchUser(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_NO_SUCH_USER.String() && e.Code == 404
}

// 用户不存在
func ErrorNoSuchUser(format string, args ...interface{}) *errors.Error {
	return errors.New(404, Errors_NO_SUCH_USER.String(), fmt.Sprintf(format, args...))
}

// 角色不存在
func IsNoSuchRole(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_NO_SUCH_ROLE.String() && e.Code == 404
}

// 角色不存在
func ErrorNoSuchRole(format string, args ...interface{}) *errors.Error {
	return errors.New(404, Errors_NO_SUCH_ROLE.String(), fmt.Sprintf(format, args...))
}

// 该权限点不存在
func IsNoSuchResource(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_NO_SUCH_RESOURCE.String() && e.Code == 404
}

// 该权限点不存在
func ErrorNoSuchResource(format string, args ...interface{}) *errors.Error {
	return errors.New(404, Errors_NO_SUCH_RESOURCE.String(), fmt.Sprintf(format, args...))
}

// 用户已存在,无需反复创建
func IsUserAlreadyExist(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_USER_ALREADY_EXIST.String() && e.Code == 200
}

// 用户已存在,无需反复创建
func ErrorUserAlreadyExist(format string, args ...interface{}) *errors.Error {
	return errors.New(200, Errors_USER_ALREADY_EXIST.String(), fmt.Sprintf(format, args...))
}

// 角色已存在，无需反复创建
func IsRoleAlreadyExist(err error) bool {
	if err == nil {
		return false
	}
	e := errors.FromError(err)
	return e.Reason == Errors_ROLE_ALREADY_EXIST.String() && e.Code == 200
}

// 角色已存在，无需反复创建
func ErrorRoleAlreadyExist(format string, args ...interface{}) *errors.Error {
	return errors.New(200, Errors_ROLE_ALREADY_EXIST.String(), fmt.Sprintf(format, args...))
}