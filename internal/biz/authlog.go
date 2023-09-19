package biz

import (
	"context"

	"game-config/internal/biz/bo"

	"github.com/go-kratos/kratos/v2/log"
)

type SearchAuthLogRequest struct {
	AppId          uint32
	ContentKey     string // 用于搜索权限模块操作日志的关键字
	OperatorKey    string // 用于按照管理员用户姓名/昵称搜索日志的关键字
	StartId        uint32
	StartOperation uint32 // 枚举值范围，为左开区间
	EndOperation   uint32 // 枚举值范围，为右闭区间
	PageSize       uint32
}

type ListAuthLogRequest struct {
	AppId          uint32
	ContentKey     string   // 用于搜索权限模块操作日志的关键字
	OperatorIds    []uint32 // 按照操作员id进行检索
	StartId        uint32
	StartOperation uint32 // 枚举值范围，为左闭区间
	EndOperation   uint32 // 枚举值范围，为右闭区间
	PageSize       uint32
}

type AuthLogRepo interface {
	CreateAuthLog(ctx context.Context, log *bo.AuthLog) error
	SearchAuthLog(ctx context.Context, req *ListAuthLogRequest) ([]*bo.AuthLog, error)
}

type AuthLogLogic struct {
	log           *log.Helper
	logRepo       AuthLogRepo
	adminUserRepo AdminUserRepo
}

func NewAuthLogUsecase(logger log.Logger, logRepo AuthLogRepo, adminUserRepo AdminUserRepo) *AuthLogLogic {
	return &AuthLogLogic{log: log.NewHelper(logger), logRepo: logRepo, adminUserRepo: adminUserRepo}
}

func (uc *AuthLogLogic) CreateAuthLog(ctx context.Context, log *bo.AuthLog) error {
	return uc.logRepo.CreateAuthLog(ctx, log)
}

func (uc *AuthLogLogic) SearchAuthLog(ctx context.Context, req SearchAuthLogRequest) ([]*bo.AuthLog, error) {
	operatorIds := make([]uint32, 0)
	// if need to search operators bu their names
	if len(req.OperatorKey) != 0 {
		operators, err := uc.adminUserRepo.SearchAdminUsersByName(ctx, req.OperatorKey)
		if err != nil {
			return nil, err
		}
		// if it did not hit any users
		if len(operators) == 0 {
			return nil, nil
		}
		for i := 0; i < len(operators); i++ {
			operatorIds = append(operatorIds, operators[i].Id)
		}
	}

	return uc.logRepo.SearchAuthLog(ctx, &ListAuthLogRequest{
		AppId:          req.AppId,
		ContentKey:     req.ContentKey,
		OperatorIds:    operatorIds,
		StartId:        req.StartId,
		StartOperation: req.StartOperation,
		EndOperation:   req.EndOperation,
		PageSize:       req.PageSize,
	})
}
