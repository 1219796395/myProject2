package data

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/1219796395/myProject2/api/errorcode"
	"github.com/1219796395/myProject2/internal/biz"

	"github.com/go-sql-driver/mysql"
	"github.com/patrickmn/go-cache"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type remoteConfigRepo struct {
	data *Data
	log  *log.Helper
}

// NewRemoteConfigRepo .
func NewRemoteConfigRepo(data *Data, logger log.Logger) biz.RemoteConfigRepo {
	return &remoteConfigRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// RemoteConfig 远程配置-数据层 实体类
type RemoteConfig struct {
	Id              int64     `gorm:"column:id;primaryKey"` // id, auto increment
	Appid           uint32    `gorm:"column:appid"`
	Env             string    `gorm:"column:env"`
	Channel         string    `gorm:"column:channel"`
	Platform        string    `gorm:"column:platform"`
	Name            string    `gorm:"column:name"`
	State           uint32    `gorm:"column:state"`
	Data            string    `gorm:"column:data"`
	DataVersion     uint64    `gorm:"column:data_version"`
	ModifyData      string    `gorm:"column:modify_data"`
	NotDistChannel  int32     `gorm:"column:not_dist_channel"`
	NotDistPlatform int32     `gorm:"column:not_dist_platform"`
	Operator        string    `gorm:"column:operator"`
	OperatorID      uint32    `gorm:"column:operator_id"`
	Ctime           time.Time `gorm:"column:ctime;->"`
	Mtime           time.Time `gorm:"column:mtime;->"`
}

func (RemoteConfig) TableName() string {
	return "remote_config"
}

func convToBizRemoteConfig(rc *RemoteConfig) *biz.RemoteConfig {
	return &biz.RemoteConfig{
		Appid:           rc.Appid,
		Env:             rc.Env,
		Channel:         rc.Channel,
		Platform:        rc.Platform,
		Name:            rc.Name,
		State:           rc.State,
		Data:            rc.Data,
		DataVersion:     rc.DataVersion,
		ModifyData:      rc.ModifyData,
		NotDistChannel:  rc.NotDistChannel,
		NotDistPlatform: rc.NotDistPlatform,
		Operator:        rc.Operator,
		OperatorID:      rc.OperatorID,
		Ctime:           rc.Ctime,
		Mtime:           rc.Mtime,
	}
}

func batchConvToBizRemoteConfig(rcs []*RemoteConfig) []*biz.RemoteConfig {
	rcBizs := make([]*biz.RemoteConfig, 0, len(rcs))
	for _, rc := range rcs {
		rcBizs = append(rcBizs, convToBizRemoteConfig(rc))
	}
	return rcBizs
}

func convToDBRemoteConfig(rc *biz.RemoteConfig) *RemoteConfig {
	return &RemoteConfig{
		Appid:           rc.Appid,
		Env:             rc.Env,
		Channel:         rc.Channel,
		Platform:        rc.Platform,
		Name:            rc.Name,
		State:           rc.State,
		Data:            rc.Data,
		DataVersion:     rc.DataVersion,
		ModifyData:      rc.ModifyData,
		NotDistChannel:  rc.NotDistChannel,
		NotDistPlatform: rc.NotDistPlatform,
		Operator:        rc.Operator,
		OperatorID:      rc.OperatorID,
		Ctime:           rc.Ctime,
		Mtime:           rc.Mtime,
	}
}

func (r *remoteConfigRepo) GetRemoteConfig(ctx context.Context, rc *biz.RemoteConfig) (*biz.RemoteConfig, error) {
	rcDB := &RemoteConfig{}
	// 查询db的配置信息
	if err := r.data.db.WithContext(ctx).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `name`=? and `state` != ?", rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name, biz.RemoteConfigDeleted).First(rcDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorcode.ErrorRemoteConfigNotFound("remote config not found")
		}
		return nil, err
	}
	rc = convToBizRemoteConfig(rcDB)
	return rc, nil
}

// GetRemoteConfigList 批量获取远程配置列表
func (r *remoteConfigRepo) GetRemoteConfigList(ctx context.Context, rc *biz.RemoteConfig) (
	[]*biz.RemoteConfig, error) {
	rcDBs := []*RemoteConfig{}
	find := " `appid`=? "
	args := []interface{}{rc.Appid}
	if rc.Env != "" {
		find += " and `env`=? "
		args = append(args, rc.Env)
	}
	if rc.Channel != "" {
		find += " and `channel`=? "
		args = append(args, rc.Channel)
	}
	if rc.Platform != "" {
		find += " and `platform`=? "
		args = append(args, rc.Platform)
	}
	if rc.Name != "" {
		find += " and `name`=? "
		args = append(args, rc.Name)
	}
	find += " and `state`!= ?"
	args = append(args, biz.RemoteConfigDeleted)
	if err := r.data.db.WithContext(ctx).Where(find, args...).Find(&rcDBs).Error; err != nil {
		return nil, err
	}
	return batchConvToBizRemoteConfig(rcDBs), nil
}

// PageGetRemoteConfigList 分页批量获取远程配置列表
func (r *remoteConfigRepo) PageGetRemoteConfigList(ctx context.Context, rc *biz.RemoteConfig, start, size int) (
	[]*biz.RemoteConfig, error) {
	rcDBs := []*RemoteConfig{}
	find := " `appid`=? "
	args := []interface{}{rc.Appid}
	if rc.Env != "" {
		find += " and `env`=? "
		args = append(args, rc.Env)
	}
	if rc.Channel != "" {
		find += " and `channel`=? "
		args = append(args, rc.Channel)
	}
	if rc.Platform != "" {
		find += " and `platform`=? "
		args = append(args, rc.Platform)
	}
	find += " and `state`!= ?"
	args = append(args, biz.RemoteConfigDeleted)
	if err := r.data.db.WithContext(ctx).Where(find, args...).Offset(start).Limit(size).Find(&rcDBs).Error; err != nil {
		return nil, err
	}
	return batchConvToBizRemoteConfig(rcDBs), nil
}

// PageGetRemoteConfigList 分页批量获取远程配置列表
func (r *remoteConfigRepo) PageGetRemoteConfigListWithStates(ctx context.Context, rc *biz.RemoteConfig, states []uint32, start, size int) (
	[]*biz.RemoteConfig, error) {
	rcDBs := []*RemoteConfig{}
	findList := []string{}
	args := []interface{}{}
	if rc.Appid != 0 {
		findList = append(findList, " `appid`=? ")
		args = append(args, rc.Appid)
	}
	if rc.Env != "" {
		findList = append(findList, " `env`=? ")
		args = append(args, rc.Env)
	}
	if rc.Channel != "" {
		findList = append(findList, " `channel`=? ")
		args = append(args, rc.Channel)
	}
	if rc.Platform != "" {
		findList = append(findList, " `platform`=? ")
		args = append(args, rc.Platform)
	}

	if len(states) > 0 {
		findList = append(findList, " `state` in (?) ")
		args = append(args, states)
	}
	find := strings.Join(findList, "and")
	if err := r.data.db.WithContext(ctx).Where(find, args...).Offset(start).Limit(size).Find(&rcDBs).Error; err != nil {
		return nil, err
	}
	return batchConvToBizRemoteConfig(rcDBs), nil
}

// CreateRemoteConfig 创建远程配置
func (r *remoteConfigRepo) CreateRemoteConfig(ctx context.Context, rc *biz.RemoteConfig) error {
	rcDB := &RemoteConfig{}
	isCreate := true
	// 查询db的配置信息
	if err := r.data.db.WithContext(ctx).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `name`=?", rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name).First(rcDB).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		isCreate = false
	}

	if isCreate && rcDB.State != biz.RemoteConfigDeleted {
		return errorcode.ErrorRemoteConfigAlreadyExist("remote config already exist")
	}
	rc.State = biz.RemoteConfigNeverPublished
	rc.ModifyData = rc.Data
	rc.Data = "{}"
	rc.DataVersion = 0
	now := time.Now()
	rc.Ctime, rc.Mtime = now, now
	rcDB = convToDBRemoteConfig(rc)
	if !isCreate {
		if err := r.data.db.WithContext(ctx).Create(rcDB).Error; err != nil {
			if errMySQL, ok := err.(*mysql.MySQLError); ok && errMySQL.Number == 1062 {
				return errorcode.ErrorRemoteConfigAlreadyExist("remote config already exist")
			}
			return err
		}
	}
	return r.data.db.WithContext(ctx).Model(rcDB).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `name`=?", rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name).Updates(rcDB).Error
}

// UpdateRemoteConfig 修改远程配置
func (r *remoteConfigRepo) UpdateRemoteConfig(ctx context.Context, beforeRc, rc *biz.RemoteConfig) error {
	// 根据db的配置信息，判断修改后的状态
	rc.ModifyData = rc.Data
	rc.State = biz.RemoteConfigWaitPublish
	if beforeRc != nil {
		if beforeRc.State == biz.RemoteConfigNeverPublished {
			rc.State = biz.RemoteConfigNeverPublished
		}
		rc.Data = beforeRc.Data
		rc.DataVersion = beforeRc.DataVersion
		rc.Ctime, rc.Mtime = beforeRc.Ctime, time.Now()
	}
	// 在db上修改远程配置
	if err := r.data.db.WithContext(ctx).Model(&RemoteConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `name`=?", rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name).Updates(map[string]interface{}{
		"state":       rc.State,
		"modify_data": rc.ModifyData,
		"operator":    rc.Operator,
	}).Error; err != nil {
		return err
	}
	return nil
}

// DeleteRemoteConfig 删除远程配置
func (r *remoteConfigRepo) DeleteRemoteConfig(ctx context.Context, beforeRc, rc *biz.RemoteConfig) error {
	// 在db上删除远程配置
	rc.State = biz.RemoteConfigDeleted
	if beforeRc != nil {
		rc.Data = beforeRc.Data
		rc.DataVersion = beforeRc.DataVersion
		rc.ModifyData = beforeRc.ModifyData
		rc.Ctime, rc.Mtime = beforeRc.Ctime, time.Now()
	}
	if err := r.data.db.WithContext(ctx).Model(&RemoteConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `name`=?", rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name).Updates(map[string]interface{}{
		"state":    rc.State,
		"operator": rc.Operator,
	}).Error; err != nil {
		return err
	}
	return nil
}

// PublishRemoteConfig 发布远程配置
func (r *remoteConfigRepo) PublishRemoteConfig(ctx context.Context, beforeRc, rc *biz.RemoteConfig) error {
	// 发布远程配置，必须传入发布前的远程配置信息
	if beforeRc == nil {
		return errorcode.ErrorInternalServerError("internal server err")
	}
	// 在db上发布远程配置
	rc.State = biz.RemoteConfigOnline
	rc.Data = beforeRc.ModifyData
	now := time.Now()
	rc.DataVersion = uint64(now.UnixMilli())
	rc.ModifyData = beforeRc.ModifyData
	rc.Ctime, rc.Mtime = beforeRc.Ctime, now
	if err := r.data.db.WithContext(ctx).Model(&RemoteConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `name`=?", rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name).Updates(map[string]interface{}{
		"state":        rc.State,
		"data":         rc.Data,
		"data_version": rc.DataVersion,
		"operator":     rc.Operator,
	}).Error; err != nil {
		return err
	}

	return nil
}

// CancelPublishRemoteConfig 取消发布远程配置
func (r *remoteConfigRepo) CancelPublishRemoteConfig(ctx context.Context, beforeRc, rc *biz.RemoteConfig) error {
	// 在db上取消远程配置
	rc.State = biz.RemoteConfigNeverPublished
	rc.Data = "{}"
	now := time.Now()
	rc.DataVersion = uint64(now.UnixMilli())
	if beforeRc != nil {
		rc.ModifyData = beforeRc.ModifyData
		rc.Ctime, rc.Mtime = beforeRc.Ctime, now
	}
	if err := r.data.db.WithContext(ctx).Model(&RemoteConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `name`=?", rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name).Updates(map[string]interface{}{
		"state":        rc.State,
		"data":         rc.Data,
		"data_version": rc.DataVersion,
		"operator":     rc.Operator,
	}).Error; err != nil {
		return err
	}
	return nil
}

func genRemoteConfigKey(appid uint32, env, channel, platform, configName string) string {
	return fmt.Sprintf("game_config_remote_config_%d_%s_%s_%s_%s",
		appid, env, channel, platform, configName)
}

// GetRemoteConfig C端获取单个远程配置
func (r *remoteConfigRepo) GetRemoteConfigCache(ctx context.Context, appid uint32, env, channel, platform, configName string) (string, error) {
	key := genRemoteConfigKey(appid, env, channel, platform, configName)
	data, ok := r.data.cache.Get(key)
	if !ok {
		return "", errorcode.ErrorRemoteConfigNotFound("remote config not found")
	}
	rc, ok := data.(*biz.RemoteConfig)
	if !ok {
		return "", errorcode.ErrorRemoteConfigNotFound("remote config not found")
	}
	return rc.Data, nil
}

func (r *remoteConfigRepo) RemoveRemoteConfigCache(ctx context.Context, rc *biz.RemoteConfig) error {
	// 在redis上删除远程配置
	key := genRemoteConfigKey(rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name)
	r.data.cache.Delete(key)
	return nil
}

func (r *remoteConfigRepo) SetRemoteConfigCache(ctx context.Context, rc *biz.RemoteConfig) error {
	// 在redis上更新远程配置
	key := genRemoteConfigKey(rc.Appid, rc.Env, rc.Channel, rc.Platform, rc.Name)
	r.data.cache.Set(key, rc, cache.NoExpiration)
	return nil
}
