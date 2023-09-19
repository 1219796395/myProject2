package data

import (
	"context"
	"game-config/internal/biz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type networkConfigLogRepo struct {
	data *Data
	log  *log.Helper
}

// NewNetworkConfigLogRepo .
func NewNetworkConfigLogRepo(data *Data, logger log.Logger) biz.NetworkConfigLogRepo {
	return &networkConfigLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// NetworkConfigLog network_config日志-数据层 实体类
type NetworkConfigLog struct {
	Id           int64     `gorm:"column:id;primaryKey"` // id, auto increment
	Appid        uint32    `gorm:"column:appid"`
	Env          string    `gorm:"column:env"`
	Channel      string    `gorm:"column:channel"`
	Platform     string    `gorm:"column:platform"`
	Version      string    `gorm:"column:version"`
	MajorVersion uint32    `gorm:"column:major_version"`
	MinorVersion uint32    `gorm:"column:minor_version"`
	Operation    uint32    `gorm:"column:operation"`
	MethodName   string    `gorm:"column:method_name"`
	Operator     string    `gorm:"column:operator"`
	UpdateBefore string    `gorm:"column:update_before"`
	UpdateAfter  string    `gorm:"column:update_after"`
	Ctime        time.Time `gorm:"column:ctime;->"`
}

func (NetworkConfigLog) TableName() string {
	return "network_config_log"
}

func convToDBNetworkConfigLog(ncLog *biz.NetworkConfigLog) *NetworkConfigLog {
	return &NetworkConfigLog{
		Appid:        ncLog.Appid,
		Env:          ncLog.Env,
		Channel:      ncLog.Channel,
		Platform:     ncLog.Platform,
		Version:      ncLog.Version,
		MajorVersion: ncLog.MajorVersion,
		MinorVersion: ncLog.MinorVersion,
		Operation:    ncLog.Operation,
		MethodName:   ncLog.MethodName,
		Operator:     ncLog.Operator,
		UpdateBefore: ncLog.UpdateBefore,
		UpdateAfter:  ncLog.UpdateAfter,
		Ctime:        ncLog.Ctime,
	}
}

func batchConvToBizNetworkConfigLog(ncLogs []*NetworkConfigLog) []*biz.NetworkConfigLog {
	ncLogBizs := make([]*biz.NetworkConfigLog, 0, len(ncLogs))
	for _, ncLog := range ncLogs {
		ncLogBizs = append(ncLogBizs, convToBizNetworkConfigLog(ncLog))
	}
	return ncLogBizs
}

func convToBizNetworkConfigLog(ncLog *NetworkConfigLog) *biz.NetworkConfigLog {
	return &biz.NetworkConfigLog{
		Appid:        ncLog.Appid,
		Env:          ncLog.Env,
		Channel:      ncLog.Channel,
		Platform:     ncLog.Platform,
		Version:      ncLog.Version,
		MajorVersion: ncLog.MajorVersion,
		MinorVersion: ncLog.MinorVersion,
		Operation:    ncLog.Operation,
		MethodName:   ncLog.MethodName,
		Operator:     ncLog.Operator,
		UpdateBefore: ncLog.UpdateBefore,
		UpdateAfter:  ncLog.UpdateAfter,
		Ctime:        ncLog.Ctime,
	}
}

func (r *networkConfigLogRepo) CreateNetworkConfigLog(ctx context.Context, log *biz.NetworkConfigLog) error {
	logDB := convToDBNetworkConfigLog(log)
	return r.data.db.WithContext(ctx).Create(logDB).Error
}

func (r *networkConfigLogRepo) GetNetworkConfigLogList(ctx context.Context, ncLog *biz.NetworkConfigLog) (
	[]*biz.NetworkConfigLog, error) {
	ncLogDBs := []*NetworkConfigLog{}
	find := "`appid`=? "
	args := []interface{}{ncLog.Appid}
	if ncLog.Env != "" {
		find += " and `env`=? "
		args = append(args, ncLog.Env)
	}
	if ncLog.Channel != "" {
		find += " and `channel`=? "
		args = append(args, ncLog.Channel)
	}
	if ncLog.Platform != "" {
		find += " and `platform`=? "
		args = append(args, ncLog.Platform)
	}
	if ncLog.Version != "" {
		find += " and `config_name`=? "
		args = append(args, ncLog.Version)
	}
	if ncLog.Operation != 0 {
		find += " and `operation`=? "
		args = append(args, ncLog.Operation)
	}
	if ncLog.Operator != "" {
		find += " and `operator`=? "
		args = append(args, ncLog.Operator)
	}
	if ncLog.StartTime != 0 {
		find += " and `ctime`>=? "
		args = append(args, time.Unix(0, int64(ncLog.StartTime)*1e6))
	}
	if ncLog.EndTime != 0 {
		find += " and `ctime`<? "
		args = append(args, time.Unix(0, int64(ncLog.EndTime)*1e6))
	}

	if err := r.data.db.WithContext(ctx).Where(find, args...).Find(&ncLogDBs).Error; err != nil {
		return nil, err
	}
	return batchConvToBizNetworkConfigLog(ncLogDBs), nil
}
