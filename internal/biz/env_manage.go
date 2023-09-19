package biz

import (
	"context"
	"time"

	"github.com/1219796395/myProject2/api/errorcode"

	"github.com/go-kratos/kratos/v2/log"
)

// EnvManageLogic 环境管理-逻辑层
type EnvManageLogic struct {
	repo   EnvManageRepo
	rcRepo RemoteConfigRepo
	log    *log.Helper
}

// NewEnvManageLogic 新建一个 环境管理-逻辑层 对象
func NewEnvManageLogic(repo EnvManageRepo, rcRepo RemoteConfigRepo, logger log.Logger) *EnvManageLogic {
	return &EnvManageLogic{repo: repo, rcRepo: rcRepo, log: log.NewHelper(logger)}
}

// EnvManage 环境-逻辑层实体类
type Env struct {
	Appid    uint32
	Field    string
	Name     string
	NewName  string
	Remark   string
	IsPreset bool
	Operator string
	State    uint32
	Ctime    time.Time
	Mtime    time.Time
}

// EnvManageRepo 环境管理-数据层
type EnvManageRepo interface {
	GetEnvList(context.Context, *Env) ([]*Env, error)
	GetEnv(context.Context, *Env) (*Env, error)
	CreateEnv(context.Context, *Env) error
	UpdateEnv(context.Context, *Env) error
	DeleteEnv(context.Context, *Env) error
}

func (l *EnvManageLogic) GetEnvList(ctx context.Context, env *Env) ([]*Env, error) {
	envs, err := l.repo.GetEnvList(ctx, env)
	if err != nil {
		return nil, err
	}
	return envs, nil
}

func (l *EnvManageLogic) CreateEnv(ctx context.Context, env *Env) error {
	// 1.查看是否存在该环境
	_, err := l.repo.GetEnv(ctx, env)
	if err == nil { // 存在该环境，返回错误
		return errorcode.ErrorRemoteConfigAlreadyExist("env name already exist")
	}
	if !errorcode.IsEnvManageNotFound(err) { // 若不是 not found 报错，则直接返回
		return err
	}
	return l.repo.CreateEnv(ctx, env)
}

func (l *EnvManageLogic) UpdateEnv(ctx context.Context, env *Env) error {
	// 1.查看是否存在该环境
	envDB, err := l.repo.GetEnv(ctx, env)
	if err != nil {
		return err
	}
	// 2.修改或删除前需确认，该环境是否存在配置
	rc := &RemoteConfig{
		Appid: env.Appid,
		Env:   env.Name,
	}
	rcs, err := l.rcRepo.GetRemoteConfigList(ctx, rc)
	if err != nil {
		return err
	}
	if len(rcs) > 0 {
		return errorcode.ErrorEnvManageNotModifyEnvHasConf("can not modify env which has conf")
	}

	// 3. 判断请求是否超时
	if env.Mtime.UnixMilli() != 0 && env.Mtime.UnixMilli() < envDB.Mtime.UnixMilli() {
		return errorcode.ErrorRemoteConfigExpireRequest("expire request")
	}

	// 4. 修改环境信息
	return l.repo.UpdateEnv(ctx, env)
}

func (l *EnvManageLogic) DeleteEnv(ctx context.Context, env *Env) error {
	// 1.查看是否存在该环境
	envDB, err := l.repo.GetEnv(ctx, env)
	if err != nil {
		return err
	}
	// 2.修改或删除前需确认，该环境是否存在配置
	rc := &RemoteConfig{
		Appid: env.Appid,
		Env:   env.Name,
	}
	rcs, err := l.rcRepo.GetRemoteConfigList(ctx, rc)
	if err != nil {
		return err
	}
	if len(rcs) > 0 {
		return errorcode.ErrorEnvManageNotModifyEnvHasConf("can not modify env which has conf")
	}

	// 3. 判断请求是否超时
	if env.Mtime.UnixMilli() != 0 && env.Mtime.UnixMilli() < envDB.Mtime.UnixMilli() {
		return errorcode.ErrorRemoteConfigExpireRequest("expire request")
	}

	// 4. 删除环境信息
	return l.repo.DeleteEnv(ctx, env)
}

// 校验是否存在该环境名
func existEnv(ctx context.Context, repo EnvManageRepo, appid uint32, envName string) error {
	env := &Env{
		Appid: appid,
		Name:  envName,
	}
	if _, err := repo.GetEnv(ctx, env); err != nil {
		if errorcode.IsEnvManageNotFound(err) {
			return errorcode.ErrorBadRequest("invalid env")
		}
		return err
	}
	return nil
}
