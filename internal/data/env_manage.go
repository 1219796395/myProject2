package data

import (
	"context"
	"game-config/api/errorcode"
	"game-config/internal/biz"
	"game-config/internal/conf"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type envManageRepo struct {
	data *Data
	log  *log.Helper
}

// NewEnvManageRepo .
func NewEnvManageRepo(data *Data, logger log.Logger) biz.EnvManageRepo {
	return &envManageRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Env 环境管理-数据层 实体类
type Env struct {
	Id       int64     `gorm:"column:id;primaryKey"` // id, auto increment
	Appid    uint32    `gorm:"column:appid"`
	Field    string    `gorm:"column:field"`
	Name     string    `gorm:"column:name"`
	Remark   string    `gorm:"column:remark"`
	Operator string    `gorm:"column:operator"`
	Ctime    time.Time `gorm:"column:ctime;->"`
	Mtime    time.Time `gorm:"column:mtime;->"`
}

func (Env) TableName() string {
	return "env_manage"
}

func convToBizEnv(env *Env) *biz.Env {
	return &biz.Env{
		Appid:    env.Appid,
		Field:    env.Field,
		Name:     env.Name,
		Remark:   env.Remark,
		Operator: env.Operator,
		Ctime:    env.Ctime,
		Mtime:    env.Mtime,
	}
}

func batchConvToBizEnv(envs []*Env) []*biz.Env {
	envBizs := make([]*biz.Env, 0, len(envs))
	for _, env := range envs {
		envBizs = append(envBizs, convToBizEnv(env))
	}
	return envBizs
}

func convToDBEnv(env *biz.Env) *Env {
	return &Env{
		Appid:    env.Appid,
		Field:    env.Field,
		Name:     env.Name,
		Remark:   env.Remark,
		Operator: env.Operator,
		Ctime:    env.Ctime,
		Mtime:    env.Mtime,
	}
}

func (r *envManageRepo) GetEnvList(ctx context.Context, env *biz.Env) ([]*biz.Env, error) {
	envDBs := []*Env{}
	find := "`appid`=? "
	args := []interface{}{env.Appid}
	if env.Field != "" {
		find += " and `field`=? "
		args = append(args, env.Field)
	}
	if env.Name != "" {
		find += " and `name`=? "
		args = append(args, env.Name)
	}
	if err := r.data.db.WithContext(ctx).Where(find, args...).Find(&envDBs).Error; err != nil {
		return nil, err
	}
	envs := batchConvToBizEnv(envDBs)
	res, err := getPresetEnvs(env)
	if err != nil {
		return nil, err
	}
	res = append(res, envs...)
	return res, nil
}

// 获取预设的环境列表
func getPresetEnvs(env *biz.Env) ([]*biz.Env, error) {
	bc, err := conf.GetConf()
	if err != nil {
		return nil, err
	}
	res := []*biz.Env{}
	fields := bc.GetBiz().GetEnvFieldList()

	for _, f := range fields {
		if (env.Field == "" || env.Field == f) && (env.Name == "" || env.Name == f) {
			res = append(res, genPresetEnv(env.Appid, f))
		}
	}

	return res, nil
}

// 生成预设的环境数据
func genPresetEnv(appid uint32, field string) *biz.Env {
	return &biz.Env{
		Appid:    appid,
		Field:    field,
		Name:     field,
		IsPreset: true,
		Remark:   "preset env",
		Ctime:    time.Unix(1693497600, 0), //2023.9.1 00:00:00
		Mtime:    time.Unix(1693497600, 0), //2023.9.1 00:00:00,
	}
}

func (r *envManageRepo) GetEnv(ctx context.Context, env *biz.Env) (*biz.Env, error) {

	envs, err := getPresetEnvs(env)
	if err != nil {
		return nil, err
	}
	if len(envs) == 1 {
		return envs[0], nil
	}

	find := "`appid`=? and `name`=?"
	args := []interface{}{env.Appid, env.Name}
	if env.Field != "" {
		find += " and `field`=? "
		args = append(args, env.Field)
	}
	envDB := &Env{}
	if err := r.data.db.WithContext(ctx).Where(find, args...).First(envDB).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorcode.ErrorEnvManageNotFound("env not found")
		}
	}
	return convToBizEnv(envDB), nil
}

func (r *envManageRepo) CreateEnv(ctx context.Context, env *biz.Env) error {
	envDB := convToDBEnv(env)
	if err := r.data.db.WithContext(ctx).Create(envDB).Error; err != nil {
		if errMySQL, ok := err.(*mysql.MySQLError); ok && errMySQL.Number == 1062 {
			return errorcode.ErrorEnvManageNameAlreadyExist("env name already exist")
		}
		return err
	}
	return nil
}

func (r *envManageRepo) UpdateEnv(ctx context.Context, env *biz.Env) error {
	update := make(map[string]interface{})
	if env.NewName != "" {
		update["name"] = env.NewName
	}
	if env.Remark != "" {
		update["remark"] = env.Remark
	}
	update["operator"] = env.Operator
	// 在db上修改环境信息
	if err := r.data.db.WithContext(ctx).Model(&Env{}).Where("`appid`=? and `field`=? and `name`=?",
		env.Appid, env.Field, env.Name).Updates(update).Error; err != nil {
		return err
	}
	return nil
}

func (r *envManageRepo) DeleteEnv(ctx context.Context, env *biz.Env) error {
	// 在db上删除环境信息
	if err := r.data.db.WithContext(ctx).Where("`appid`=? and `field`=? and `name`=?",
		env.Appid, env.Field, env.Name).Delete(&Env{}).Error; err != nil {
		return err
	}
	return nil
}
