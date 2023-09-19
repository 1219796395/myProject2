package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/1219796395/myProject2/api/errorcode"
	"github.com/1219796395/myProject2/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type networkConfigRepo struct {
	data *Data
	log  *log.Helper
}

// NewEnvManageRepo .
func NewNetworkConfigRepo(data *Data, logger log.Logger) biz.NetworkConfigRepo {
	return &networkConfigRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// NetworkConfig network_config-数据层 实体类
type NetworkConfig struct {
	Id           int64     `gorm:"column:id;primaryKey"` // id, auto increment
	Appid        uint32    `gorm:"column:appid"`
	Env          string    `gorm:"column:env"`
	Channel      string    `gorm:"column:channel"`
	Platform     string    `gorm:"column:platform"`
	Version      string    `gorm:"column:version"`
	MajorVersion uint32    `gorm:"column:major_version"`
	MinorVersion uint32    `gorm:"column:minor_version"`
	State        uint32    `gorm:"column:state"`
	AuditStart   uint64    `gorm:"column:audit_start"`
	AuditEnd     uint64    `gorm:"column:audit_end"`
	OnlineStart  uint64    `gorm:"column:online_start"`
	OnlineConfig string    `gorm:"column:online_config"`
	AuditConfig  string    `gorm:"column:audit_config"`
	ModifyData   string    `gorm:"column:modify_data"`
	Operator     string    `gorm:"column:operator"`
	Ctime        time.Time `gorm:"column:ctime;->"`
	Mtime        time.Time `gorm:"column:mtime;->"`
}

func (NetworkConfig) TableName() string {
	return "network_config"
}

func batchConvToBizNetworkConfig(ncs []*NetworkConfig) []*biz.NetworkConfig {
	res := make([]*biz.NetworkConfig, 0, len(ncs))
	for _, nc := range ncs {
		res = append(res, convToBizNetworkConfig(nc))
	}
	return res
}

func convToBizNetworkConfig(nc *NetworkConfig) *biz.NetworkConfig {
	modifyData := biz.NetworkConfigModify{}
	json.Unmarshal([]byte(nc.ModifyData), &modifyData)
	return &biz.NetworkConfig{
		Appid:        nc.Appid,
		Env:          nc.Env,
		Channel:      nc.Channel,
		Platform:     nc.Platform,
		Version:      nc.Version,
		MajorVersion: nc.MajorVersion,
		MinorVersion: nc.MinorVersion,
		State:        nc.State,
		AuditStart:   nc.AuditStart,
		AuditEnd:     nc.AuditEnd,
		OnlineStart:  nc.OnlineStart,
		OnlineConfig: nc.OnlineConfig,
		AuditConfig:  nc.AuditConfig,
		ModifyData:   modifyData,
		Operator:     nc.Operator,
		Ctime:        nc.Ctime,
		Mtime:        nc.Mtime,
	}
}

func convToDBNetworkConfig(nc *biz.NetworkConfig) *NetworkConfig {
	modifyData, _ := json.Marshal(&nc.ModifyData)
	return &NetworkConfig{
		Appid:        nc.Appid,
		Env:          nc.Env,
		Channel:      nc.Channel,
		Platform:     nc.Platform,
		Version:      nc.Version,
		MajorVersion: nc.MajorVersion,
		MinorVersion: nc.MinorVersion,
		State:        nc.State,
		AuditStart:   nc.AuditStart,
		AuditEnd:     nc.AuditEnd,
		OnlineStart:  nc.OnlineStart,
		OnlineConfig: nc.OnlineConfig,
		AuditConfig:  nc.AuditConfig,
		ModifyData:   string(modifyData),
		Operator:     nc.Operator,
		Ctime:        nc.Ctime,
		Mtime:        nc.Mtime,
	}
}

func (r *networkConfigRepo) GetNetworkConfig(ctx context.Context, nc *biz.NetworkConfig) (*biz.NetworkConfig, error) {
	ncDB := &NetworkConfig{}
	// 查询db的配置信息
	if err := r.data.db.WithContext(ctx).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `version`=? and `state` != ?", nc.Appid, nc.Env, nc.Channel, nc.Platform, nc.Version, biz.NetworkConfigDeleted).First(ncDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorcode.ErrorNetworkConfigNotFound("network config not found")
		}
		return nil, err
	}
	return convToBizNetworkConfig(ncDB), nil
}

func (r *networkConfigRepo) GetNetworkConfigList(ctx context.Context, nc *biz.NetworkConfig) ([]*biz.NetworkConfig, error) {
	ncDBs := []*NetworkConfig{}
	find := "`appid`=? and `env`=? "
	args := []interface{}{nc.Appid, nc.Env}
	if nc.Channel != "" {
		find += " and `channel`=? "
		args = append(args, nc.Channel)
	}
	if nc.Platform != "" {
		find += " and `platform`=? "
		args = append(args, nc.Platform)
	}
	find += " and `state`!= ?"
	args = append(args, biz.NetworkConfigDeleted)
	if err := r.data.db.WithContext(ctx).Where(find, args...).Find(&ncDBs).Error; err != nil {
		return nil, err
	}
	return batchConvToBizNetworkConfig(ncDBs), nil
}

func (r *networkConfigRepo) GetNetworkConfigListWithStates(ctx context.Context, nc *biz.NetworkConfig, states []uint32) ([]*biz.NetworkConfig, error) {
	ncDBs := []*NetworkConfig{}
	find := "`appid`=? and `env`=? "
	args := []interface{}{nc.Appid, nc.Env}
	if nc.Channel != "" {
		find += " and `channel`=? "
		args = append(args, nc.Channel)
	}
	if nc.Platform != "" {
		find += " and `platform`=? "
		args = append(args, nc.Platform)
	}
	if len(states) != 0 {
		find += " and `state` in (?) "
		args = append(args, states)
	}
	if err := r.data.db.WithContext(ctx).Where(find, args...).Find(&ncDBs).Error; err != nil {
		return nil, err
	}
	return batchConvToBizNetworkConfig(ncDBs), nil
}

func (r *networkConfigRepo) CreateNetworkConfig(ctx context.Context, nc *biz.NetworkConfig) error {
	ncDB := &NetworkConfig{}
	isCreate := true
	// 查询db的配置信息
	if err := r.data.db.WithContext(ctx).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `version`=?", nc.Appid, nc.Env, nc.Channel, nc.Platform, nc.Version).First(ncDB).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return err
		}
		isCreate = false
	}

	if isCreate && ncDB.State != biz.NetworkConfigDeleted {
		return errorcode.ErrorNetworkConfigAlreadyExist("network config already exist")
	}
	nc.State = biz.NetworkConfigNeverPublished
	nc.ModifyData = biz.NetworkConfigModify{
		AuditStart:   nc.AuditStart,
		AuditEnd:     nc.AuditEnd,
		OnlineStart:  nc.OnlineStart,
		OnlineConfig: nc.OnlineConfig,
		AuditConfig:  nc.AuditConfig,
	}
	nc.AuditStart = 0
	nc.AuditEnd = 0
	nc.OnlineStart = 0
	nc.OnlineConfig = "{}"
	nc.AuditConfig = "{}"
	now := time.Now()
	nc.Ctime, nc.Mtime = now, now
	ncDB = convToDBNetworkConfig(nc)
	if !isCreate {
		if err := r.data.db.WithContext(ctx).Create(ncDB).Error; err != nil {
			if errMySQL, ok := err.(*mysql.MySQLError); ok && errMySQL.Number == 1062 {
				return errorcode.ErrorNetworkConfigAlreadyExist("network config already exist")
			}
			return err
		}
	}
	return r.data.db.WithContext(ctx).Model(ncDB).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `version`=?", nc.Appid, nc.Env, nc.Channel, nc.Platform, nc.Version).Updates(ncDB).Error
}

func (r *networkConfigRepo) UpdateNetworkConfig(ctx context.Context, beforeNc *biz.NetworkConfig, nc *biz.NetworkConfig) error {
	// 根据db的配置信息，判断修改后的状态
	nc.ModifyData = biz.NetworkConfigModify{
		AuditStart:   nc.AuditStart,
		AuditEnd:     nc.AuditEnd,
		OnlineStart:  nc.OnlineStart,
		OnlineConfig: nc.OnlineConfig,
		AuditConfig:  nc.AuditConfig,
	}
	nc.State = biz.NetworkConfigWaitPublish
	if beforeNc != nil {
		if beforeNc.State == biz.NetworkConfigNeverPublished {
			nc.State = biz.NetworkConfigNeverPublished
		}
		nc.MajorVersion = beforeNc.MajorVersion
		nc.MinorVersion = beforeNc.MinorVersion
		nc.AuditStart = beforeNc.AuditStart
		nc.AuditEnd = beforeNc.AuditEnd
		nc.OnlineStart = beforeNc.OnlineStart
		nc.OnlineConfig = beforeNc.OnlineConfig
		nc.AuditConfig = beforeNc.AuditConfig
		nc.Ctime, nc.Mtime = beforeNc.Ctime, time.Now()
	}
	// 在db上修改远程配置
	modifyData, _ := json.Marshal(nc.ModifyData)
	if err := r.data.db.WithContext(ctx).Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `version`=?", nc.Appid, nc.Env, nc.Channel, nc.Platform, nc.Version).Updates(map[string]interface{}{
		"state":       nc.State,
		"modify_data": string(modifyData),
		"operator":    nc.Operator,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *networkConfigRepo) DeleteNetworkConfig(ctx context.Context, beforeNc *biz.NetworkConfig, nc *biz.NetworkConfig) error {
	// 在db上删除远程配置
	nc.State = biz.NetworkConfigDeleted
	if beforeNc != nil {
		nc.MajorVersion = beforeNc.MajorVersion
		nc.MinorVersion = beforeNc.MinorVersion
		nc.AuditStart = beforeNc.AuditStart
		nc.AuditEnd = beforeNc.AuditEnd
		nc.OnlineStart = beforeNc.OnlineStart
		nc.OnlineConfig = beforeNc.OnlineConfig
		nc.AuditConfig = beforeNc.AuditConfig
		nc.ModifyData = beforeNc.ModifyData
		nc.Ctime, nc.Mtime = beforeNc.Ctime, time.Now()
	}
	if err := r.data.db.WithContext(ctx).Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and ``=? and `platform`=? "+
		"and `version`=?", nc.Appid, nc.Env, nc.Channel, nc.Platform, nc.Version).Updates(map[string]interface{}{
		"state":    nc.State,
		"operator": nc.Operator,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *networkConfigRepo) PublishNetworkConfig(ctx context.Context, beforeNc *biz.NetworkConfig, nc *biz.NetworkConfig) error {
	// 发布远程配置，必须传入发布前的远程配置信息
	if beforeNc == nil {
		return errorcode.ErrorInternalServerError("internal server err")
	}
	now := time.Now()
	nc.State = biz.GetPublishedNetworkConfigState(uint64(now.UnixMilli()),
		beforeNc.ModifyData.AuditStart, beforeNc.ModifyData.AuditEnd, beforeNc.ModifyData.OnlineStart) // 需判断发布状态
	nc.MajorVersion = beforeNc.MajorVersion
	nc.MinorVersion = beforeNc.MinorVersion
	nc.AuditStart = beforeNc.ModifyData.AuditStart
	nc.AuditEnd = beforeNc.ModifyData.AuditEnd
	nc.OnlineStart = beforeNc.ModifyData.OnlineStart
	nc.OnlineConfig = beforeNc.ModifyData.OnlineConfig
	nc.AuditConfig = beforeNc.ModifyData.AuditConfig
	nc.ModifyData = beforeNc.ModifyData
	nc.Ctime, nc.Mtime = beforeNc.Ctime, now

	// 开启事务
	tx := r.data.db.WithContext(ctx).Begin()
	// 若状态是已上线，需把旧的已上线数据改为已过期
	if nc.State == biz.NetworkConfigOnline {
		if err := tx.Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
			" and `state`=? ", nc.Appid, nc.Env, nc.Channel, nc.Platform, biz.NetworkConfigOnline).
			Updates(map[string]interface{}{
				"state": biz.NetworkConfigExpire,
			}).Error; err != nil {
			return err
		}
	}
	// 发布db的network_config
	if err := tx.Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `version`=?", nc.Appid, nc.Env, nc.Channel, nc.Platform, nc.Version).Updates(map[string]interface{}{
		"state":         nc.State,
		"audit_start":   nc.AuditStart,
		"audit_end":     nc.AuditEnd,
		"online_start":  nc.OnlineStart,
		"online_config": nc.OnlineConfig,
		"audit_config":  nc.AuditConfig,
		"operator":      nc.Operator,
	}).Error; err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (r *networkConfigRepo) CancelPublishNetworkConfig(ctx context.Context, beforeNc *biz.NetworkConfig, nc *biz.NetworkConfig) error {
	// 在db上取消network_config
	nc.State = biz.NetworkConfigNeverPublished
	nc.AuditStart = 0
	nc.AuditEnd = 0
	nc.OnlineStart = 0
	nc.OnlineConfig = "{}"
	nc.AuditConfig = "{}"
	if beforeNc != nil {
		nc.MajorVersion = beforeNc.MajorVersion
		nc.MinorVersion = beforeNc.MinorVersion
		nc.ModifyData = beforeNc.ModifyData
		nc.Ctime, nc.Mtime = beforeNc.Ctime, time.Now()
	}
	if err := r.data.db.WithContext(ctx).Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
		"and `version`=?", nc.Appid, nc.Env, nc.Channel, nc.Platform, nc.Version).Updates(map[string]interface{}{
		"state":         nc.State,
		"audit_start":   nc.AuditStart,
		"audit_end":     nc.AuditEnd,
		"online_start":  nc.OnlineStart,
		"online_config": nc.OnlineConfig,
		"audit_config":  nc.AuditConfig,
		"operator":      nc.Operator,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (r *networkConfigRepo) PageGetCategoryList(ctx context.Context, start, size int) (
	[]*biz.NetworkConfig, error) {
	res := []*NetworkConfig{}
	// 查询db的配置信息
	if err := r.data.db.WithContext(ctx).Raw(
		"select appid,env,channel,platform "+
			"from network_config "+
			"group by appid,env,channel,platform "+
			"limit ?,?", start, size).Scan(&res).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorcode.ErrorNetworkConfigNotFound("network config not found")
		}
		return nil, err
	}
	return batchConvToBizNetworkConfig(res), nil
}

func (r *networkConfigRepo) TransferState(ctx context.Context, toAuditingNcs, toWaitOnlineNcs,
	toOnlineNcs, toExpireNcs []*biz.NetworkConfig) error {
	// 开启事务
	tx := r.data.db.WithContext(ctx).Begin()
	// 扭转至【审核中】状态
	if len(toAuditingNcs) > 0 {
		appid, env, channel, platform, versions := getTransferStateArgs(toAuditingNcs)
		r.log.WithContext(ctx).Infof("将 %v 扭转至【审核中】状态 \n\n", versions)
		if err := tx.Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
			" and `version` in (?) ", appid, env, channel, platform, versions).
			Updates(map[string]interface{}{
				"state": biz.NetworkConfigAuditing,
			}).Error; err != nil {
			return err
		}
	}

	// 扭转至【待上线】状态
	if len(toWaitOnlineNcs) > 0 {
		appid, env, channel, platform, versions := getTransferStateArgs(toWaitOnlineNcs)
		r.log.WithContext(ctx).Infof("将 %v 扭转至【待上线】状态 \n\n", versions)
		if err := tx.Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
			" and `version` in (?) ", appid, env, channel, platform, versions).
			Updates(map[string]interface{}{
				"state": biz.NetworkConfigWaitOnline,
			}).Error; err != nil {
			return err
		}
	}

	// 扭转至【已上线】状态
	if len(toOnlineNcs) > 0 {
		appid, env, channel, platform, versions := getTransferStateArgs(toOnlineNcs)
		r.log.WithContext(ctx).Infof("将 %v 扭转至【已上线】状态 \n\n", versions)
		// 先将之前【已上线】的network_config修改为【已过期】
		if err := tx.Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
			" and `state`=? ", appid, env, channel, platform, biz.NetworkConfigOnline).
			Updates(map[string]interface{}{
				"state": biz.NetworkConfigExpire,
			}).Error; err != nil {
			return err
		}

		// 再扭转至【已上线】状态
		if err := tx.Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
			" and `version`=? ", appid, env, channel, platform, versions[0]).
			Updates(map[string]interface{}{
				"state": biz.NetworkConfigOnline,
			}).Error; err != nil {
			return err
		}
	}

	// 扭转至【已过期】状态
	if len(toExpireNcs) > 0 {
		appid, env, channel, platform, versions := getTransferStateArgs(toExpireNcs)
		r.log.WithContext(ctx).Infof("将 %v 扭转至【已过期】状态 \n\n", versions)
		if err := tx.Model(&NetworkConfig{}).Where("`appid`=? and `env`=? and `channel`=? and `platform`=? "+
			" and `version` in (?) ", appid, env, channel, platform, versions).
			Updates(map[string]interface{}{
				"state": biz.NetworkConfigExpire,
			}).Error; err != nil {
			return err
		}
	}
	tx.Commit()
	return nil
}

func genNetworkConfigOnlineKey(appid uint32, env, channel, platform string) string {
	return fmt.Sprintf("game_config_network_config_online_%d_%s_%s_%s",
		appid, env, channel, platform)
}

func genNetworkConfigAuditKey(appid uint32, env, channel, platform string) string {
	return fmt.Sprintf("game_config_network_config_audit_%d_%s_%s_%s",
		appid, env, channel, platform)
}

func getTransferStateArgs(ncs []*biz.NetworkConfig) (uint32, string, string, string, []string) {
	appid := uint32(0)
	env := ""
	channel := ""
	platform := ""
	versions := make([]string, 0, len(ncs))
	if len(ncs) > 0 {
		appid = ncs[0].Appid
		env = ncs[0].Env
		channel = ncs[0].Channel
		platform = ncs[0].Platform
	}
	for _, nc := range ncs {
		versions = append(versions, nc.Version)
	}
	return appid, env, channel, platform, versions
}

// cache
func (r *networkConfigRepo) GetNetworkConfigOnlineCache(ctx context.Context, appid uint32, env, channel, platform string) (string, string, error) {
	key := genNetworkConfigOnlineKey(appid, env, channel, platform)
	datas, err := r.data.rdb.HMGet(ctx, key, "version", "online_config").Result()
	if err != nil {
		if err == redis.Nil {
			return "", "{}", errorcode.ErrorNetworkConfigNotFound("network config not found")
		}
		return "", "{}", err
	}
	if len(datas) != 2 {
		return "", "{}", errorcode.ErrorNetworkConfigNotFound("network config not found")
	}
	if datas[0] == nil || datas[1] == nil {
		return "", "{}", errorcode.ErrorNetworkConfigNotFound("network config not found")
	}
	version, _ := datas[0].(string)
	onlineConfig, _ := datas[1].(string)
	if version == "" || onlineConfig == "" {
		return "", "{}", errorcode.ErrorNetworkConfigNotFound("network config not found")
	}
	return version, onlineConfig, nil
}

func (r *networkConfigRepo) SetNetworkConfigOnlineCache(ctx context.Context, nc *biz.NetworkConfig) error {
	key := genNetworkConfigOnlineKey(nc.Appid, nc.Env, nc.Channel, nc.Platform)
	if _, err := r.data.rdb.HSet(ctx, key, "version", nc.Version, "online_config", nc.OnlineConfig).Result(); err != nil {
		return err
	}
	return nil
}

func (r *networkConfigRepo) GetNetworkConfigAuditCache(ctx context.Context, nc *biz.NetworkConfig) (*biz.NetworkConfig, error) {
	key := genNetworkConfigAuditKey(nc.Appid, nc.Env, nc.Channel, nc.Platform)
	data, err := r.data.rdb.HGet(ctx, key, nc.Version).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errorcode.ErrorNetworkConfigNotFound("network config not found")
		}
		return nil, err
	}
	if data == "" {
		return nil, errorcode.ErrorNetworkConfigNotFound("network config not found")
	}
	nc.AuditConfig = data
	return nc, nil
}

func (r *networkConfigRepo) GetAllNetworkConfigAuditCache(ctx context.Context, nc *biz.NetworkConfig) ([]*biz.NetworkConfig, error) {
	key := genNetworkConfigAuditKey(nc.Appid, nc.Env, nc.Channel, nc.Platform)
	datas, err := r.data.rdb.HGetAll(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, errorcode.ErrorNetworkConfigNotFound("network config not found")
		}
		return nil, err
	}
	ncs := make([]*biz.NetworkConfig, 0, len(datas)/2)
	for k, v := range datas {
		ncs = append(ncs, &biz.NetworkConfig{
			Appid:       nc.Appid,
			Env:         nc.Env,
			Channel:     nc.Channel,
			Platform:    nc.Platform,
			Version:     k,
			AuditConfig: v,
		})
	}
	return ncs, nil
}

func (r *networkConfigRepo) SetNetworkConfigAuditCache(ctx context.Context, nc *biz.NetworkConfig) error {
	key := genNetworkConfigAuditKey(nc.Appid, nc.Env, nc.Channel, nc.Platform)
	if _, err := r.data.rdb.HSet(ctx, key, nc.Version, nc.AuditConfig).Result(); err != nil {
		return err
	}
	return nil
}

func (r *networkConfigRepo) RemoveNetworkConfigAuditCache(ctx context.Context, nc *biz.NetworkConfig) error {
	key := genNetworkConfigAuditKey(nc.Appid, nc.Env, nc.Channel, nc.Platform)
	if _, err := r.data.rdb.HDel(ctx, key, nc.Version).Result(); err != nil {
		return err
	}
	return nil
}

func (r *networkConfigRepo) Lock(ctx context.Context, key string, expireSecond uint64) error {
	res, err := r.data.rdb.SetNX(ctx, key, "lock", time.Second*time.Duration(expireSecond)).Result()
	if err != nil {
		return err
	}
	if !res {
		return errors.New("require lock fail")
	}
	return nil
}

func (r *networkConfigRepo) Unlock(ctx context.Context, key string) error {
	_, err := r.data.rdb.Del(ctx, key).Result()
	return err
}

func (r *networkConfigRepo) ResetNetworkConfigCache(ctx context.Context, category *biz.NetworkConfig, resetOnline bool, onlineNc *biz.NetworkConfig,
	resetAuditNcs []*biz.NetworkConfig, delAuditVers []interface{}) error {
	onlineKey := genNetworkConfigOnlineKey(category.Appid, category.Env, category.Channel, category.Platform)
	onlineValues := make([]interface{}, 0)
	if resetOnline {
		onlineValues = []interface{}{"version", onlineNc.Version, "online_config", onlineNc.OnlineConfig}
	}

	auditKey := genNetworkConfigAuditKey(category.Appid, category.Env, category.Channel, category.Platform)
	auditValues := make([]interface{}, 0)
	if len(resetAuditNcs) > 0 {
		for _, nc := range resetAuditNcs {
			auditValues = append(auditValues, nc.Version)
			auditValues = append(auditValues, nc.AuditConfig)
		}
	}

	delAuditKey := genNetworkConfigAuditKey(category.Appid, category.Env, category.Channel, category.Platform)
	delAuditValues := make([]interface{}, 0)
	if len(delAuditVers) > 0 {
		delAuditValues = delAuditVers
	}

	if resetOnline {
		r.log.WithContext(ctx).Infof("需要在redis中设置【在线】network_config， key = %s, values = %+v", onlineKey, onlineValues)
	}

	if len(resetAuditNcs) > 0 {
		r.log.WithContext(ctx).Infof("需要在redis中设置【审核中】network_config， key = %s, values = %+v", auditKey, auditValues)
	}

	if len(delAuditVers) > 0 {
		r.log.WithContext(ctx).Infof("需要在redis中删除多余的【审核中】network_config， key = %s, values = %+v", delAuditKey, delAuditValues)
	}

	script := `
		local onlineLen = tonumber(ARGV[1])
		local onlineStart = 2
		local onlineEnd = onlineStart+onlineLen-1
		local result = ""
		if onlineLen > 0 then
			local onlineArgs = {}
			for i = onlineStart, onlineEnd do 
				table.insert(onlineArgs, ARGV[i])
			end
			result = redis.call('hmset',KEYS[1],unpack(onlineArgs))
		end

		local auditLen = tonumber(ARGV[onlineEnd+1])
		local auditStart = onlineEnd+2
		local auditEnd = auditStart+auditLen-1
		if auditLen > 0 then
			local auditArgs = {}
			for i = auditStart, auditEnd do 
				table.insert(auditArgs, ARGV[i])
			end
			result = redis.call('hmset',KEYS[2],unpack(auditArgs))
		end

		local delAuditLen = tonumber(ARGV[auditEnd+1])
		local delAuditStart = auditEnd+2
		local delAuditEnd = delAuditStart+delAuditLen-1
		if delAuditLen > 0 then
			local delAuditArgs = {}
			for i = delAuditStart, delAuditEnd do 
				table.insert(delAuditArgs, ARGV[i])
			end
			result = redis.call('hdel',KEYS[3],unpack(delAuditArgs))
		end
		return result
	`
	keys := []string{onlineKey, auditKey, delAuditKey}
	args := make([]interface{}, 0)
	args = append(args, len(onlineValues))
	args = append(args, onlineValues...)
	args = append(args, len(auditValues))
	args = append(args, auditValues...)
	args = append(args, len(delAuditValues))
	args = append(args, delAuditVers...)

	fmt.Printf("\n\n args = %+v \n\n", args)

	_, err := r.data.rdb.Eval(ctx, script, keys, args).Result()
	return err
}
