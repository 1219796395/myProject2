package biz

import (
	"context"
	"path/filepath"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// RemoteConfigLogLogic 配置日志-逻辑层
type RemoteConfigLogLogic struct {
	repo RemoteConfigLogRepo
	log  *log.Helper
}

// NewRemoteConfigLogLogic 新建一个 远程配置日志-逻辑层 对象
func NewRemoteConfigLogLogic(repo RemoteConfigLogRepo, logger log.Logger) *RemoteConfigLogLogic {
	return &RemoteConfigLogLogic{repo: repo, log: log.NewHelper(logger)}
}

// RemoteConfigLog 远程配置日志-逻辑层实体类
type RemoteConfigLog struct {
	Appid        uint32
	Env          string
	Channel      string
	Platform     string
	ConfigName   string
	Operation    uint32
	MethodName   string
	Operator     string
	UpdateBefore string
	UpdateAfter  string
	Ctime        time.Time
	StartTime    uint64
	EndTime      uint64
}

// RemoteConfigLogRepo 配置日志-数据层
type RemoteConfigLogRepo interface {
	GetRemoteConfigLogList(context.Context, *RemoteConfigLog) ([]*RemoteConfigLog, error)
	CreateRemoteConfigLog(context.Context, *RemoteConfigLog) error
}

const (
	operationCreate        = 1
	operationUpdate        = 2
	operationDelete        = 3
	operationPublish       = 4
	operationCancelPublish = 5
)

func convRemoteConfigToLog(ctx context.Context, operation uint32, before, after *RemoteConfig) *RemoteConfigLog {
	appid := before.Appid
	if before.Appid == 0 {
		appid = after.Appid
	}
	env := before.Env
	if before.Env == "" {
		env = after.Env
	}
	channel := before.Channel
	if before.Channel == "" {
		channel = after.Channel
	}
	platform := before.Platform
	if before.Platform == "" {
		platform = after.Platform
	}
	updateBefore, err := json.Marshal(before)
	if err != nil {
		log.Errorf("json marshal err = %+v", err)
	}
	updateAfter, err := json.Marshal(after)
	if err != nil {
		log.Errorf("json marshal err = %+v", err)
	}
	methodName := ""
	if ts, ok := transport.FromServerContext(ctx); ok {
		methodName = filepath.Base(ts.Operation())
	}
	return &RemoteConfigLog{
		Appid:        appid,
		Env:          env,
		Channel:      channel,
		Platform:     platform,
		ConfigName:   after.Name,
		Operation:    operation,
		MethodName:   methodName,
		Operator:     after.Operator,
		UpdateBefore: string(updateBefore),
		UpdateAfter:  string(updateAfter),
	}
}

func convNetworkConfigToLog(ctx context.Context, operation uint32, before, after *NetworkConfig) *NetworkConfigLog {
	appid := before.Appid
	if before.Appid == 0 {
		appid = after.Appid
	}
	env := before.Env
	if before.Env == "" {
		env = after.Env
	}
	channel := before.Channel
	if before.Channel == "" {
		channel = after.Channel
	}
	platform := before.Platform
	if before.Platform == "" {
		platform = after.Platform
	}
	updateBefore, err := json.Marshal(before)
	if err != nil {
		log.Errorf("json marshal err = %+v", err)
	}
	updateAfter, err := json.Marshal(after)
	if err != nil {
		log.Errorf("json marshal err = %+v", err)
	}
	methodName := ""
	if ts, ok := transport.FromServerContext(ctx); ok {
		methodName = filepath.Base(ts.Operation())
	}
	return &NetworkConfigLog{
		Appid:        appid,
		Env:          env,
		Channel:      channel,
		Platform:     platform,
		Version:      after.Version,
		MajorVersion: after.MajorVersion,
		MinorVersion: after.MinorVersion,
		Operation:    operation,
		MethodName:   methodName,
		Operator:     after.Operator,
		UpdateBefore: string(updateBefore),
		UpdateAfter:  string(updateAfter),
	}
}

func (l *RemoteConfigLogLogic) GetRemoteConfigLogList(ctx context.Context, rcLog *RemoteConfigLog) (
	[]*RemoteConfigLog, error) {
	return l.repo.GetRemoteConfigLogList(ctx, rcLog)
}
