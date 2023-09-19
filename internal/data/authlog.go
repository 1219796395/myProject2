package data

import (
	"context"
	"time"

	"github.com/1219796395/myProject2/internal/biz"
	"github.com/1219796395/myProject2/internal/biz/bo"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/jinzhu/copier"
)

// daos
// 日志表model
type AuthLog struct {
	Id         uint32    `gorm:"column:id;primaryKey"`       //日志主键
	AppId      uint32    `gorm:"column:app_id"`              //产品id,1为明日方舟，2为玩家社区
	OperatorId uint32    `gorm:"column:operator_id"`         //操作员id
	Operation  uint32    `gorm:"column:operation"`           //操作类型
	Content    string    `gorm:"column:content"`             //操作详情，权限模块日志中该字段为一段可读模版，其他日志中该字段为空
	UserName   string    `gorm:"column:user_name"`           //冗余字段，用于模糊搜索
	RoleName   string    `gorm:"column:role_name"`           //冗余字段，用于模糊搜索
	UserId     uint32    `gorm:"column:user_id"`             //冗余字段
	RoleId     uint32    `gorm:"column:role_id"`             //冗余字段
	CreatedAt  time.Time `gorm:"column:created_at;<-:false"` //日志创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;<-:false"` //日志更新时间
}

func (AuthLog) TableName() string {
	return "auth_log"
}

func convertToDBAuthLog(authLog *bo.AuthLog) *AuthLog {
	var dbAuthLog AuthLog
	copier.CopyWithOption(&dbAuthLog, authLog, copier.Option{DeepCopy: true})
	return &dbAuthLog
}

func batchConvertToBizAuthLogs(authLogs []*AuthLog) []*bo.AuthLog {
	var bizAuthLogs []*bo.AuthLog
	copier.CopyWithOption(&bizAuthLogs, &authLogs, copier.Option{DeepCopy: true})
	return bizAuthLogs
}

type AuthLogRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthLogRepo(data *Data, logger log.Logger) biz.AuthLogRepo {
	return &AuthLogRepo{data: data, log: log.NewHelper(logger)}
}

func (r *AuthLogRepo) CreateAuthLog(ctx context.Context, log *bo.AuthLog) error {
	if err := r.data.db.WithContext(ctx).Create(convertToDBAuthLog(log)).Error; err != nil {
		return err
	}
	return nil
}

func (r *AuthLogRepo) SearchAuthLog(ctx context.Context, req *biz.ListAuthLogRequest) ([]*bo.AuthLog, error) {
	var logs []*AuthLog
	stat := r.data.db.WithContext(ctx)
	// if need to filter by operator id
	if len(req.OperatorIds) != 0 {
		stat = stat.Where("operator_id in ?", req.OperatorIds)
	}
	if req.StartId != 0 {
		stat = stat.Where("id < ?", req.StartId)
	}
	if len(req.ContentKey) != 0 {
		stat = stat.Where("role_name like ? OR user_name like ?", "%"+req.ContentKey+"%", "%"+req.ContentKey+"%")
	}
	if err := stat.Order("id desc").Limit(int(req.PageSize)).Find(&logs).Error; err != nil {
		return nil, err
	}
	return batchConvertToBizAuthLogs(logs), nil
}
