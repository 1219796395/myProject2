package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// NetworkConfigLogLogic 配置日志-逻辑层
type NetworkConfigLogLogic struct {
	repo NetworkConfigLogRepo
	log  *log.Helper
}

// NewNetworkConfigLogLogic 新建一个 network_config日志-逻辑层 对象
func NewNetworkConfigLogLogic(repo RemoteConfigLogRepo, logger log.Logger) *RemoteConfigLogLogic {
	return &RemoteConfigLogLogic{repo: repo, log: log.NewHelper(logger)}
}

// NetworkConfigLog network_config日志-逻辑层实体类
type NetworkConfigLog struct {
	Appid        uint32
	Env          string
	Channel      string
	Platform     string
	Version      string
	MajorVersion uint32
	MinorVersion uint32
	Operation    uint32
	MethodName   string
	Operator     string
	UpdateBefore string
	UpdateAfter  string
	Ctime        time.Time
	StartTime    uint64
	EndTime      uint64
}

// NetworkConfigLogRepo 配置日志-数据层
type NetworkConfigLogRepo interface {
	CreateNetworkConfigLog(context.Context, *NetworkConfigLog) error
}
