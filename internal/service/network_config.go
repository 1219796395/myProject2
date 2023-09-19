package service

import (
	"context"
	"time"

	"github.com/1219796395/myProject2/api/errorcode"
	pb "github.com/1219796395/myProject2/api/networkconfig"
	"github.com/1219796395/myProject2/internal/biz"
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

type NetworkConfigService struct {
	pb.UnimplementedNetworkConfigServer
	logic *biz.NetworkConfigLogic
	log   *log.Helper
	bc    *conf.Bootstrap
}

// 新建一个network_config Service
func NewNetworkConfigService(logic *biz.NetworkConfigLogic, logger log.Logger, bc *conf.Bootstrap) *NetworkConfigService {
	return &NetworkConfigService{logic: logic, log: log.NewHelper(logger), bc: bc}
}

func convToPBNetworkConfig(nc *biz.NetworkConfig) *pb.NetworkConfigDetail {
	return &pb.NetworkConfigDetail{
		AppId:        nc.Appid,
		Env:          nc.Env,
		Channel:      nc.Channel,
		Platform:     nc.Platform,
		Version:      nc.Version,
		State:        nc.State,
		AuditStart:   nc.ModifyData.AuditStart,
		AuditEnd:     nc.ModifyData.AuditEnd,
		OnlineStart:  nc.ModifyData.OnlineStart,
		OnlineConfig: nc.ModifyData.OnlineConfig,
		AuditConfig:  nc.ModifyData.AuditConfig,
		Operator:     nc.Operator,
		UpdateTime:   uint64(nc.Mtime.UnixMilli()),
	}
}

func batchConvToPBNetworkConfig(ncs []*biz.NetworkConfig) []*pb.NetworkConfigDetail {
	ncPBs := make([]*pb.NetworkConfigDetail, 0, len(ncs))
	for _, nc := range ncs {
		ncPBs = append(ncPBs, convToPBNetworkConfig(nc))
	}
	return ncPBs
}

// GetNetworkConfigList 获取network_config 列表
func (s *NetworkConfigService) GetNetworkConfigList(ctx context.Context, req *pb.GetNetworkConfigListReq) (
	*pb.GetNetworkConfigListRsp, error) {
	if !checkInnerEnv(req.Env) || !checkBatchReadPlatform(req.Platform) || !checkBatchReadChannel(req.Channel) {
		s.log.WithContext(ctx).Errorf("[GetNetworkConfigList] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	nc := &biz.NetworkConfig{
		Appid:    req.AppId,
		Env:      req.Env,
		Channel:  req.Channel,
		Platform: req.Platform,
	}
	ncs, err := s.logic.GetNetworkConfigList(ctx, nc)
	if err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[GetNetworkConfigList] fail! req = %+v, err = %+v", req, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[GetNetworkConfigList] fail! req = %+v, err = %+v", req, err)
		return nil, errorcode.ErrorInternalServerError("get network config list fail")
	}

	return &pb.GetNetworkConfigListRsp{
		NetworkConfigs: batchConvToPBNetworkConfig(ncs),
	}, nil
}

// CreateNetworkConfig 创建network_config
func (s *NetworkConfigService) CreateNetworkConfig(ctx context.Context, req *pb.CreateNetworkConfigReq) (*pb.CreateNetworkConfigRsp, error) {
	if !checkInnerEnv(req.Env) || !checkPlatform(req.Platform, s.bc) || !checkChannel(req.Channel, s.bc) ||
		!checkTime(req.AuditStart, req.AuditEnd, req.OnlineStart) || !checkConfData(req.AuditConfig) || !checkConfData(req.OnlineConfig) {
		s.log.WithContext(ctx).Errorf("[CreateNetworkConfig] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	nc := &biz.NetworkConfig{
		Appid:        req.AppId,
		Env:          req.Env,
		Channel:      req.Channel,
		Platform:     req.Platform,
		MajorVersion: req.MajorVersion,
		MinorVersion: req.MinorVersion,
		AuditStart:   req.AuditStart,
		AuditEnd:     req.AuditEnd,
		OnlineStart:  req.OnlineStart,
		OnlineConfig: req.OnlineConfig,
		AuditConfig:  req.AuditConfig,
		Operator:     req.Operator,
	}
	if err := s.logic.CreateNetworkConfig(ctx, nc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[CreateNetworkConfig] fail! req = %+v, err = %+v", req, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[CreateNetworkConfig] fail! req = %+v, err = %+v", req, err)
		return nil, errorcode.ErrorInternalServerError("create network config fail")

	}
	return &pb.CreateNetworkConfigRsp{
		NetworkConfig: convToPBNetworkConfig(nc),
	}, nil
}

func (s *NetworkConfigService) DeleteNetworkConfig(ctx context.Context, req *pb.DeleteNetworkConfigReq) (*pb.DeleteNetworkConfigRsp, error) {
	if !checkInnerEnv(req.Env) || !checkPlatform(req.Platform, s.bc) || !checkChannel(req.Channel, s.bc) {
		s.log.WithContext(ctx).Errorf("[DeleteNetworkConfig] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	nc := &biz.NetworkConfig{
		Appid:    req.AppId,
		Env:      req.Env,
		Channel:  req.Channel,
		Platform: req.Platform,
		Version:  req.Version,
		Operator: req.Operator,
		Mtime:    time.Unix(0, int64(req.GetUpdateTime())*1e6),
	}
	if err := s.logic.DeleteNetworkConfig(ctx, nc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[DeleteNetworkConfig] fail! req = %+v, err = %+v", req, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[DeleteNetworkConfig] fail! req = %+v, err = %+v", req, err)
		return nil, errorcode.ErrorInternalServerError("delete network config fail")
	}
	return &pb.DeleteNetworkConfigRsp{}, nil
}

func (s *NetworkConfigService) UpdateNetworkConfig(ctx context.Context, req *pb.UpdateNetworkConfigReq) (*pb.UpdateNetworkConfigRsp, error) {
	if !checkInnerEnv(req.Env) || !checkPlatform(req.Platform, s.bc) || !checkChannel(req.Channel, s.bc) ||
		!checkTime(req.AuditStart, req.AuditEnd, req.OnlineStart) || !checkConfData(req.AuditConfig) || !checkConfData(req.OnlineConfig) {
		s.log.WithContext(ctx).Errorf("[UpdateNetworkConfig] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}

	nc := &biz.NetworkConfig{
		Appid:        req.AppId,
		Env:          req.Env,
		Channel:      req.Channel,
		Platform:     req.Platform,
		Version:      req.Version,
		AuditStart:   req.AuditStart,
		AuditEnd:     req.AuditEnd,
		OnlineStart:  req.OnlineStart,
		OnlineConfig: req.OnlineConfig,
		AuditConfig:  req.AuditConfig,
		Operator:     req.Operator,
		Mtime:        time.Unix(0, int64(req.GetUpdateTime())*1e6),
	}
	if err := s.logic.UpdateNetworkConfig(ctx, nc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[UpdateNetworkConfig] fail! req = %+v, err = %+v", req, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[UpdateNetworkConfig] fail! req = %+v, err = %+v", req, err)
		return nil, errorcode.ErrorInternalServerError("update network config fail")
	}
	return &pb.UpdateNetworkConfigRsp{}, nil
}

func (s *NetworkConfigService) PublishNetworkConfig(ctx context.Context, req *pb.PublishNetworkConfigReq) (*pb.PublishNetworkConfigRsp, error) {
	if !checkInnerEnv(req.Env) || !checkPlatform(req.Platform, s.bc) || !checkChannel(req.Channel, s.bc) {
		s.log.WithContext(ctx).Errorf("[PublishNetworkConfig] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}

	nc := &biz.NetworkConfig{
		Appid:    req.AppId,
		Env:      req.Env,
		Channel:  req.Channel,
		Platform: req.Platform,
		Version:  req.Version,
		Operator: req.Operator,
		Mtime:    time.Unix(0, int64(req.GetUpdateTime())*1e6),
	}
	if err := s.logic.PublishNetworkConfig(ctx, nc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[PublishNetworkConfig] fail! req = %+v, err = %+v", req, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[PublishNetworkConfig] fail! req = %+v, err = %+v", req, err)
		return nil, errorcode.ErrorInternalServerError("publish network config fail")
	}
	return &pb.PublishNetworkConfigRsp{}, nil
}

func (s *NetworkConfigService) CancelPublishNetworkConfig(ctx context.Context, req *pb.CancelPublishNetworkConfigReq) (*pb.CancelPublishNetworkConfigRsp, error) {
	if !checkInnerEnv(req.Env) || !checkPlatform(req.Platform, s.bc) || !checkChannel(req.Channel, s.bc) {
		s.log.WithContext(ctx).Errorf("[CancelPublishNetworkConfig] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	nc := &biz.NetworkConfig{
		Appid:    req.AppId,
		Env:      req.Env,
		Channel:  req.Channel,
		Platform: req.Platform,
		Version:  req.Version,
		Operator: req.Operator,
		Mtime:    time.Unix(0, int64(req.GetUpdateTime())*1e6),
	}
	if err := s.logic.CancelPublishNetworkConfig(ctx, nc); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[CancelPublishNetworkConfig] fail! req = %+v, err = %+v", req, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[CancelPublishNetworkConfig] fail! req = %+v, err = %+v", req, err)
		return nil, errorcode.ErrorInternalServerError("cancel publish network config fail")
	}
	return &pb.CancelPublishNetworkConfigRsp{}, nil
}

func checkTime(auditStart, auditEnd, onlineStart uint64) bool {
	return onlineStart >= auditEnd && auditEnd >= auditStart
}

func (s *NetworkConfigService) GetNetworkConfig(ctx context.Context, req *pb.GetNetworkConfigReq) (*pb.GetNetworkConfigRsp, error) {
	if !checkChannel(req.Channel, s.bc) || !checkPlatform(req.Platform, s.bc) || !checkVersion(req.Version) {
		s.log.WithContext(ctx).Errorf("[GetNetworkConfig] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}

	// 获取http请求头中的x-hg-internal字段，若值为1，则表示内网访问；否则是外网访问
	isInner := false
	if ts, ok := transport.FromServerContext(ctx); ok && ts.Kind() == transport.KindHTTP {
		isInner = ts.RequestHeader().Get("x-hg-internal") == "1" // true表示内网访问
	}

	if !isInner && !checkOuterEnv(req.Env, s.bc) {
		s.log.WithContext(ctx).Errorf("[GetNetworkConfig] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	if isInner && !checkInnerEnv(req.Env) {
		s.log.WithContext(ctx).Errorf("[GetNetworkConfig] invalid param! req = %+v", req)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}

	nc := &biz.NetworkConfig{
		Appid:    req.AppId,
		Env:      req.Env,
		Channel:  req.Channel,
		Platform: req.Platform,
		Version:  req.Version,
	}
	state, data, err := s.logic.GetNetworkConfig(ctx, nc)
	if err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[GetRemoteConfig] fail! req = %+v, err = %+v", req, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[GetRemoteConfig] fail! req = %+v, err = %+v", req, err)
		return nil, errorcode.ErrorInternalServerError("get remote config fail")
	}
	return &pb.GetNetworkConfigRsp{
		VersionState: state,
		ConfigData:   data,
	}, nil
}

func checkVersion(version string) bool {
	return biz.VersionRegx.MatchString(version)
}
