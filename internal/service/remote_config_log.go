package service

import (
	"context"

	"github.com/1219796395/myProject2/api/errorcode"
	pb "github.com/1219796395/myProject2/api/operationlog/remoteconfiglog"
	"github.com/1219796395/myProject2/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

// RemoteConfigLogService is a remote config log service.
type RemoteConfigLogService struct {
	pb.UnimplementedRemoteConfigLogServer

	logic *biz.RemoteConfigLogLogic
	log   *log.Helper
}

// NewRemoteConfigLogService ...
func NewRemoteConfigLogService(logic *biz.RemoteConfigLogLogic, logger log.Logger) *RemoteConfigLogService {
	return &RemoteConfigLogService{logic: logic, log: log.NewHelper(logger)}
}

func convToPBRemoteConfigLog(log *biz.RemoteConfigLog) *pb.RemoteConfigLogDetail {
	return &pb.RemoteConfigLogDetail{
		Appid:        log.Appid,
		Env:          log.Env,
		Channel:      log.Channel,
		Platform:     log.Platform,
		ConfigName:   log.ConfigName,
		Operation:    log.Operation,
		Operator:     log.Operator,
		UpdateBefore: log.UpdateBefore,
		UpdateAfter:  log.UpdateAfter,
		CreateTime:   uint64(log.Ctime.UnixMilli()),
	}
}

func batchConvToPBRemoteConfigLog(logs []*biz.RemoteConfigLog) []*pb.RemoteConfigLogDetail {
	logPBs := make([]*pb.RemoteConfigLogDetail, 0, len(logs))
	for _, log := range logs {
		logPBs = append(logPBs, convToPBRemoteConfigLog(log))
	}
	return logPBs
}

func (s *RemoteConfigLogService) GetRemoteConfigLogList(ctx context.Context, in *pb.GetRemoteConfigLogListReq) (
	*pb.GetRemoteConfigLogListRsp, error) {
	if !checkBatchReadPlatform(in.Platform) || !checkBatchReadChannel(in.Channel) {
		s.log.WithContext(ctx).Errorf("[GetRemoteConfigLogList] invalid param! req = %+v", in)
		return nil, errorcode.ErrorBadRequest("invalid param")
	}
	rcLog := &biz.RemoteConfigLog{
		Appid:      in.Appid,
		Env:        in.Env,
		Channel:    in.Channel,
		Platform:   in.Platform,
		ConfigName: in.ConfigName,
		Operation:  in.Operation,
		Operator:   in.Operator,
		StartTime:  in.StartTime,
		EndTime:    in.EndTime,
	}
	logs, err := s.logic.GetRemoteConfigLogList(ctx, rcLog)
	if err != nil {
		if _, ok := errorcode.Errors_value[errors.FromError(err).GetReason()]; ok {
			s.log.WithContext(ctx).Errorf("[GetRemoteConfigLogList] fail! req = %+v, err = %+v", in, err)
			return nil, err
		}
		s.log.WithContext(ctx).Errorf("[GetRemoteConfigLogList] fail! req = %+v, err = %+v", in, err)
		return nil, errorcode.ErrorInternalServerError("get remote config log list fail")
	}

	return &pb.GetRemoteConfigLogListRsp{
		Logs: batchConvToPBRemoteConfigLog(logs),
	}, nil
}
