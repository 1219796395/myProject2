package data

import (
	"context"
	stderr "errors"
	"fmt"
	"game-config/internal/biz"
	"game-config/internal/biz/bo"
	"math/rand"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/copier"
	"github.com/spf13/cast"
)

// daos
type AdminUser struct {
	Id        uint32    `json:"id" gorm:"column:id;primaryKey"`
	HgId      string    `json:"hg_id" gorm:"column:hg_id"`
	Email     string    `json:"email" gorm:"column:email"`
	Nickname  string    `json:"nickname" gorm:"column:nick_name"`
	Name      string    `json:"name" gorm:"column:name"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;<-:false"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updated_at;<-:false"`
}

func (AdminUser) TableName() string {
	return "admin_user"
}

func convertToDBAdminUser(adminUser *bo.AdminUser) *AdminUser {
	var dbAdminUser AdminUser
	copier.CopyWithOption(&dbAdminUser, adminUser, copier.Option{DeepCopy: true})
	return &dbAdminUser
}

func convertToBizAdminUser(adminUser *AdminUser) *bo.AdminUser {
	var bizAdminUser bo.AdminUser
	copier.CopyWithOption(&bizAdminUser, adminUser, copier.Option{DeepCopy: true})
	return &bizAdminUser
}

func batchConvertToBizAdminUsers(adminUsers []*AdminUser) []*bo.AdminUser {
	var bizAdminUsers []*bo.AdminUser
	copier.CopyWithOption(&bizAdminUsers, &adminUsers, copier.Option{DeepCopy: true})
	return bizAdminUsers
}

type UserAppRelation struct {
	Id         uint32    `gorm:"column:id;primaryKey"`
	UserId     uint32    `gorm:"column:user_id"`
	AppId      uint32    `gorm:"column:app_id"`
	Status     uint32    `gorm:"column:status"`
	Identity   uint32    `gorm:"column:identity"`
	CreatedAt  time.Time `gorm:"column:created_at;<-:false"`
	UpdatedAt  time.Time `gorm:"column:updated_at;<-:false"`
	AssignedAt time.Time `gorm:"column:assigned_at"`
}

func (UserAppRelation) TableName() string {
	return "user_app_relation"
}

func convertToDBUserAppRelation(userAppRelation *bo.UserAppRelation) *UserAppRelation {
	var dbUserAppRelation UserAppRelation
	copier.CopyWithOption(&dbUserAppRelation, userAppRelation, copier.Option{DeepCopy: true})
	return &dbUserAppRelation
}

func convertToBizUserAppRelation(userAppRelation *UserAppRelation) *bo.UserAppRelation {
	var bizUserAppRelation bo.UserAppRelation
	copier.CopyWithOption(&bizUserAppRelation, userAppRelation, copier.Option{DeepCopy: true})
	return &bizUserAppRelation
}

func batchConvertToBizUserAppRelations(userAppRelations []*UserAppRelation) []*bo.UserAppRelation {
	var bizUserAppRelations []*bo.UserAppRelation
	copier.CopyWithOption(&bizUserAppRelations, userAppRelations, copier.Option{DeepCopy: true})
	return bizUserAppRelations
}

type AdminRole struct {
	Id                  uint32    `gorm:"column:id;primaryKey"`
	Name                string    `gorm:"column:name"`
	Description         string    `gorm:"column:description"`
	ResourceList        string    `gorm:"column:resource_list"`
	Status              uint32    `gorm:"column:status"`
	AppId               uint32    `gorm:"column:app_id"`
	ResourceListContent []string  `gorm:"-"`
	CreatedAt           time.Time `gorm:"column:created_at;<-:false"`
	UpdatedAt           time.Time `gorm:"column:updated_at;<-:false"`
}

func (AdminRole) TableName() string {
	return "admin_role"
}

func convertToDBAdminRole(adminRole *bo.AdminRole) *AdminRole {
	var dbAdminRole AdminRole
	copier.CopyWithOption(&dbAdminRole, adminRole, copier.Option{DeepCopy: true})
	return &dbAdminRole
}

func convertToBizAdminRole(adminRole *AdminRole) *bo.AdminRole {
	var bizAdminRole bo.AdminRole
	copier.CopyWithOption(&bizAdminRole, adminRole, copier.Option{DeepCopy: true})
	return &bizAdminRole
}

func batchConvertToBizAdminRoles(adminRoles []*AdminRole) []*bo.AdminRole {
	var bizAdminRoles []*bo.AdminRole
	copier.CopyWithOption(&bizAdminRoles, &adminRoles, copier.Option{DeepCopy: true})
	return bizAdminRoles
}

type UserRoleRelation struct {
	Id        uint32    `gorm:"column:id;primaryKey"`
	UserId    uint32    `gorm:"column:user_id"`
	RoleId    uint32    `gorm:"column:role_id"`
	AppId     uint32    `gorm:"column:app_id"`
	CreatedAt time.Time `gorm:"column:created_at;<-:false"`
	UpdatedAt time.Time `gorm:"column:updated_at;<-:false"`
}

func (UserRoleRelation) TableName() string {
	return "user_role_relation"
}

func convertToDbUserRoleRelation(userRoleRelation *bo.UserRoleRelation) *UserAppRelation {
	var dbUserRoleRelation UserAppRelation
	copier.CopyWithOption(&dbUserRoleRelation, userRoleRelation, copier.Option{DeepCopy: true})
	return &dbUserRoleRelation
}

func convertToBizUserRoleRelation(userRoleRelation *UserAppRelation) *bo.UserAppRelation {
	var bizUserRoleRelation bo.UserAppRelation
	copier.CopyWithOption(&bizUserRoleRelation, userRoleRelation, copier.Option{DeepCopy: true})
	return &bizUserRoleRelation
}

func batchConvertToDBUserRoleRelations(userRoleRelations []*bo.UserRoleRelation) []*UserRoleRelation {
	var dbUserRoleRelation []*UserRoleRelation
	copier.CopyWithOption(&dbUserRoleRelation, &userRoleRelations, copier.Option{DeepCopy: true})
	return dbUserRoleRelation
}

// repo
type adminUserRepo struct {
	data *Data
	log  *log.Helper
}

func NewAdminUserRepo(data *Data, logger log.Logger) biz.AdminUserRepo {
	return &adminUserRepo{data: data, log: log.NewHelper(logger)}
}

func (r *adminUserRepo) CreateAdminUser(ctx context.Context, adminUser *bo.AdminUser) (*bo.AdminUser, error) {
	dbAdminUser := convertToDBAdminUser(adminUser)
	if err := r.data.db.WithContext(ctx).Create(dbAdminUser).Error; err != nil {
		return nil, err
	}
	return convertToBizAdminUser(dbAdminUser), nil
}

func (r *adminUserRepo) UpdateAdminUserStatus(ctx context.Context, userId uint32, status uint32) error {
	if err := r.data.db.WithContext(ctx).Table("admin_user").
		Where("id = ?", userId).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminUserRepo) GetAllAdminUsers(ctx context.Context) ([]*bo.AdminUser, error) {
	var users []*AdminUser
	if err := r.data.db.WithContext(ctx).Where("1=1").Find(&users).Error; err != nil {
		return nil, err
	}
	return batchConvertToBizAdminUsers(users), nil
}

func (r *adminUserRepo) GetAdminUserByEmail(ctx context.Context, email string) (*bo.AdminUser, error) {
	var adminUser AdminUser
	if err := r.data.db.WithContext(ctx).Where("email like ?", email).Find(&adminUser).Limit(1).Error; err != nil {
		return nil, err
	}
	if adminUser.Id == 0 {
		return nil, nil
	}
	return convertToBizAdminUser(&adminUser), nil
}

func (r *adminUserRepo) GetAdminUserByHgId(ctx context.Context, hgId string) (*bo.AdminUser, error) {
	var adminUser AdminUser
	if err := r.data.db.WithContext(ctx).Where("hg_id = ?", hgId).Find(&adminUser).Limit(1).Error; err != nil {
		return nil, err
	}
	if adminUser.Id == 0 {
		return nil, nil
	}
	return convertToBizAdminUser(&adminUser), nil
}

func (r *adminUserRepo) GetAdminUserById(ctx context.Context, id uint32) (*bo.AdminUser, error) {
	var adminUser AdminUser
	if err := r.data.db.WithContext(ctx).Where("id = ?", id).Find(&adminUser).Limit(1).Error; err != nil {
		return nil, err
	}
	if adminUser.Id == 0 {
		return nil, nil
	}
	return convertToBizAdminUser(&adminUser), nil
}

func (r *adminUserRepo) GetAnyAdminUserByIds(ctx context.Context, userIds []uint32) ([]*bo.AdminUser, error) {
	if len(userIds) == 0 {
		return nil, nil
	}
	var adminUsers []*AdminUser
	if err := r.data.db.WithContext(ctx).Where("id in ?", userIds).Find(&adminUsers).Error; err != nil {
		return nil, err
	}
	return batchConvertToBizAdminUsers(adminUsers), nil
}

func (r *adminUserRepo) SearchAdminUsersByName(ctx context.Context, key string) ([]*bo.AdminUser, error) {
	var users []*AdminUser
	if err := r.data.db.Where("nick_name like ? OR name like ?", "%"+key+"%", "%"+key+"%").
		Find(&users).Error; err != nil {
		return nil, err
	}
	return batchConvertToBizAdminUsers(users), nil
}

func (r *adminUserRepo) UpdateAdminUser(ctx context.Context, name string, nickname string, email string, id uint32) error {
	updateMap := map[string]interface{}{
		"nick_name": nickname,
		"name":      name,
		"email":     email,
	}
	if err := r.data.db.WithContext(ctx).Table("admin_user").
		Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminUserRepo) ListAdminUser(ctx context.Context, appId uint32, pageNum uint32, pageSize uint32) (users []*bo.AdminUser, total uint32, err error) {
	var (
		adminUsers []*AdminUser
		userIds    []uint32
		totalCount uint32
	)
	if pageNum == 1 {
		if err := r.data.db.WithContext(ctx).Table("user_app_relation").
			Where("app_id = ?", &appId).Select("COUNT(*)").Find(&totalCount).Error; err != nil {
			return nil, 0, err
		}
	}

	// get user ids in one specific app id by pages
	if err := r.data.db.WithContext(ctx).Table("user_app_relation").
		Where("app_id = ?", appId).Select("user_id").
		Offset(int((pageNum - 1) * pageSize)).Limit(int(pageSize)).Find(&userIds).Error; err != nil {
		return nil, 0, err
	}

	// grab those users by ids
	if err := r.data.db.WithContext(ctx).Where("id in ?", userIds).Find(&adminUsers).Error; err != nil {
		return nil, 0, err
	}
	return batchConvertToBizAdminUsers(adminUsers), totalCount, nil
}

/*
*
redis
*/
func (r *adminUserRepo) GetAdminTokenCache(ctx context.Context, token string) (uint32, error) {
	key := fmt.Sprintf(bo.AdminTokenKey, token)
	str, err := r.data.rdb.Get(ctx, key).Result()
	switch {
	case stderr.Is(err, redis.Nil):
		r.log.WithContext(ctx).Errorf("[adminUserRepo.GetAdminTokenCache] redis nil %s", token)
		return 0, nil
	case err != nil:
		r.log.WithContext(ctx).Errorf("[GetAdminTokenCache]: %v", err)
		return 0, err
	case str == "":
		r.log.WithContext(ctx).Errorf("[adminUserRepo.GetAdminTokenCache] redis str null %s", token)
		return 0, nil
	}

	return cast.ToUint32(str), nil
}

func (r *adminUserRepo) GetAndRefreshAdminTokenCache(ctx context.Context, token string) (userId uint32, err error) {
	key := fmt.Sprintf(bo.AdminTokenKey, token)
	str, err := r.data.rdb.Get(ctx, key).Result()
	switch {
	case err == redis.Nil:
		return 0, nil
	case err != nil:
		r.log.WithContext(ctx).Errorf("[GetAndRefreshAdminTokenCache]: %v", err)
		return 0, err
	case str == "":
		return 0, nil
	}

	expire := bo.AdminTokenValidSeconds + rand.Intn(3600)
	r.data.rdb.Expire(ctx, key, time.Duration(expire)*time.Second)

	return cast.ToUint32(str), nil
}

func (r *adminUserRepo) CacheAdminToken(ctx context.Context, token string, id uint32) error {
	key := fmt.Sprintf(bo.AdminTokenKey, token)
	expire := bo.AdminTokenValidSeconds + rand.Intn(600)
	result, err := r.data.rdb.SetNX(ctx, key, id, time.Duration(expire)*time.Second).Result()
	if err != nil {
		r.log.WithContext(ctx).Errorf("[CacheAdminToken]: %v", err)
		return err
	}
	if !result {
		return stderr.New("token exists")
	}

	return nil
}

/*
role
*/
func (r *adminUserRepo) CreateAdminRole(ctx context.Context, role *bo.AdminRole) (*bo.AdminRole, error) {
	dbRole := convertToDBAdminRole(role)
	if err := r.data.db.WithContext(ctx).Save(role).Error; err != nil {
		return nil, err
	}
	return convertToBizAdminRole(dbRole), nil
}

func (r *adminUserRepo) DeleteAdminRole(ctx context.Context, id uint32) error {
	if err := r.data.db.WithContext(ctx).Table("admin_role").
		Where("id = ?", id).Update("status", bo.AdminRoleDeleted).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminUserRepo) UpdateAdminRole(ctx context.Context, name string, description string, resources string, id uint32) error {
	// update name, description, resorce list
	updateMap := map[string]interface{}{
		"name":          name,
		"description":   description,
		"resource_list": resources,
	}
	if err := r.data.db.WithContext(ctx).Table("admin_role").
		Where("id = ?", id).Updates(updateMap).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminUserRepo) GetAdminRole(ctx context.Context, id uint32) (*bo.AdminRole, error) {
	var role AdminRole
	if err := r.data.db.WithContext(ctx).
		Where("id = ? AND status != ?", id, bo.AdminRoleDeleted).Find(&role).Limit(1).Error; err != nil {
		return nil, err
	}
	if role.Id == 0 {
		return nil, nil
	}

	// unserialize resource list
	var resourceList []string
	json.UnmarshalFromString(role.ResourceList, &resourceList)
	role.ResourceListContent = resourceList
	return convertToBizAdminRole(&role), nil
}

func (r *adminUserRepo) GetAdminRolesByName(ctx context.Context, name string, appId uint32) ([]*bo.AdminRole, error) {
	var roles []*AdminRole
	if err := r.data.db.WithContext(ctx).Where("name = ? AND app_id = ?", name, appId).Find(&roles).Error; err != nil {
		return nil, err
	}

	// unserialize resource list
	for i := 0; i < len(roles); i++ {
		var resourceList []string
		json.UnmarshalFromString(roles[i].ResourceList, &resourceList)
		roles[i].ResourceListContent = resourceList
	}
	return batchConvertToBizAdminRoles(roles), nil
}

func (r *adminUserRepo) GetAdminRolesByIds(ctx context.Context, ids []uint32, appId uint32) ([]*bo.AdminRole, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var roles []*AdminRole
	if err := r.data.db.WithContext(ctx).Where("id in ? AND app_id = ?", ids, appId).Where("status != ?", bo.AdminRoleDeleted).
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return batchConvertToBizAdminRoles(roles), nil
}

func (r *adminUserRepo) SearchAdminRoles(ctx context.Context, key string, appId uint32, startId uint32, pageSize int) ([]*bo.AdminRole, error) {
	var roles []*AdminRole
	tx := r.data.db.WithContext(ctx).Table("admin_role").Where("`status` != ?", bo.AdminRoleDeleted).
		Where("app_id = ?", appId).Order("id desc")
	// search by survey's name on both prefix and postfix
	if len(key) != 0 {
		tx = tx.Where("name like '" + key + "%'")
	}
	// paging
	if startId != 0 {
		tx = tx.Where("id < ?", startId)
	}
	if err := tx.Limit(pageSize).Find(&roles).Error; err != nil {
		return nil, err
	}
	// unserializing resources
	for i := 0; i < len(roles); i++ {
		var resources []string
		json.UnmarshalFromString(roles[i].ResourceList, &resources)
		roles[i].ResourceListContent = resources
	}
	return batchConvertToBizAdminRoles(roles), nil
}

func (r *adminUserRepo) ListAdminRolesByUserId(ctx context.Context, userId uint32, appId uint32) ([]*bo.AdminRole, error) {
	// get role ids
	var userRoles []*UserRoleRelation
	tx := r.data.db.WithContext(ctx).Where("user_id = ?", userId)
	if appId != 0 {
		tx = tx.Where("app_id = ?", appId)
	}
	if err := tx.Find(&userRoles).Error; err != nil {
		return nil, err
	}
	ids := make([]uint32, 0)
	for i := 0; i < len(userRoles); i++ {
		ids = append(ids, userRoles[i].RoleId)
	}

	// get role models
	var roles []*AdminRole
	if err := r.data.db.WithContext(ctx).Table("admin_role").
		Where("id in ? AND `status` = ?", ids, bo.AdminRoleNormal).Find(&roles).Error; err != nil {
		return nil, err
	}
	// unserializing resources
	for i := 0; i < len(roles); i++ {
		var resources []string
		json.UnmarshalFromString(roles[i].ResourceList, &resources)
		roles[i].ResourceListContent = resources
	}
	return batchConvertToBizAdminRoles(roles), nil
}

func (r *adminUserRepo) DeleteUserRoleRelationsByUserId(ctx context.Context, userId uint32, appId uint32) error {
	if err := r.data.db.Where("user_id = ? AND app_id = ?", userId, appId).
		Delete(&UserRoleRelation{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminUserRepo) DeleteUserRoleRelationsByUserIdAndRoleIds(ctx context.Context, userId uint32, appId uint32, roleIds []uint32) error {
	if err := r.data.db.Where("user_id = ? AND app_id = ?", userId, appId).
		Where("role_id IN ?", roleIds).Delete(&UserRoleRelation{}).Error; err != nil {
		return err
	}
	return nil
}

/*
user-role relation
*/
func (r *adminUserRepo) CreateUserRoleRelations(ctx context.Context, roles []*bo.UserRoleRelation) error {
	if len(roles) == 0 {
		return nil
	}
	if err := r.data.db.WithContext(ctx).CreateInBatches(batchConvertToDBUserRoleRelations(roles), 100).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminUserRepo) DeleteUserRoleRelation(ctx context.Context, userId uint32, roleId uint32) error {
	if err := r.data.db.WithContext(ctx).
		Where("user_id = ?", userId).Where("role_id = ?", roleId).Delete(&UserRoleRelation{}).Error; err != nil {
		return err
	}
	return nil
}

/*
user_app_relation
*/
func (r *adminUserRepo) CreateUserAppRelation(ctx context.Context, userId uint32, appId uint32) (isDuplicated bool, err error) {
	if err := r.data.db.WithContext(ctx).Create(&UserAppRelation{
		UserId:     userId,
		AppId:      appId,
		Status:     bo.UserAppIdNormal,
		Identity:   bo.IdentityOrdinary,
		AssignedAt: time.Now(),
	}).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if stderr.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return true, nil
		}
		return false, err
	}
	return false, nil
}

func (r *adminUserRepo) DeleteUserAppRelation(ctx context.Context, userId uint32, appId uint32) error {
	if err := r.data.db.WithContext(ctx).
		Where("user_id = ?", userId).Where("app_id = ?", appId).Delete(&UserAppRelation{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminUserRepo) GetUserAppRelation(ctx context.Context, userId uint32, appId uint32) (*bo.UserAppRelation, error) {
	var appIdModel UserAppRelation
	if err := r.data.db.WithContext(ctx).Table("user_app_relation").
		Where("user_id = ? AND app_id = ?", userId, appId).Limit(1).Find(&appIdModel).Error; err != nil {
		return nil, err
	}
	if appIdModel.Id == 0 {
		return nil, nil
	}
	return convertToBizUserAppRelation(&appIdModel), nil
}

func (r *adminUserRepo) GetUserAppRelations(ctx context.Context, userId uint32) ([]*bo.UserAppRelation, error) {
	var appIds []*UserAppRelation
	if err := r.data.db.WithContext(ctx).Where("user_id = ?", userId).Find(&appIds).Error; err != nil {
		return nil, err
	}
	return batchConvertToBizUserAppRelations(appIds), nil
}

func (r *adminUserRepo) UpdateUserAppRelationStatus(ctx context.Context, userId uint32, appId uint32, status uint32) error {
	updateMap := map[string]interface{}{
		"status": status,
	}
	if err := r.data.db.WithContext(ctx).Table("user_app_relation").
		Where("user_id = ?", userId).Where("app_id = ?", appId).Updates(updateMap).Error; err != nil {
		return err
	}
	return nil
}

func (r *adminUserRepo) UpdateUserAppRelationAssignedAt(ctx context.Context, userId uint32, appId uint32) error {
	if err := r.data.db.WithContext(ctx).Table("user_app_relation").
		Where("user_id = ? AND app_id = ?", userId, appId).Update("assigned_at", time.Now()).Error; err != nil {
		return err
	}
	return nil
}
