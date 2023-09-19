package service

import (
	"context"
	"game-config/api/errorcode"
	pb "game-config/api/projectconfig/envmanage"
	"game-config/internal/biz"
	"game-config/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// EnvManageService is a remote config service.
type EnvManageService struct {
	pb.UnimplementedEnvManageServer

	logic *biz.EnvManageLogic
	log   *log.Helper
}

// NewEnvManageService ...
func NewEnvManageService(logic *biz.EnvManageLogic, logger log.Logger) *EnvManageService {
	return &EnvManageService{logic: logic, log: log.NewHelper(logger)}
}

func (s *EnvManageService) GetEnvList(ctx context.Context, in *pb.GetEnvListReq) (
	*pb.GetEnvListRsp, error) {
	env := &biz.Env{
		Appid: in.GetCommon().GetAppId(),
	}
	envs, err := s.logic.GetEnvList(ctx, env)
	if err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[GetEnvList] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[GetEnvList] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("get env list fail")
	}
	envPBs := make([]*pb.Env, 0)
	for _, env := range envs {
		envPBs = append(envPBs, &pb.Env{
			AppId:      env.Appid,
			Field:      env.Field,
			Name:       env.Name,
			Remark:     env.Remark,
			IsPreset:   env.IsPreset,
			Operator:   env.Operator,
			UpdateTime: uint64(env.Mtime.UnixMilli()),
		})
	}
	return &pb.GetEnvListRsp{Envs: envPBs}, nil
}

func (s *EnvManageService) CreateEnv(ctx context.Context, in *pb.CreateEnvReq) (
	*pb.CreateEnvRsp, error) {
	if !checkEnvField(in.EnvField) || !checkEnvName(in.EnvName) {
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	env := &biz.Env{
		Appid:  in.GetCommon().GetAppId(),
		Field:  in.EnvField,
		Name:   in.EnvName,
		Remark: in.Remark,
	}
	if err := s.logic.CreateEnv(ctx, env); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[CreateEnv] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[CreateEnv] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("create env fail")
	}
	return &pb.CreateEnvRsp{}, nil
}

func (s *EnvManageService) UpdateEnv(ctx context.Context, in *pb.UpdateEnvReq) (
	*pb.UpdateEnvRsp, error) {
	if !checkEnvField(in.EnvField) || !checkEnvName(in.EnvName) || !checkEnvName(in.NewEnvName) {
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	env := &biz.Env{
		Appid:   in.GetCommon().GetAppId(),
		Field:   in.EnvField,
		Name:    in.EnvName,
		NewName: in.NewEnvName,
		Remark:  in.Remark,
		Mtime:   time.Unix(0, int64(in.UpdateTime)*1e6),
	}
	if err := s.logic.UpdateEnv(ctx, env); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[UpdateEnv] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[UpdateEnv] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("update env fail")
	}
	return &pb.UpdateEnvRsp{}, nil
}

func (s *EnvManageService) DeleteEnv(ctx context.Context, in *pb.DeleteEnvReq) (
	*pb.DeleteEnvRsp, error) {
	if !checkEnvField(in.EnvField) || !checkEnvName(in.EnvName) {
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	env := &biz.Env{
		Appid: in.GetCommon().GetAppId(),
		Field: in.EnvField,
		Name:  in.EnvName,
		Mtime: time.Unix(0, int64(in.UpdateTime)*1e6),
	}
	if err := s.logic.DeleteEnv(ctx, env); err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[DeleteEnv] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[DeleteEnv] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("delete env fail")
	}
	return &pb.DeleteEnvRsp{}, nil
}

func checkEnvField(envField string) bool {
	bc, err := conf.GetConf()
	if err != nil {
		return false
	}
	for _, p := range bc.GetBiz().EnvFieldList {
		if envField == p {
			return true
		}
	}
	return false
}

func checkEnvName(envName string) bool {
	bc, err := conf.GetConf()
	if err != nil {
		return false
	}
	for _, p := range bc.GetBiz().EnvFieldList {
		if envName == p {
			return false
		}
	}
	return true
}
