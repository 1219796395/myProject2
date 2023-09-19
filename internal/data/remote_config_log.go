package data

import (
	"context"
	"time"

	"github.com/1219796395/myProject2/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type remoteConfigLogRepo struct {
	data *Data
	log  *log.Helper
}

// NewRemoteConfigLogRepo .
func NewRemoteConfigLogRepo(data *Data, logger log.Logger) biz.RemoteConfigLogRepo {
	return &remoteConfigLogRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// RemoteConfigLog 远程配置日志-数据层 实体类
type RemoteConfigLog struct {
	Id           int64     `gorm:"column:id;primaryKey"` // id, auto increment
	Appid        uint32    `gorm:"column:appid"`
	Env          string    `gorm:"column:env"`
	Channel      string    `gorm:"column:channel"`
	Platform     string    `gorm:"column:platform"`
	ConfigName   string    `gorm:"column:config_name"`
	Operation    uint32    `gorm:"column:operation"`
	MethodName   string    `gorm:"column:method_name"`
	Operator     string    `gorm:"column:operator"`
	UpdateBefore string    `gorm:"column:update_before"`
	UpdateAfter  string    `gorm:"column:update_after"`
	Ctime        time.Time `gorm:"column:ctime;->"`
}

func (RemoteConfigLog) TableName() string {
	return "remote_config_log"
}

func convToDBRemoteConfigLog(rcLog *biz.RemoteConfigLog) *RemoteConfigLog {
	return &RemoteConfigLog{
		Appid:        rcLog.Appid,
		Env:          rcLog.Env,
		Channel:      rcLog.Channel,
		Platform:     rcLog.Platform,
		ConfigName:   rcLog.ConfigName,
		Operation:    rcLog.Operation,
		MethodName:   rcLog.MethodName,
		Operator:     rcLog.Operator,
		UpdateBefore: rcLog.UpdateBefore,
		UpdateAfter:  rcLog.UpdateAfter,
		Ctime:        rcLog.Ctime,
	}
}

func batchConvToBizRemoteConfigLog(rcLogs []*RemoteConfigLog) []*biz.RemoteConfigLog {
	rcLogBizs := make([]*biz.RemoteConfigLog, 0, len(rcLogs))
	for _, rcLog := range rcLogs {
		rcLogBizs = append(rcLogBizs, convToBizRemoteConfigLog(rcLog))
	}
	return rcLogBizs
}

func convToBizRemoteConfigLog(rcLog *RemoteConfigLog) *biz.RemoteConfigLog {
	return &biz.RemoteConfigLog{
		Appid:        rcLog.Appid,
		Env:          rcLog.Env,
		Channel:      rcLog.Channel,
		Platform:     rcLog.Platform,
		ConfigName:   rcLog.ConfigName,
		Operation:    rcLog.Operation,
		MethodName:   rcLog.MethodName,
		Operator:     rcLog.Operator,
		UpdateBefore: rcLog.UpdateBefore,
		UpdateAfter:  rcLog.UpdateAfter,
		Ctime:        rcLog.Ctime,
	}
}

func (r *remoteConfigLogRepo) CreateRemoteConfigLog(ctx context.Context, log *biz.RemoteConfigLog) error {
	logDB := convToDBRemoteConfigLog(log)
	return r.data.db.WithContext(ctx).Create(logDB).Error
}

func (r *remoteConfigLogRepo) GetRemoteConfigLogList(ctx context.Context, rcLog *biz.RemoteConfigLog) (
	[]*biz.RemoteConfigLog, error) {
	rcLogDBs := []*RemoteConfigLog{}
	find := "`appid`=? "
	args := []interface{}{rcLog.Appid}
	if rcLog.Env != "" {
		find += " and `env`=? "
		args = append(args, rcLog.Env)
	}
	if rcLog.Channel != "" {
		find += " and `channel`=? "
		args = append(args, rcLog.Channel)
	}
	if rcLog.Platform != "" {
		find += " and `platform`=? "
		args = append(args, rcLog.Platform)
	}
	if rcLog.ConfigName != "" {
		find += " and `config_name`=? "
		args = append(args, rcLog.ConfigName)
	}
	if rcLog.Operation != 0 {
		find += " and `operation`=? "
		args = append(args, rcLog.Operation)
	}
	if rcLog.Operator != "" {
		find += " and `operator`=? "
		args = append(args, rcLog.Operator)
	}
	if rcLog.StartTime != 0 {
		find += " and `ctime`>=? "
		args = append(args, time.Unix(0, int64(rcLog.StartTime)*1e6))
	}
	if rcLog.EndTime != 0 {
		find += " and `ctime`<? "
		args = append(args, time.Unix(0, int64(rcLog.EndTime)*1e6))
	}

	if err := r.data.db.WithContext(ctx).Where(find, args...).Find(&rcLogDBs).Error; err != nil {
		return nil, err
	}
	return batchConvToBizRemoteConfigLog(rcLogDBs), nil
}
