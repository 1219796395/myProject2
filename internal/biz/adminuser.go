package biz

import (
	"context"
	"game-config/internal/biz/bo"
	"game-config/internal/conf"
	"game-config/tool"

	"github.com/go-kratos/kratos/v2/log"
)

type AdminUserRepo interface {
	/*
		admin user
	*/
	CreateAdminUser(ctx context.Context, adminUser *bo.AdminUser) (*bo.AdminUser, error)
	UpdateAdminUserStatus(ctx context.Context, id uint32, status uint32) error
	// GetAdminUserByHgId
	// this method's results include the deleted users
	GetAdminUserByHgId(ctx context.Context, hgId string) (*bo.AdminUser, error)
	GetAdminUserById(ctx context.Context, id uint32) (*bo.AdminUser, error)
	// GetAnyAdminUserByIds
	// this method's results include the deleted users
	GetAnyAdminUserByIds(ctx context.Context, userIds []uint32) ([]*bo.AdminUser, error)
	SearchAdminUsersByName(ctx context.Context, key string) ([]*bo.AdminUser, error)
	ListAdminUser(ctx context.Context, appId uint32, pageNum uint32, pageSize uint32) (users []*bo.AdminUser, total uint32, err error)
	GetAllAdminUsers(ctx context.Context) ([]*bo.AdminUser, error)
	// for updating admin users by mdm infos
	UpdateAdminUser(ctx context.Context, name string, nickname string, email string, id uint32) error

	/*
		token cache
	*/
	GetAdminTokenCache(ctx context.Context, token string) (userId uint32, err error)
	CacheAdminToken(ctx context.Context, token string, id uint32) error
	GetAndRefreshAdminTokenCache(ctx context.Context, token string) (userId uint32, err error)

	/*
		admin roles
	*/
	CreateAdminRole(ctx context.Context, role *bo.AdminRole) (*bo.AdminRole, error)
	DeleteAdminRole(ctx context.Context, id uint32) error
	UpdateAdminRole(ctx context.Context, name string, description string, resources string, id uint32) error
	GetAdminRole(ctx context.Context, id uint32) (*bo.AdminRole, error)
	// GetAdminRolesByName
	// this method's result include roles that are deleted
	GetAdminRolesByName(ctx context.Context, name string, appId uint32) ([]*bo.AdminRole, error)
	GetAdminRolesByIds(ctx context.Context, ids []uint32, appId uint32) ([]*bo.AdminRole, error)
	SearchAdminRoles(ctx context.Context, key string, appId uint32, startId uint32, pageSize int) ([]*bo.AdminRole, error)
	ListAdminRolesByUserId(ctx context.Context, userId uint32, appId uint32) ([]*bo.AdminRole, error)

	/*
		user-role-relation
	*/
	CreateUserRoleRelations(ctx context.Context, userRoles []*bo.UserRoleRelation) error
	DeleteUserRoleRelation(ctx context.Context, userId uint32, roleId uint32) error
	DeleteUserRoleRelationsByUserIdAndRoleIds(ctx context.Context, userId uint32, appId uint32, roleIds []uint32) error
	DeleteUserRoleRelationsByUserId(ctx context.Context, userId uint32, appId uint32) error

	/*
		user-appid-relation
	*/
	CreateUserAppRelation(ctx context.Context, userId uint32, appId uint32) (isDuplicated bool, err error)
	DeleteUserAppRelation(ctx context.Context, userId uint32, appId uint32) error
	GetUserAppRelation(ctx context.Context, userId uint32, appId uint32) (*bo.UserAppRelation, error)
	GetUserAppRelations(ctx context.Context, userId uint32) ([]*bo.UserAppRelation, error)
	UpdateUserAppRelationStatus(ctx context.Context, userId uint32, appId uint32, status uint32) error
}

type AdminUserLogic struct {
	repo    AdminUserRepo
	bizConf *conf.Biz
	log     *log.Helper
}

func NewAdminUserLogic(repo AdminUserRepo, bc *conf.Bootstrap, logger log.Logger) *AdminUserLogic {
	return &AdminUserLogic{
		repo:    repo,
		bizConf: bc.Biz,
		log:     log.NewHelper(logger),
	}
}

func (uc *AdminUserLogic) CreateAdminUser(ctx context.Context, adminUser *bo.AdminUser) (*bo.AdminUser, error) {
	return uc.repo.CreateAdminUser(ctx, adminUser)
}

// delete user in one simple appid, this user may still exist in others
func (uc *AdminUserLogic) DeleteAdminUserByAppId(ctx context.Context, userId uint32, appId uint32) error {
	// delete the user's role under that app id
	if err := uc.repo.DeleteUserRoleRelationsByUserId(ctx, userId, appId); err != nil {
		return err
	}
	// delete user's relation to that app id
	if err := uc.repo.DeleteUserAppRelation(ctx, userId, appId); err != nil {
		return err
	}
	return nil
}

func (uc *AdminUserLogic) UpdateAdminUser(ctx context.Context, name string, nickName string, email string, id uint32) error {
	return uc.repo.UpdateAdminUser(ctx, name, nickName, email, id)
}

func (uc *AdminUserLogic) UpdateAdminUserStatus(ctx context.Context, userId uint32, status uint32) error {
	return uc.repo.UpdateAdminUserStatus(ctx, userId, status)
}

func (uc *AdminUserLogic) GetAdminUserByHgId(ctx context.Context, hgId string) (*bo.AdminUser, error) {
	return uc.repo.GetAdminUserByHgId(ctx, hgId)
}

func (uc *AdminUserLogic) GetAdminUserById(ctx context.Context, id uint32) (*bo.AdminUser, error) {
	return uc.repo.GetAdminUserById(ctx, id)
}

func (uc *AdminUserLogic) ListAdminUser(ctx context.Context, appId uint32, pageNum uint32, pageSize uint32) (users []*bo.AdminUser, total uint32, err error) {
	return uc.repo.ListAdminUser(ctx, appId, pageNum, pageSize)
}

func (uc *AdminUserLogic) GetAnyAdminUserByIds(ctx context.Context, userIds []uint32) ([]*bo.AdminUser, error) {
	return uc.repo.GetAnyAdminUserByIds(ctx, userIds)
}

func (uc *AdminUserLogic) SearchAdminUsersByName(ctx context.Context, key string) ([]*bo.AdminUser, error) {
	return uc.repo.SearchAdminUsersByName(ctx, key)
}

func (uc *AdminUserLogic) GetAllAdminUsers(ctx context.Context) ([]*bo.AdminUser, error) {
	return uc.repo.GetAllAdminUsers(ctx)
}

// token
func (uc *AdminUserLogic) GetAdminTokenCache(ctx context.Context, token string) (userId uint32, err error) {
	return uc.repo.GetAdminTokenCache(ctx, token)
}

func (uc *AdminUserLogic) GenerateAndCacheAdminToken(ctx context.Context, id uint32) (token string, err error) {
	token = tool.RandStr(32)
	err = uc.repo.CacheAdminToken(ctx, token, id)
	if err != nil {
		return "", err
	}
	return token, nil
}

// role
func (uc *AdminUserLogic) CreateAdminRole(ctx context.Context, role *bo.AdminRole) (*bo.AdminRole, error) {
	return uc.repo.CreateAdminRole(ctx, role)
}

func (uc *AdminUserLogic) UpdateAdminRole(ctx context.Context, name string,
	description string, resources []string, status uint32, id uint32) error {
	resourcesStr, _ := json.MarshalToString(resources)
	return uc.repo.UpdateAdminRole(ctx, name, description, resourcesStr, id)
}

func (uc *AdminUserLogic) DeleteAdminRole(ctx context.Context, id uint32) error {
	return uc.repo.DeleteAdminRole(ctx, id)
}

func (uc *AdminUserLogic) GetAdminRole(ctx context.Context, id uint32) (*bo.AdminRole, error) {
	return uc.repo.GetAdminRole(ctx, id)
}

func (uc *AdminUserLogic) GetAdminRolesByName(ctx context.Context, name string, appId uint32) ([]*bo.AdminRole, error) {
	return uc.repo.GetAdminRolesByName(ctx, name, appId)
}

func (uc *AdminUserLogic) GetAdminRolesByIds(ctx context.Context, ids []uint32, appId uint32) ([]*bo.AdminRole, error) {
	return uc.repo.GetAdminRolesByIds(ctx, ids, appId)
}

func (uc *AdminUserLogic) SearchRole(ctx context.Context, key string, appId uint32, startId uint32, pageSize int) ([]*bo.AdminRole, error) {
	return uc.repo.SearchAdminRoles(ctx, key, appId, startId, pageSize)
}

func (uc *AdminUserLogic) ListRolesByUserId(ctx context.Context, userId uint32, appId uint32) ([]*bo.AdminRole, error) {
	return uc.repo.ListAdminRolesByUserId(ctx, userId, appId)
}

func (uc *AdminUserLogic) DeleteUserRolesByUserId(ctx context.Context, userId uint32, appId uint32) error {
	return uc.repo.DeleteUserRoleRelationsByUserId(ctx, userId, appId)
}

func (uc *AdminUserLogic) DeleteUserRoleRelationsByUserIdAndRoleIds(ctx context.Context, userId uint32, appId uint32, roleIds []uint32) error {
	return uc.repo.DeleteUserRoleRelationsByUserIdAndRoleIds(ctx, userId, appId, roleIds)
}

func (uc *AdminUserLogic) CreateUserRoles(ctx context.Context, userRoles []*bo.UserRoleRelation) error {
	return uc.repo.CreateUserRoleRelations(ctx, userRoles)
}

func (uc *AdminUserLogic) DeleteRoleUser(ctx context.Context, userId uint32, roleId uint32) error {
	return uc.repo.DeleteUserRoleRelation(ctx, userId, roleId)
}

func (uc *AdminUserLogic) CreateUserAppRelation(ctx context.Context, userId uint32, appId uint32) (isDuplicated bool, err error) {
	return uc.repo.CreateUserAppRelation(ctx, userId, appId)
}

func (uc *AdminUserLogic) GetUserAppRelation(ctx context.Context, userId uint32, appId uint32) (*bo.UserAppRelation, error) {
	return uc.repo.GetUserAppRelation(ctx, userId, appId)
}

func (uc *AdminUserLogic) GetUserAppRelations(ctx context.Context, userId uint32) ([]*bo.UserAppRelation, error) {
	return uc.repo.GetUserAppRelations(ctx, userId)
}

func (uc *AdminUserLogic) UpdateUserAppIdStatus(ctx context.Context, userId uint32, appId uint32, status uint32) error {
	return uc.repo.UpdateUserAppRelationStatus(ctx, userId, appId, status)
}

func (uc *AdminUserLogic) Can(ctx context.Context, user *bo.AdminUser, operation string, appId uint32) (bool, error) {
	// get all resource point needed
	resourcesNeeded, ok := bo.OperaionResourceMap[operation]
	if !ok {
		return false, nil
	}

	// if the map is empty, means that this operation need no auth
	if len(resourcesNeeded) == 0 {
		return true, nil
	}

	// if the user is a super
	userApp, err := uc.repo.GetUserAppRelation(ctx, user.Id, appId)
	if err != nil {
		return false, err
	}
	if userApp == nil {
		return false, nil
	}
	if userApp.Identity == bo.IdentitySuper {
		return true, nil
	}
	if userApp.Status == bo.UserAppIdBanned {
		return false, nil
	}

	// get roles of this user
	roles, err := uc.ListRolesByUserId(ctx, user.Id, appId)
	if err != nil {
		return false, err
	}

	// see if any resource needed is acquired
	for i := 0; i < len(roles); i++ {
		resources := roles[i].ResourceListContent
		for j := 0; j < len(resources); j++ {
			// if any resource in the list get fullfilled, let'em pass
			if resourcesNeeded[resources[j]] {
				return true, nil
			}
		}
	}
	// no resource meet, deny the access
	return false, nil
}

func (uc *AdminUserLogic) VerifyUser(ctx context.Context, token string) (*bo.AdminUser, error) {
	if len(token) == 0 {
		return nil, nil
	}
	// token缓存
	adminUId, err := uc.GetAdminTokenCache(ctx, token)
	if err != nil {
		// 拿不到credCache，报错
		return nil, err
	}
	if adminUId == 0 {
		return nil, nil
	}
	// 校验用户原始数据
	adminUser, err := uc.GetAdminUserById(ctx, adminUId)
	if err != nil {
		// db链接错误，直接放行
		return nil, err
	}

	return adminUser, nil
}
