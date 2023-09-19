package biz

import (
	"context"
	"sync"
	"time"

	"github.com/1219796395/myProject2/api/errorcode"
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
)

// RemoteConfigLogic 远程配置-逻辑层
type RemoteConfigLogic struct {
	repo    RemoteConfigRepo
	envRepo EnvManageRepo
	logRepo RemoteConfigLogRepo
	log     *log.Helper
}

// NewRemoteConfigLogic 新建一个 远程配置-逻辑层 对象
func NewRemoteConfigLogic(repo RemoteConfigRepo, envRepo EnvManageRepo, logRepo RemoteConfigLogRepo, logger log.Logger) *RemoteConfigLogic {
	return &RemoteConfigLogic{repo: repo, envRepo: envRepo, logRepo: logRepo, log: log.NewHelper(logger)}
}

// RemoteConfig 远程配置-逻辑层实体类
type RemoteConfig struct {
	Appid           uint32    `json:"appid"`
	Env             string    `json:"env"`
	Channel         string    `json:"channel"`
	Platform        string    `json:"platform"`
	Name            string    `json:"name"`
	State           uint32    `json:"state"`
	Data            string    `json:"data"`
	DataVersion     uint64    `json:"data_version"`
	ModifyData      string    `json:"modify_data"`
	NotDistChannel  int32     `json:"not_dist_channel"`
	NotDistPlatform int32     `json:"not_dist_platform"`
	Operator        string    `json:"operator"`
	OperatorID      uint32    `json:"operator_id"`
	Ctime           time.Time `json:"ctime"`
	Mtime           time.Time `json:"mtime"`
}

// 远程配置状态
const (
	RemoteConfigNeverPublished = 1 // 未发布
	RemoteConfigWaitPublish    = 2 // 待发布
	RemoteConfigOnline         = 3 // 已上线
	RemoteConfigDeleted        = 4 // 已删除
)

// RemoteConfigRepo 远程配置-数据层
type RemoteConfigRepo interface {
	// db
	GetRemoteConfig(context.Context, *RemoteConfig) (*RemoteConfig, error)
	GetRemoteConfigList(context.Context, *RemoteConfig) ([]*RemoteConfig, error)
	PageGetRemoteConfigList(context.Context, *RemoteConfig, int, int) ([]*RemoteConfig, error)
	PageGetRemoteConfigListWithStates(context.Context, *RemoteConfig, []uint32, int, int) ([]*RemoteConfig, error)
	CreateRemoteConfig(context.Context, *RemoteConfig) error
	UpdateRemoteConfig(context.Context, *RemoteConfig, *RemoteConfig) error
	DeleteRemoteConfig(context.Context, *RemoteConfig, *RemoteConfig) error
	PublishRemoteConfig(context.Context, *RemoteConfig, *RemoteConfig) error
	CancelPublishRemoteConfig(context.Context, *RemoteConfig, *RemoteConfig) error

	// cache
	GetRemoteConfigCache(context.Context, uint32, string, string, string, string) (string, error)
	RemoveRemoteConfigCache(context.Context, *RemoteConfig) error
	SetRemoteConfigCache(context.Context, *RemoteConfig) error
}

// GetRemoteConfigList 获取远程配置列表
func (l *RemoteConfigLogic) GetRemoteConfigList(ctx context.Context, rc *RemoteConfig) (
	[]*RemoteConfig, error) {
	// 校验是否存在该环境名
	if err := existEnv(ctx, l.envRepo, rc.Appid, rc.Env); err != nil {
		return nil, err
	}
	return l.repo.GetRemoteConfigList(ctx, rc)
}

// CreateRemoteConfig 创建远程配置
func (l *RemoteConfigLogic) CreateRemoteConfig(ctx context.Context, rc *RemoteConfig) error {
	// 校验是否存在该环境名
	if err := existEnv(ctx, l.envRepo, rc.Appid, rc.Env); err != nil {
		return err
	}
	// 创建远程配置
	if err := l.repo.CreateRemoteConfig(ctx, rc); err != nil {
		return err
	}

	// 记录配置日志
	log := convRemoteConfigToLog(ctx, operationCreate, &RemoteConfig{}, rc)
	return l.logRepo.CreateRemoteConfigLog(ctx, log)
}

// DeleteRemoteConfig 删除远程配置
func (l *RemoteConfigLogic) DeleteRemoteConfig(ctx context.Context, rc *RemoteConfig) error {
	beforeRc, err := l.repo.GetRemoteConfig(ctx, rc)
	if err != nil {
		return err
	}

	// 未发布状态可用
	if beforeRc.State != RemoteConfigNeverPublished {
		return errorcode.ErrorBadRequest("not allow delete")
	}

	if rc.Mtime.UnixMilli() != 0 && rc.Mtime.UnixMilli() < beforeRc.Mtime.UnixMilli() {
		return errorcode.ErrorRemoteConfigExpireRequest("expire request")
	}

	// 删除远程配置
	if err := l.repo.DeleteRemoteConfig(ctx, beforeRc, rc); err != nil {
		return err
	}

	// 记录日志
	log := convRemoteConfigToLog(ctx, operationDelete, beforeRc, rc)
	return l.logRepo.CreateRemoteConfigLog(ctx, log)
}

// UpdateRemoteConfig 修改远程配置
func (l *RemoteConfigLogic) UpdateRemoteConfig(ctx context.Context, rc *RemoteConfig) error {
	beforeRc, err := l.repo.GetRemoteConfig(ctx, rc)
	if err != nil {
		return err
	}
	if rc.Mtime.UnixMilli() != 0 && rc.Mtime.UnixMilli() < beforeRc.Mtime.UnixMilli() {
		return errorcode.ErrorRemoteConfigExpireRequest("expire request")
	}

	// 修改远程配置
	if err := l.repo.UpdateRemoteConfig(ctx, beforeRc, rc); err != nil {
		return err
	}

	// 记录日志
	log := convRemoteConfigToLog(ctx, operationUpdate, beforeRc, rc)
	return l.logRepo.CreateRemoteConfigLog(ctx, log)
}

// PublishRemoteConfig 发布远程配置
func (l *RemoteConfigLogic) PublishRemoteConfig(ctx context.Context, rc *RemoteConfig) error {
	beforeRc, err := l.repo.GetRemoteConfig(ctx, rc)
	if err != nil {
		return err
	}

	// 未发布、待发布状态可用
	if beforeRc.State == RemoteConfigOnline {
		return errorcode.ErrorBadRequest("not allow publish")
	}

	if rc.Mtime.UnixMilli() != 0 && rc.Mtime.UnixMilli() < beforeRc.Mtime.UnixMilli() {
		return errorcode.ErrorRemoteConfigExpireRequest("expire request")
	}

	if err := l.repo.PublishRemoteConfig(ctx, beforeRc, rc); err != nil {
		return err
	}

	// 记录日志
	log := convRemoteConfigToLog(ctx, operationPublish, beforeRc, rc)
	return l.logRepo.CreateRemoteConfigLog(ctx, log)
}

// CancelPublishRemoteConfig 取消发布远程配置
func (l *RemoteConfigLogic) CancelPublishRemoteConfig(ctx context.Context, rc *RemoteConfig) error {
	beforeRc, err := l.repo.GetRemoteConfig(ctx, rc)
	if err != nil {
		return err
	}

	// 待发布、已上线状态可用
	if beforeRc.State == RemoteConfigNeverPublished {
		return errorcode.ErrorBadRequest("not allow cancel publish")
	}

	if rc.Mtime.UnixMilli() != 0 && rc.Mtime.UnixMilli() < beforeRc.Mtime.UnixMilli() {
		return errorcode.ErrorRemoteConfigExpireRequest("expire request")
	}

	// 取消发布配置
	if err := l.repo.CancelPublishRemoteConfig(ctx, beforeRc, rc); err != nil {
		return err
	}

	// 记录日志
	log := convRemoteConfigToLog(ctx, operationCancelPublish, beforeRc, rc)
	return l.logRepo.CreateRemoteConfigLog(ctx, log)
}

// GetRemoteConfig C端获取单个远程配置
func (l *RemoteConfigLogic) GetRemoteConfig(ctx context.Context, appid uint32, env, channel, platform, configName string) (string, error) {
	return l.repo.GetRemoteConfigCache(ctx, appid, env, channel, platform, configName)
}

// db与cache数据对账
func (l *RemoteConfigLogic) CheckCacheByDB(ctx context.Context) error {
	bc, err := conf.GetConf()
	if err != nil {
		return err
	}

	start := 0
	size := int(bc.GetBiz().RemoteConfigCheckCacheByDbTask.GetBatch())
	states := []uint32{
		RemoteConfigNeverPublished,
		RemoteConfigWaitPublish,
		RemoteConfigOnline,
		RemoteConfigDeleted,
	}

	wg := sync.WaitGroup{}
	for {
		// 分页获取remote_config信息
		rcs, err := l.repo.PageGetRemoteConfigListWithStates(ctx, &RemoteConfig{}, states, start, size)
		if err != nil {
			return err
		}
		if len(rcs) == 0 {
			break
		}

		expireRemoteConfigs := make([]*RemoteConfig, 0)
		onlineRemoteConfigs := make([]*RemoteConfig, 0)
		for _, rc := range rcs {
			if rc.State == RemoteConfigNeverPublished || rc.State == RemoteConfigDeleted {
				expireRemoteConfigs = append(expireRemoteConfigs, rc)
			}
			if rc.State == RemoteConfigWaitPublish || rc.State == RemoteConfigOnline {
				onlineRemoteConfigs = append(onlineRemoteConfigs, rc)
			}
		}

		for _, onlineRc := range onlineRemoteConfigs {
			wg.Add(1)
			rc := onlineRc
			go func() {
				defer wg.Done()
				if err := l.checkOnlineCache(ctx, rc); err != nil {
					l.log.Errorf("check online cache fail! onlineRemoteConfig = %+v, err = %+v", rc, err)
				}
			}()
		}

		for _, expireRc := range expireRemoteConfigs {
			wg.Add(1)
			rc := expireRc
			go func() {
				defer wg.Done()
				if err := l.checkExpireCache(ctx, rc); err != nil {
					l.log.Errorf("check expire cache fail! expireRemoteConfig = %+v, err = %+v", rc, err)
				}
			}()
		}
		wg.Wait()

		start += size
	}
	return nil
}

func (l *RemoteConfigLogic) checkOnlineCache(ctx context.Context, onlineRcDB *RemoteConfig) error {
	// 获取online的cache
	data, err := l.repo.GetRemoteConfigCache(ctx, onlineRcDB.Appid, onlineRcDB.Env, onlineRcDB.Channel,
		onlineRcDB.Platform, onlineRcDB.Name)
	if err != nil && !errorcode.IsRemoteConfigNotFound(err) {
		return err
	}
	if onlineRcDB.Data == data {
		return nil
	}
	return l.repo.SetRemoteConfigCache(ctx, onlineRcDB)
}

func (l *RemoteConfigLogic) checkExpireCache(ctx context.Context, expireRcDB *RemoteConfig) error {
	// 获取expire的cache
	_, err := l.repo.GetRemoteConfigCache(ctx, expireRcDB.Appid, expireRcDB.Env, expireRcDB.Channel,
		expireRcDB.Platform, expireRcDB.Name)
	if errorcode.IsRemoteConfigNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return l.repo.RemoveRemoteConfigCache(ctx, expireRcDB)
}
