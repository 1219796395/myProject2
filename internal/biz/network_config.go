package biz

import (
	"context"
	"fmt"
	"game-config/api/errorcode"
	"game-config/internal/conf"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

// NetworkConfigLogic network_config-逻辑层
type NetworkConfigLogic struct {
	repo    NetworkConfigRepo
	envRepo EnvManageRepo
	logRepo NetworkConfigLogRepo
	log     *log.Helper
}

var VersionRegx *regexp.Regexp

// NewNetworkConfigLogic 新建一个 network_config-逻辑层 对象
func NewNetworkConfigLogic(repo NetworkConfigRepo, envRepo EnvManageRepo, logRepo NetworkConfigLogRepo, logger log.Logger, bc *conf.Bootstrap) *NetworkConfigLogic {
	var err error
	VersionRegx, err = regexp.Compile(bc.GetBiz().GetVersionPattern())
	if err != nil {
		panic(err)
	}
	return &NetworkConfigLogic{repo: repo, envRepo: envRepo, logRepo: logRepo, log: log.NewHelper(logger)}
}

// NetworkConfig network_config-逻辑层实体类
type NetworkConfig struct {
	Appid        uint32              `json:"appid"`
	Env          string              `json:"env"`
	Channel      string              `json:"channel"`
	Platform     string              `json:"platform"`
	Version      string              `json:"version"`
	MajorVersion uint32              `json:"major_version"`
	MinorVersion uint32              `json:"minor_version"`
	State        uint32              `json:"state"`
	AuditStart   uint64              `json:"audit_start"`
	AuditEnd     uint64              `json:"audit_end"`
	OnlineStart  uint64              `json:"online_start"`
	OnlineConfig string              `json:"online_config"`
	AuditConfig  string              `json:"audit_config"`
	ModifyData   NetworkConfigModify `json:"modify_data"`
	Operator     string              `json:"operator"`
	Ctime        time.Time           `json:"ctime"`
	Mtime        time.Time           `json:"mtime"`
}

type NetworkConfigModify struct {
	AuditStart   uint64 `json:"audit_start"`
	AuditEnd     uint64 `json:"audit_end"`
	OnlineStart  uint64 `json:"online_start"`
	OnlineConfig string `json:"online_config"`
	AuditConfig  string `json:"audit_config"`
}

// network_config状态
const (
	NetworkConfigNeverPublished = 1 // 未发布
	NetworkConfigWaitPublish    = 2 // 待发布
	NetworkConfigWaitAudit      = 3 // 待审核
	NetworkConfigAuditing       = 4 // 审核中
	NetworkConfigWaitOnline     = 5 // 待上线
	NetworkConfigOnline         = 6 // 已上线
	NetworkConfigExpire         = 7 // 已过期
	NetworkConfigDeleted        = 8 // 已删除
)

// NetworkConfigRepo network_config-数据层
type NetworkConfigRepo interface {
	// db
	GetNetworkConfig(context.Context, *NetworkConfig) (*NetworkConfig, error)
	GetNetworkConfigList(context.Context, *NetworkConfig) ([]*NetworkConfig, error)
	GetNetworkConfigListWithStates(ctx context.Context, nc *NetworkConfig, states []uint32) ([]*NetworkConfig, error)
	CreateNetworkConfig(context.Context, *NetworkConfig) error
	UpdateNetworkConfig(context.Context, *NetworkConfig, *NetworkConfig) error
	DeleteNetworkConfig(context.Context, *NetworkConfig, *NetworkConfig) error
	PublishNetworkConfig(context.Context, *NetworkConfig, *NetworkConfig) error
	CancelPublishNetworkConfig(context.Context, *NetworkConfig, *NetworkConfig) error
	PageGetCategoryList(ctx context.Context, start, size int) ([]*NetworkConfig, error)
	TransferState(ctx context.Context, toAuditingNcs, toWaitOnlineNcs, toOnlineNcs, toExpireNcs []*NetworkConfig) error

	// cache
	GetNetworkConfigOnlineCache(ctx context.Context, appid uint32, env, channel, platform string) (string, string, error)
	SetNetworkConfigOnlineCache(context.Context, *NetworkConfig) error
	GetNetworkConfigAuditCache(context.Context, *NetworkConfig) (*NetworkConfig, error)
	GetAllNetworkConfigAuditCache(ctx context.Context, nc *NetworkConfig) ([]*NetworkConfig, error)
	SetNetworkConfigAuditCache(context.Context, *NetworkConfig) error
	RemoveNetworkConfigAuditCache(context.Context, *NetworkConfig) error
	ResetNetworkConfigCache(context.Context, *NetworkConfig, bool, *NetworkConfig, []*NetworkConfig, []interface{}) error
	Unlock(ctx context.Context, key string) error
	Lock(ctx context.Context, key string, expireSecond uint64) error
}

// GetNetworkConfigList 获取network_config列表
func (l *NetworkConfigLogic) GetNetworkConfigList(ctx context.Context, nc *NetworkConfig) (
	[]*NetworkConfig, error) {
	// 校验是否存在该环境名
	if err := existEnv(ctx, l.envRepo, nc.Appid, nc.Env); err != nil {
		return nil, err
	}

	return l.repo.GetNetworkConfigList(ctx, nc)
}

// CreateNetworkConfig 创建network_config配置
func (l *NetworkConfigLogic) CreateNetworkConfig(ctx context.Context, nc *NetworkConfig) error {
	// 校验是否存在该环境名
	if err := existEnv(ctx, l.envRepo, nc.Appid, nc.Env); err != nil {
		return err
	}

	// 创建network_config
	randStr := randString(10)
	nc.Version = fmt.Sprintf("%d.%d-%s", nc.MajorVersion, nc.MinorVersion, randStr)
	if err := l.repo.CreateNetworkConfig(ctx, nc); err != nil {
		return err
	}

	// 记录配置日志
	log := convNetworkConfigToLog(ctx, operationCreate, &NetworkConfig{}, nc)
	return l.logRepo.CreateNetworkConfigLog(ctx, log)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// DeleteNetworkConfig 删除network_config配置
func (l *NetworkConfigLogic) DeleteNetworkConfig(ctx context.Context, nc *NetworkConfig) error {
	beforeNc, err := l.repo.GetNetworkConfig(ctx, nc)
	if err != nil {
		return err
	}

	// 未发布状态可用
	if beforeNc.State != NetworkConfigNeverPublished {
		return errorcode.ErrorBadRequest("not allow delete")
	}

	if nc.Mtime.UnixMilli() != 0 && nc.Mtime.UnixMilli() < beforeNc.Mtime.UnixMilli() {
		return errorcode.ErrorNetworkConfigExpireRequest("expire request")
	}

	// 删除network_config
	if err := l.repo.DeleteNetworkConfig(ctx, beforeNc, nc); err != nil {
		return err
	}

	// 记录日志
	log := convNetworkConfigToLog(ctx, operationDelete, beforeNc, nc)
	return l.logRepo.CreateNetworkConfigLog(ctx, log)
}

// UpdateNetworkConfig 修改network_config配置
func (l *NetworkConfigLogic) UpdateNetworkConfig(ctx context.Context, nc *NetworkConfig) error {
	beforeNc, err := l.repo.GetNetworkConfig(ctx, nc)
	if err != nil {
		return err
	}
	if nc.Mtime.UnixMilli() != 0 && nc.Mtime.UnixMilli() < beforeNc.Mtime.UnixMilli() {
		return errorcode.ErrorNetworkConfigExpireRequest("expire request")
	}

	// 修改network_config
	if err := l.repo.UpdateNetworkConfig(ctx, beforeNc, nc); err != nil {
		return err
	}

	// 记录日志
	log := convNetworkConfigToLog(ctx, operationUpdate, beforeNc, nc)
	return l.logRepo.CreateNetworkConfigLog(ctx, log)
}

// PublishNetworkConfig 发布network_config
func (l *NetworkConfigLogic) PublishNetworkConfig(ctx context.Context, nc *NetworkConfig) error {
	beforeNc, err := l.repo.GetNetworkConfig(ctx, nc)
	if err != nil {
		return err
	}

	// 未发布、待发布状态可用
	if beforeNc.State != NetworkConfigNeverPublished && beforeNc.State != NetworkConfigWaitPublish {
		return errorcode.ErrorBadRequest("not allow publish")
	}

	if nc.Mtime.UnixMilli() != 0 && nc.Mtime.UnixMilli() < beforeNc.Mtime.UnixMilli() {
		return errorcode.ErrorNetworkConfigExpireRequest("expire request")
	}

	// 发布需要写个事务，若发布的配置是上线，需要把另一个上线的配置改为【已过期】
	if err := l.repo.PublishNetworkConfig(ctx, beforeNc, nc); err != nil {
		return err
	}

	if err := l.transferCacheState(ctx, nc); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	// 记录日志
	log := convNetworkConfigToLog(ctx, operationPublish, beforeNc, nc)
	return l.logRepo.CreateNetworkConfigLog(ctx, log)
}

// CancelPublishNetworkConfig 取消发布network_config
func (l *NetworkConfigLogic) CancelPublishNetworkConfig(ctx context.Context, nc *NetworkConfig) error {
	beforeNc, err := l.repo.GetNetworkConfig(ctx, nc)
	if err != nil {
		return err
	}
	// 发布上线后，与上线前可用
	if beforeNc.State == NetworkConfigOnline || beforeNc.State == NetworkConfigNeverPublished ||
		beforeNc.State == NetworkConfigExpire {
		return errorcode.ErrorBadRequest("not allow cancel publish")
	}

	if nc.Mtime.UnixMilli() != 0 && nc.Mtime.UnixMilli() < beforeNc.Mtime.UnixMilli() {
		return errorcode.ErrorNetworkConfigExpireRequest("expire request")
	}

	// 取消发布配置
	if err := l.repo.CancelPublishNetworkConfig(ctx, beforeNc, nc); err != nil {
		return err
	}

	if err := l.transferCacheState(ctx, nc); err != nil {
		return err
	}

	if err != nil {
		return err
	}
	// 记录日志
	log := convNetworkConfigToLog(ctx, operationCancelPublish, beforeNc, nc)
	return l.logRepo.CreateNetworkConfigLog(ctx, log)
}

const (
	VersionOnline    = 1
	VersionExpire    = 2
	VersionNotOnline = 3
)

// GetNetworkConfig C端获取单个network_config
func (l *NetworkConfigLogic) GetNetworkConfig(ctx context.Context, nc *NetworkConfig) (uint32, string, error) {
	onlineVer, onlineConfig, err := l.repo.GetNetworkConfigOnlineCache(ctx, nc.Appid, nc.Env, nc.Channel, nc.Platform)
	if err != nil && !errorcode.IsNetworkConfigNotFound(err) {
		return 0, "{}", err
	}
	if onlineVer != "" {
		localVers := strings.Split(strings.Split(nc.Version, "-")[0], ".")
		localMajorVer, _ := strconv.Atoi(localVers[0])
		localMinorVer, _ := strconv.Atoi(localVers[1])
		onlineVers := strings.Split(strings.Split(onlineVer, "-")[0], ".")
		onlineMajorVer, _ := strconv.Atoi(onlineVers[0])
		onlineMinorVer, _ := strconv.Atoi(onlineVers[0])
		// 若本地大版本 == 在线大版本 && 本地小版本 <= 在线小版本
		if localMajorVer == onlineMajorVer && localMinorVer <= onlineMinorVer {
			return VersionOnline, onlineConfig, nil
		}
		// 若本地大版本 < 在线大版本
		if localMajorVer < onlineMajorVer {
			return VersionExpire, onlineConfig, nil
		}
	}

	auditConf, err := l.repo.GetNetworkConfigAuditCache(ctx, nc)
	// 若命中 审核中 版本
	if err == nil {
		return VersionNotOnline, auditConf.AuditConfig, nil
	}
	// 若未命中 审核中 版本
	if errorcode.IsNetworkConfigNotFound(err) {
		return VersionNotOnline, onlineConfig, nil
	}

	return 0, "{}", err
}

const networkConfigTransferStateLockKey = "network_config_tranfer_state_lock_key"

// 自动扭转状态
func (l *NetworkConfigLogic) AutoTransferState(ctx context.Context) error {
	bc, err := conf.GetConf()
	if err != nil {
		return err
	}

	lockKey := networkConfigTransferStateLockKey
	defer func() {
		// 释放分布式锁
		l.repo.Unlock(ctx, lockKey)
	}()

	// 获取分布式锁
	if err := l.repo.Lock(ctx, lockKey, bc.Biz.NetworkConfigTranferStateTask.LockExpire); err != nil {
		return err
	}

	start := 0
	size := int(bc.GetBiz().NetworkConfigTranferStateTask.GetBatch())
	now := uint64(time.Now().UnixMilli())
	wg := sync.WaitGroup{}
	for {
		// 获取 network_config 对分类列表， 分类包括 appid、env、channel、platform
		categroies, err := l.repo.PageGetCategoryList(ctx, start, size)
		if err != nil {
			return err
		}
		if len(categroies) == 0 {
			break
		}
		for _, c := range categroies {
			// 并发扭转db与cache中的状态
			category := c
			wg.Add(1)
			go func() {
				defer wg.Done()
				// 扭转db中的network_config状态
				if err := l.transferDBState(ctx, now, category); err != nil {
					log.Errorf("err = %+v", err)
					return
				}
				// 扭转cache中的network_config状态
				if err := l.transferCacheState(ctx, category); err != nil {
					log.Errorf("err = %+v", err)
				}
			}()
		}
		wg.Wait()
		start += size
	}
	return nil
}

// 扭转db中的状态
func (l *NetworkConfigLogic) transferDBState(ctx context.Context, now uint64, category *NetworkConfig) error {
	// 获取该分类下的【待审核】、【审核中】、【待上线】、【已上线】network_config列表
	dbStates := []uint32{
		NetworkConfigWaitAudit,
		NetworkConfigAuditing,
		NetworkConfigWaitOnline,
		NetworkConfigOnline,
	}
	ncs, err := l.repo.GetNetworkConfigListWithStates(ctx, category, dbStates)
	if err != nil {
		return err
	}
	if len(ncs) == 0 {
		return nil
	}
	// 将这些network_config 根据上线时间，从大到小排列
	sort.SliceStable(ncs, func(i, j int) bool {
		return ncs[i].OnlineStart >= ncs[j].OnlineStart
	})
	// 对于这些network_config，扭转状态
	toAuditingNcs := make([]*NetworkConfig, 0)
	toWaitOnlineNcs := make([]*NetworkConfig, 0)
	toExpireNcs := make([]*NetworkConfig, 0)
	toOnlineNcs := make([]*NetworkConfig, 0)
	hasOnline := false
	for _, nc := range ncs {
		state := GetPublishedNetworkConfigState(now, nc.AuditStart, nc.AuditEnd, nc.OnlineStart)
		// 只取上线时间最晚的network_config扭转为上线状态
		if state == NetworkConfigOnline {
			if !hasOnline {
				state = NetworkConfigOnline
				hasOnline = true
			} else {
				state = NetworkConfigExpire
			}
		}
		// 若状态未扭转，跳过该循环
		if nc.State == state {
			continue
		}
		nc.State = state
		if nc.State == NetworkConfigAuditing {
			toAuditingNcs = append(toAuditingNcs, nc)
		}
		if nc.State == NetworkConfigWaitOnline {
			toWaitOnlineNcs = append(toWaitOnlineNcs, nc)
		}
		if nc.State == NetworkConfigOnline {
			toOnlineNcs = append(toOnlineNcs, nc)
		}
		if nc.State == NetworkConfigExpire {
			toExpireNcs = append(toExpireNcs, nc)
		}
	}
	return l.repo.TransferState(ctx, toAuditingNcs, toWaitOnlineNcs, toOnlineNcs, toExpireNcs)
}

// 获取已发布的network_config状态
func GetPublishedNetworkConfigState(now, auditStart, auditEnd, onlineStart uint64) uint32 {
	state := uint32(0)
	// 判断network_config是否处于【待审核】状态
	if now < auditStart {
		state = NetworkConfigWaitAudit
	}
	// 判断network_config是否处于【审核中】状态
	if now >= auditStart && now < auditEnd {
		state = NetworkConfigAuditing
	}
	// 判断network_config是否处于【待上线】状态
	if now >= auditEnd && now < onlineStart {
		state = NetworkConfigWaitOnline
	}
	// 判断network_config是否处于【已上线】状态
	if now >= onlineStart {
		state = NetworkConfigOnline
	}
	return state
}

// 扭转cache中的state
func (l *NetworkConfigLogic) transferCacheState(ctx context.Context, category *NetworkConfig) error {

	// 获取该分类下的【审核中】和【已上线】的network_config列表
	cacheStates := []uint32{
		NetworkConfigAuditing,
		NetworkConfigOnline,
	}
	ncDBs, err := l.repo.GetNetworkConfigListWithStates(ctx, category, cacheStates)
	if err != nil {
		return err
	}
	var onlineNcDB *NetworkConfig
	auditNcDBs := make([]*NetworkConfig, 0)
	for _, nc := range ncDBs {
		if nc.State == NetworkConfigAuditing {
			auditNcDBs = append(auditNcDBs, nc)
		}
		if nc.State == NetworkConfigOnline {
			onlineNcDB = nc
		}
	}

	// 获取cache中的network_config信息，做对账
	onlineCacheVersion, onlineCacheConf, err :=
		l.repo.GetNetworkConfigOnlineCache(ctx, category.Appid, category.Env, category.Channel, category.Platform)
	if err != nil && !errorcode.IsNetworkConfigNotFound(err) {
		return err
	}
	auditNcCaches, err := l.repo.GetAllNetworkConfigAuditCache(ctx, category)
	if err != nil && !errorcode.IsNetworkConfigNotFound(err) {
		return err
	}

	// 对账
	auditNcCacheMap := make(map[string]string)
	for _, nc := range auditNcCaches {
		key := nc.Version
		auditNcCacheMap[key] = nc.AuditConfig
	}
	// 判断是否重置	cache中【在线】的network_config
	resetOnline := false
	if onlineNcDB != nil && onlineNcDB.Version != "" && onlineNcDB.OnlineConfig != "" {
		resetOnline = onlineNcDB.Version != onlineCacheVersion || onlineNcDB.OnlineConfig != onlineCacheConf
	}
	// 判断是否重置cache中【审核中】的network_config
	resetAuditNcs := make([]*NetworkConfig, 0)
	for _, nc := range auditNcDBs {
		if data, ok := auditNcCacheMap[nc.Version]; !ok || data != nc.AuditConfig {
			resetAuditNcs = append(resetAuditNcs, nc)
		}
	}

	if !resetOnline && len(resetAuditNcs) == 0 && len(auditNcDBs) == len(auditNcCaches) {
		l.log.Infof("appid = %d, env = %s, channel = %s, platform = %s, db与cache数据一致，无需修改cache数据", category.Appid, category.Env, category.Channel, category.Platform)
		return nil
	}

	// 获取cache中冗余的【审核中】network_config
	for _, nc := range auditNcDBs {
		delete(auditNcCacheMap, nc.Version)
	}
	delAuditVers := []interface{}{}
	for k := range auditNcCacheMap {
		delAuditVers = append(delAuditVers, k)
	}

	// 根据这些network_config信息，重置cache
	return l.repo.ResetNetworkConfigCache(ctx, category, resetOnline, onlineNcDB, resetAuditNcs, delAuditVers)
}
