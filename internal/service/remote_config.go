package service

import (
	"context"
	"time"

	"github.com/1219796395/myProject2/api/common"
	"github.com/1219796395/myProject2/api/errorcode"
	pb "github.com/1219796395/myProject2/api/remoteconfig"
	"github.com/1219796395/myProject2/internal/biz"
	"github.com/1219796395/myProject2/internal/conf"
	"github.com/1219796395/myProject2/internal/middleware"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

// RemoteConfigService is a remote config service.
type RemoteConfigService struct {
	pb.UnimplementedRemoteConfigServer

	logic *biz.RemoteConfigLogic
	log   *log.Helper
	bc    *conf.Bootstrap
}

// NewRemoteConfigService ...
func NewRemoteConfigService(logic *biz.RemoteConfigLogic, logger log.Logger, bc *conf.Bootstrap) *RemoteConfigService {
	return &RemoteConfigService{logic: logic, log: log.NewHelper(logger), bc: bc}
}

func trans2PBRemoteConfig(rc *biz.RemoteConfig) *pb.RemoteConfigDetail {
	return &pb.RemoteConfigDetail{
		AppId:           rc.Appid,
		Env:             rc.Env,
		Channel:         rc.Channel,
		Platform:        rc.Platform,
		Name:            rc.Name,
		State:           rc.State,
		Data:            rc.Data,
		ModifyData:      rc.ModifyData,
		NotDistChannel:  rc.NotDistChannel,
		NotDistPlatform: rc.NotDistPlatform,
		Operator:        rc.Operator,
		UpdateTime:      uint64(rc.Mtime.UnixMilli()),
	}
}

func batchTrans2PBRemoteConfig(rcs []*biz.RemoteConfig) []*pb.RemoteConfigDetail {
	rcPBs := make([]*pb.RemoteConfigDetail, 0, len(rcs))
	for _, rc := range rcs {
		rcPBs = append(rcPBs, trans2PBRemoteConfig(rc))
	}
	return rcPBs
}

// GetRemoteConfigList 获取远程配置列表
func (s *RemoteConfigService) GetRemoteConfigList(ctx context.Context, in *pb.GetRemoteConfigListReq) (
	*pb.GetRemoteConfigListRsp, error) {
	if !checkInnerEnv(in.Env) || !checkBatchReadPlatform(in.Platform) || !checkBatchReadChannel(in.Channel) {
		s.log.WithContext(ctx).Errorf("[GetRemoteConfigList] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	channel := in.GetChannel()
	if channel == "all" {
		channel = ""
	}
	platform := in.GetPlatform()
	if platform == "all" {
		platform = ""
	}
	rc := &biz.RemoteConfig{
		Appid:    in.GetCommon().GetAppId(),
		Env:      in.GetEnv(),
		Channel:  channel,
		Platform: platform,
		Name:     in.ConfigName,
	}
	rcs, err := s.logic.GetRemoteConfigList(ctx, rc)
	if err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[GetRemoteConfigList] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[GetRemoteConfigList] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("get remote config list fail")
	}
	rcPBs := batchTrans2PBRemoteConfig(rcs)
	return &pb.GetRemoteConfigListRsp{
		RemoteConfigs: rcPBs,
	}, nil
}

// CreateRemoteConfig 创建远程配置
func (s *RemoteConfigService) CreateRemoteConfig(ctx context.Context, in *pb.CreateRemoteConfigReq) (
	*pb.CreateRemoteConfigRsp, error) {
	if !checkInnerEnv(in.Env) || !checkPlatform(in.Platform, s.bc) || !checkChannel(in.Channel, s.bc) ||
		!checkConfName(in.ConfigName) || !checkConfData(in.ConfigData) {
		s.log.WithContext(ctx).Errorf("[CreateRemoteConfig] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	if in.GetNotDistChannel() == 1 {
		in.Channel = "default"
	}
	if in.GetNotDistPlatform() == 1 {
		in.Platform = "default"
	}
	user := middleware.GetAdminUser(ctx)
	if user.SkipAuth {
		user.Name = "game ci api"
		user.Id = 0
	}
	rc := &biz.RemoteConfig{
		Appid:           in.GetCommon().GetAppId(),
		Env:             in.GetEnv(),
		Channel:         in.GetChannel(),
		Platform:        in.GetPlatform(),
		Name:            in.GetConfigName(),
		Data:            in.GetConfigData(),
		NotDistChannel:  in.GetNotDistChannel(),
		NotDistPlatform: in.GetNotDistPlatform(),
		Operator:        user.Name,
		OperatorID:      user.Id,
	}
	if err := s.logic.CreateRemoteConfig(ctx, rc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[CreateRemoteConfig] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[CreateRemoteConfig] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("create remote config fail")
	}
	return &pb.CreateRemoteConfigRsp{}, nil
}

// CreateRemoteConfigV1 旧的创建远程配置
func (s *RemoteConfigService) CreateRemoteConfigV1(ctx context.Context, in *pb.CreateRemoteConfigV1Req) (
	*pb.CreateRemoteConfigRsp, error) {
	req := &pb.CreateRemoteConfigReq{
		Common:     &common.Common{AppId: in.GetAppid()},
		Env:        in.GetEnv(),
		Channel:    in.GetChannel(),
		Platform:   in.GetPlatform(),
		ConfigName: in.GetConfigName(),
		ConfigData: in.GetConfigData(),
	}
	return s.CreateRemoteConfig(ctx, req)
}

// DeleteRemoteConfig 删除远程配置
func (s *RemoteConfigService) DeleteRemoteConfig(ctx context.Context, in *pb.DeleteRemoteConfigReq) (
	*pb.DeleteRemoteConfigRsp, error) {
	if !checkInnerEnv(in.Env) || !checkPlatform(in.Platform, s.bc) || !checkChannel(in.Channel, s.bc) || !checkConfName(in.ConfigName) {
		s.log.WithContext(ctx).Errorf("[DeleteRemoteConfig] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	user := middleware.GetAdminUser(ctx)
	if user.SkipAuth {
		user.Name = "game ci api"
		user.Id = 0
	}
	rc := &biz.RemoteConfig{
		Appid:      in.GetCommon().GetAppId(),
		Env:        in.GetEnv(),
		Channel:    in.GetChannel(),
		Platform:   in.GetPlatform(),
		Name:       in.GetConfigName(),
		Operator:   user.Name,
		OperatorID: user.Id,
		Mtime:      time.Unix(0, int64(in.GetUpdateTime())*1e6),
	}
	if err := s.logic.DeleteRemoteConfig(ctx, rc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[DeleteRemoteConfig] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[DeleteRemoteConfig] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("delete remote config fail")
	}
	return &pb.DeleteRemoteConfigRsp{}, nil
}

// UpdateRemoteConfig 修改远程配置
func (s *RemoteConfigService) UpdateRemoteConfig(ctx context.Context, in *pb.UpdateRemoteConfigReq) (
	*pb.UpdateRemoteConfigRsp, error) {
	if !checkInnerEnv(in.Env) || !checkPlatform(in.Platform, s.bc) || !checkChannel(in.Channel, s.bc) ||
		!checkConfName(in.ConfigName) || !checkConfData(in.ConfigData) {
		s.log.WithContext(ctx).Errorf("[UpdateRemoteConfig] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	user := middleware.GetAdminUser(ctx)
	if user.SkipAuth {
		user.Name = "game ci api"
		user.Id = 0
	}
	rc := &biz.RemoteConfig{
		Appid:      in.GetCommon().GetAppId(),
		Env:        in.GetEnv(),
		Channel:    in.GetChannel(),
		Platform:   in.GetPlatform(),
		Name:       in.GetConfigName(),
		Data:       in.GetConfigData(),
		Operator:   user.Name,
		OperatorID: user.Id,
		Mtime:      time.Unix(0, int64(in.GetUpdateTime())*1e6),
	}
	if err := s.logic.UpdateRemoteConfig(ctx, rc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[UpdateRemoteConfig] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[UpdateRemoteConfig] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("update remote config fail")
	}
	return &pb.UpdateRemoteConfigRsp{}, nil
}

// UpdateRemoteConfigV1 旧的修改远程配置
func (s *RemoteConfigService) UpdateRemoteConfigV1(ctx context.Context, in *pb.UpdateRemoteConfigV1Req) (
	*pb.UpdateRemoteConfigRsp, error) {
	req := &pb.UpdateRemoteConfigReq{
		Common:     &common.Common{AppId: in.GetAppid()},
		Env:        in.GetEnv(),
		Channel:    in.GetChannel(),
		Platform:   in.GetPlatform(),
		ConfigName: in.GetConfigName(),
		ConfigData: in.GetConfigData(),
		UpdateTime: in.GetUpdateTime(),
	}
	return s.UpdateRemoteConfig(ctx, req)
}

// PublishRemoteConfig 发布远程配置
func (s *RemoteConfigService) PublishRemoteConfig(ctx context.Context, in *pb.PublishRemoteConfigReq) (
	*pb.PublishRemoteConfigRsp, error) {
	if !checkInnerEnv(in.Env) || !checkPlatform(in.Platform, s.bc) || !checkChannel(in.Channel, s.bc) || !checkConfName(in.ConfigName) {
		s.log.WithContext(ctx).Errorf("[PublishRemoteConfig] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	user := middleware.GetAdminUser(ctx)
	if user.SkipAuth {
		user.Name = "game ci api"
		user.Id = 0
	}
	rc := &biz.RemoteConfig{
		Appid:      in.GetCommon().GetAppId(),
		Env:        in.GetEnv(),
		Channel:    in.GetChannel(),
		Platform:   in.GetPlatform(),
		Name:       in.GetConfigName(),
		Operator:   user.Name,
		OperatorID: user.Id,
		Mtime:      time.Unix(0, int64(in.GetUpdateTime())*1e6),
	}
	if err := s.logic.PublishRemoteConfig(ctx, rc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[PublishRemoteConfig] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[PublishRemoteConfig] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("publish remote config fail")
	}
	return &pb.PublishRemoteConfigRsp{}, nil
}

// PublishRemoteConfigV1 旧的发布远程配置
func (s *RemoteConfigService) PublishRemoteConfigV1(ctx context.Context, in *pb.PublishRemoteConfigV1Req) (
	*pb.PublishRemoteConfigRsp, error) {
	req := &pb.PublishRemoteConfigReq{
		Common:     &common.Common{AppId: in.GetAppid()},
		Env:        in.GetEnv(),
		Channel:    in.GetChannel(),
		Platform:   in.GetPlatform(),
		ConfigName: in.GetConfigName(),
		UpdateTime: in.GetUpdateTime(),
	}
	return s.PublishRemoteConfig(ctx, req)
}

// CancelPublishRemoteConfig 取消发布远程配置
func (s *RemoteConfigService) CancelPublishRemoteConfig(ctx context.Context, in *pb.CancelPublishRemoteConfigReq) (
	*pb.CancelPublishRemoteConfigRsp, error) {
	if !checkInnerEnv(in.Env) || !checkPlatform(in.Platform, s.bc) || !checkChannel(in.Channel, s.bc) || !checkConfName(in.ConfigName) {
		s.log.WithContext(ctx).Errorf("[CancelPublishRemoteConfig] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	user := middleware.GetAdminUser(ctx)
	if user.SkipAuth {
		user.Name = "game ci api"
		user.Id = 0
	}
	rc := &biz.RemoteConfig{
		Appid:      in.GetCommon().GetAppId(),
		Env:        in.GetEnv(),
		Channel:    in.GetChannel(),
		Platform:   in.GetPlatform(),
		Name:       in.GetConfigName(),
		Operator:   user.Name,
		OperatorID: user.Id,
		Mtime:      time.Unix(0, int64(in.GetUpdateTime())*1e6),
	}
	if err := s.logic.CancelPublishRemoteConfig(ctx, rc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[CancelPublishRemoteConfig] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[CancelPublishRemoteConfig] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("cancel publish remote config fail")
	}
	return &pb.CancelPublishRemoteConfigRsp{}, nil
}

// GetRemoteConfig C端获取单个远程配置（请求头字段区分公网和内网）
func (s *RemoteConfigService) GetRemoteConfig(ctx context.Context, in *pb.GetRemoteConfigReq) (
	*pb.GetRemoteConfigRsp, error) {
	if !checkChannel(in.Channel, s.bc) || !checkPlatform(in.Platform, s.bc) {
		s.log.WithContext(ctx).Errorf("[GetRemoteConfig] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	// 获取http请求头中的x-hg-internal字段，若值为1，则表示内网访问；否则是外网访问
	isInner := false
	if ts, ok := transport.FromServerContext(ctx); ok && ts.Kind() == transport.KindHTTP {
		isInner = ts.RequestHeader().Get("x-hg-internal") == "1" // true表示内网访问
	}

	if !isInner && !checkOuterEnv(in.Env, s.bc) {
		s.log.WithContext(ctx).Errorf("[GetRemoteConfig] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	if isInner && !checkInnerEnv(in.Env) {
		s.log.WithContext(ctx).Errorf("[GetRemoteConfig] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}

	data, err := s.logic.GetRemoteConfig(ctx, in.AppId, in.Env, in.Channel, in.Platform, in.ConfigName)
	if err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[GetRemoteConfig] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[GetRemoteConfig] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("get remote config fail")
	}
	return &pb.GetRemoteConfigRsp{
		ConfigData: data,
	}, nil
}

// GetRemoteConfigV1 旧的C端获取单个远程配置（请求头字段区分公网和内网）
func (s *RemoteConfigService) GetRemoteConfigV1(ctx context.Context, in *pb.GetRemoteConfigV1Req) (
	*pb.GetRemoteConfigRsp, error) {
	req := &pb.GetRemoteConfigReq{
		AppId:      in.GetAppid(),
		Env:        in.GetEnv(),
		Channel:    in.GetChannel(),
		Platform:   in.GetPlatform(),
		ConfigName: in.GetConfigName(),
	}
	return s.GetRemoteConfig(ctx, req)
}

func checkOuterEnv(env string, bc *conf.Bootstrap) bool {
	for _, e := range bc.GetBiz().OuterEnvList {
		if env == e {
			return true
		}
	}
	return false
}

func checkInnerEnv(env string) bool {
	return env != ""
}

func checkBatchReadPlatform(platform string) bool {
	bc, err := conf.GetConf()
	if err != nil {
		return false
	}
	for _, p := range bc.GetBiz().BatchReadPlatformList {
		if platform == p {
			return true
		}
	}
	return false
}

func checkPlatform(platform string, bc *conf.Bootstrap) bool {
	for _, p := range bc.GetBiz().PlatformList {
		if platform == p {
			return true
		}
	}
	return false
}

func checkBatchReadChannel(channel string) bool {
	bc, err := conf.GetConf()
	if err != nil {
		return false
	}
	for _, c := range bc.GetBiz().BatchReadChannelList {
		if channel == c {
			return true
		}
	}
	return false
}

func checkChannel(channel string, bc *conf.Bootstrap) bool {
	for _, c := range bc.GetBiz().ChannelList {
		if channel == c {
			return true
		}
	}
	return false
}

func checkConfName(confName string) bool {
	return confName != ""
}

func checkConfData(confData string) bool {
	return json.Valid([]byte(confData))
}
