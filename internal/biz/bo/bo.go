package bo

import (
	"time"
)

const (
	AdminTokenKey          = "admin:token:%s"
	AdminTokenValidSeconds = 3600 * 24 * 5

	// 用户app id关系状态
	UserAppIdNormal = 1
	UserAppIdBanned = 2

	// 管理角色状态
	AdminRoleNormal  = 1
	AdminRoleDeleted = 2

	// 管理员身份
	IdentityOrdinary = 1 // 一般角色，可以被有修改权限的用户增删修改
	IdentitySuper    = 2 // 超管角色，除开发者外不能修改，一个超管可对应一个或全部项目

	// TODO, 定义权限点枚举值

	// 信息化sso相关
	SSO_GET_TOKEN         = "/oauth/token"
	SSO_GET_INFO_BY_TOKEN = "/api/mdm/profile/v1/me"

	// 信息化sso登录状态
	SSOLogInSuccess = 1
	SSOLogInFailed  = 2

	// 信息化mdm相关
	MDM_SEARCH_USERS = "/external/employeeSearch"
	MDM_QUERY_USERS  = "/external/employeeDetail"

	// 权限日志
	OperationCreateUser  = 1 //添加用户
	OperationDeleteUser  = 2 //删除用户
	OperationCreateRole  = 3 //添加角色标签
	OperationUpdateRole  = 4 //更新角色标签
	OperationDeleteRole  = 5 //删除角色标签
	OperationRestoreUser = 6 //解禁用户
	OperationBanUser     = 7 //封禁用户
	OperationAssignRole  = 8 //用户被分配角色
	OperationDepriveRole = 9 //用户被删除角色

	// 日志详情模版
	LogContentTempCreateUser  = "用户%s被创建 "
	LogContentTempDeleteUser  = "用户%s被删除"
	LogContentTempCreateRole  = "角色%s被创建"
	LogContentTempUpdateRole  = "角色%s被更新"
	LogContentTempDeleteRole  = "角色%s被删除"
	LogContentTempRestoreUser = "用户%s被解除禁用"
	LogContentTempBanUser     = "用户%s被禁用"
	LogContentTempAssignRole  = "用户%s被分配%s角色"
	LogContentTempDepriveRole = "用户%s被删除%s角色"

	// 权限点
	REMOTE_CONFIG_QUERY          = "remote-config-query"
	REMOTE_CONFIG_CREATE         = "remote-config-create"
	REMOTE_CONFIG_UPDATE         = "remote-config-update"
	REMOTE_CONFIG_DELETE         = "remote-config-delete"
	REMOTE_CONFIG_PUBLISH        = "remote-config-publish"
	REMOTE_CONFIG_CANCEL_PUBLISH = "remote-config-cancel-publish"

	ENV_MANAGE_QUERY  = "env-manage-query"
	ENV_MANAGE_CREATE = "env-manage-create"
	ENV_MANAGE_UPDATE = "env-manage-update"
	ENV_MANAGE_DELETE = "env-manage-delete"
)

// TODO, define resource urls
var OperaionResourceMap = map[string]map[string]bool{
	// remote config
	"/api.remoteconfig.RemoteConfig/CancelPublishRemoteConfig": {REMOTE_CONFIG_CANCEL_PUBLISH: true},
	"/api.remoteconfig.RemoteConfig/CreateRemoteConfig":        {REMOTE_CONFIG_CREATE: true},
	"/api.remoteconfig.RemoteConfig/DeleteRemoteConfig":        {REMOTE_CONFIG_DELETE: true},
	"/api.remoteconfig.RemoteConfig/GetRemoteConfigList":       {REMOTE_CONFIG_QUERY: true},
	"/api.remoteconfig.RemoteConfig/PublishRemoteConfig":       {REMOTE_CONFIG_PUBLISH: true},
	"/api.remoteconfig.RemoteConfig/UpdateRemoteConfig":        {REMOTE_CONFIG_UPDATE: true},

	// env manage
	"/api.projectconfig.envmanage.EnvManage/CreateEnv":  {ENV_MANAGE_CREATE: true},
	"/api.projectconfig.envmanage.EnvManage/DeleteEnv":  {ENV_MANAGE_DELETE: true},
	"/api.projectconfig.envmanage.EnvManage/GetEnvList": {ENV_MANAGE_QUERY: true},
	"/api.projectconfig.envmanage.EnvManage/UpdateEnv":  {ENV_MANAGE_UPDATE: true},
}

// TODO, 定义全权限点枚举数组
var AllResources []string = []string{
	REMOTE_CONFIG_QUERY,
	REMOTE_CONFIG_CREATE,
	REMOTE_CONFIG_UPDATE,
	REMOTE_CONFIG_DELETE,
	REMOTE_CONFIG_PUBLISH,
	REMOTE_CONFIG_CANCEL_PUBLISH,
}

type AdminUser struct {
	Id        uint32
	HgId      string
	Email     string
	Nickname  string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserAppRelation struct {
	Id         uint32
	UserId     uint32
	AppId      uint32
	Status     uint32
	Identity   uint32
	CreatedAt  time.Time
	UpdatedAt  time.Time
	AssignedAt time.Time
}

type AdminRole struct {
	Id                  uint32
	Name                string
	Description         string
	ResourceList        string
	Status              uint32
	AppId               uint32
	ResourceListContent []string
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type UserRoleRelation struct {
	Id        uint32
	UserId    uint32
	RoleId    uint32
	AppId     uint32
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SSOTokenData struct {
	Token string `json:"token"`
}

type SSOTokenRes struct {
	Code    uint32        `json:"code"`
	Message string        `json:"message"`
	Data    *SSOTokenData `json:"data"`
}

type SSOInfoRes struct {
	Code    uint32       `json:"code"`
	Message string       `json:"message"`
	Data    *SSOInfoData `json:"data"`
}

type SSOInfoData struct {
	Name      string `json:"name"`
	HgId      string `json:"hgid"`
	HgAccount string `json:"hgAccount"`
}

type MdmUserInfo struct {
	Name     string `json:"name"`
	Nickname string `json:"nickname"`
	Email    string `json:"enterpriseEmail"`
	HgId     string `json:"hgId"`
}

// MdmRes
// either searching users or batch pull user has same data
type MdmRes struct {
	Code    uint32         `json:"code"`
	Message string         `json:"message"`
	Data    []*MdmUserInfo `json:"data"`
}

type AuthLog struct {
	Id         uint32    //日志主键
	AppId      uint32    //产品id,1为明日方舟，2为玩家社区
	OperatorId uint32    //操作员id
	Operation  uint32    //操作类型
	Content    string    //操作详情，权限模块日志中该字段为一段可读模版，其他日志中该字段为空
	UserName   string    //冗余字段，用于模糊搜索
	RoleName   string    //冗余字段，用于模糊搜索
	UserId     uint32    //冗余字段
	RoleId     uint32    //冗余字段
	CreatedAt  time.Time //日志创建时间
	UpdatedAt  time.Time //日志更新时间
}
