package service

import (
	"context"
	"fmt"
	"strconv"

	pb "github.com/1219796395/myProject2/api/auth"
	errorPb "github.com/1219796395/myProject2/api/errorcode"
	"github.com/1219796395/myProject2/internal/biz"
	"github.com/1219796395/myProject2/internal/biz/bo"
	"github.com/1219796395/myProject2/internal/client"
	"github.com/1219796395/myProject2/internal/conf"

	"github.com/1219796395/myProject2/internal/middleware"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

const (
	// user identities
	UserIdentitySuper    = 1
	UserIdentityNotSuper = 0

	// user status
	UserStatusBanned    = 1
	UserStatusNotBanned = 0
)

type AdminAuthService struct {
	pb.UnimplementedAuthServer

	log          *log.Helper
	conf         *conf.Bootstrap
	adminUserBiz *biz.AdminUserLogic
	logBiz       *biz.AuthLogLogic
	client       *client.AdminUserClient
}

func NewAdminAuthService(
	logger log.Logger,
	conf *conf.Bootstrap,
	adminUserBiz *biz.AdminUserLogic,
	logBiz *biz.AuthLogLogic,
	client *client.AdminUserClient,
) *AdminAuthService {
	return &AdminAuthService{
		log:          log.NewHelper(logger),
		conf:         conf,
		adminUserBiz: adminUserBiz,
		logBiz:       logBiz,
		client:       client,
	}
}

func (s *AdminAuthService) GenerateTokenByCode(ctx context.Context, req *pb.GenerateTokenByCodeRequest) (*pb.GenerateTokenByCodeResponse, error) {
	// verify through calling sso system
	ssoToken, status, err := s.client.VerifySSOCode(ctx, req.AuthCode, req.RedirectUrl)
	if err != nil {
		s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode] VerifySSOCode err: ", err.Error())
		return nil, errorPb.ErrorInternalServerError("INTER_SERVER_ERROR")
	}
	if status != bo.SSOLogInSuccess {
		s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode] sso login err")
		return nil, errorPb.ErrorInternalServerError("sso error")
	}

	// get user info by sso token
	ssoInfo, status, err := s.client.GetUserInfoBySSOToken(ctx, ssoToken)
	if err != nil {
		s.log.WithContext(ctx).Errorf("[AdminAuthService.GenerateTokenByCode] get info by sso token err: %s, code is: %d", err.Error(), status)
		return nil, errorPb.ErrorInternalServerError("INTER_SERVER_ERROR")
	}

	// query our db for this user
	adminUser, err := s.adminUserBiz.GetAdminUserByHgId(ctx, ssoInfo.HgId)
	if err != nil {
		s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode] query hg from db: ", err.Error())
		return nil, errors.InternalServer("INTER_SERVER_ERROR", "database error")
	}

	// pull info from mdm
	infos, err := s.client.QueryUserInfosFromMdm(ctx, []string{ssoInfo.HgId})
	if err != nil {
		s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode] mdm query error ", err.Error())

		// if can aquire user info from mdm
		// if the user exist, generate token, or reject
		if adminUser != nil {
			// generate token
			token, err := s.adminUserBiz.GenerateAndCacheAdminToken(ctx, adminUser.Id)
			if err != nil {
				s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode]", err.Error())
				return nil, errorPb.ErrorInternalServerError("generate token error")
			}
			return &pb.GenerateTokenByCodeResponse{Token: token}, nil
		}

		// unable to create new user, can not response
		s.log.WithContext(ctx).Errorf("[AdminAuthService.GenerateTokenByCode] failed to create user")
		return nil, errors.InternalServer("INTER_SERVER_ERROR", "mdm error")
	}

	// if this user not exist in mdm
	if len(infos) == 0 {
		s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode] no such user")
		return nil, errorPb.ErrorNoSuchUser("no such user")
	}
	info := infos[0]

	if adminUser == nil {
		// if there is no such user in our system, create one with info from mdm
		adminUser = &bo.AdminUser{
			HgId:     info.HgId,
			Email:    info.Email,
			Nickname: info.Nickname,
			Name:     info.Name,
		}
		adminUser, err = s.adminUserBiz.CreateAdminUser(ctx, adminUser)
		if err != nil {
			s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode] create admin user ", err.Error())
			return nil, errorPb.ErrorInternalServerError("INTERNAL SERVER ERROR")
		}
	} else {
		// if we have this user in our system, update it with info from mdm
		if info.Name != adminUser.Name || info.Nickname != adminUser.Nickname || info.Email != adminUser.Email {
			err := s.adminUserBiz.UpdateAdminUser(ctx, info.Name, info.Nickname, info.Email, adminUser.Id)
			if err != nil {
				s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode] update user info ", err.Error())
			}
		}
	}

	// generate token
	token, err := s.adminUserBiz.GenerateAndCacheAdminToken(ctx, adminUser.Id)
	if err != nil {
		s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode] generate token in redis ", err.Error())
		return nil, errorPb.ErrorInternalServerError("INTER_SERVER_ERROR")
	}

	return &pb.GenerateTokenByCodeResponse{Token: token}, nil
}

func (s *AdminAuthService) SearchUserMdmInfo(ctx context.Context, req *pb.SearchUserMdmInfoRequest) (*pb.SearchUserMdmInfoResponse, error) {
	// search
	users, err := s.client.SearchUsersFromMdm(ctx, req.Key)
	if err != nil {
		s.log.WithContext(ctx).Error("[AdminAuthService.SearchUserMdmInfo]", err.Error())
		return nil, errorPb.ErrorInternalServerError("mdm error")
	}

	// render response
	userPbs := make([]*pb.MdmUserInfo, 0)
	for i := 0; i < len(users); i++ {
		userPbs = append(userPbs, &pb.MdmUserInfo{
			HgId:     users[i].HgId,
			Name:     users[i].Name,
			Nickname: users[i].Nickname,
			Email:    users[i].Email,
		})
	}
	return &pb.SearchUserMdmInfoResponse{Users: userPbs}, nil
}

func (s *AdminAuthService) Me(ctx context.Context, req *pb.MeRequest) (*pb.MeResponse, error) {
	// convert the token into user model
	user, err := s.adminUserBiz.VerifyUser(ctx, req.Token)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.GetMe] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("sso error")
	}
	if user == nil {
		s.log.WithContext(ctx).Errorf("[adminservice.GetMe] user %s not exist ", req.Token)
		return nil, errorPb.ErrorNoSuchUser("no such user")
	}

	// get user info
	userAgg, err := s.PackUserAgg(ctx, user)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.GetMe] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("use case error")
	}
	return &pb.MeResponse{User: userAgg}, nil
}

func (s *AdminAuthService) CreateAdminUser(ctx context.Context, req *pb.CreateAdminUserRequest) (*pb.CreateAdminUserResponse, error) {
	// get common context
	adminUserContext := middleware.GetAdminUser(ctx)

	// check if this user already exist or not
	user, err := s.adminUserBiz.GetAdminUserByHgId(ctx, req.HgId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.CreateAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// if user not exist
	if user == nil {
		// pull user info from mdm
		infos, err := s.client.QueryUserInfosFromMdm(ctx, []string{req.HgId})
		if err != nil || len(infos) == 0 {
			s.log.WithContext(ctx).Error("[AdminAuthService.GenerateTokenByCode]", err.Error())
		}
		if len(infos) == 0 {
			return nil, errorPb.ErrorNoSuchUser("no such user")
		}
		info := infos[0]

		// create a new user
		user = &bo.AdminUser{
			HgId:     info.HgId,
			Email:    info.Email,
			Nickname: info.Nickname,
			Name:     info.Name,
		}
		if user, err = s.adminUserBiz.CreateAdminUser(ctx, user); err != nil {
			s.log.WithContext(ctx).Error("[adminservice.CreateAdminUser] ", err.Error())
			return nil, errorPb.ErrorInternalServerError("databse error")
		}
	}

	// added a app id user this user account
	isDuplicated, err := s.adminUserBiz.CreateUserAppRelation(ctx, user.Id, req.Common.AppId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.CreateAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	if isDuplicated {
		return nil, errorPb.ErrorUserAlreadyExist("user already exist")
	}

	// generate an auth log
	// create log
	if err := s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
		AppId:      req.Common.AppId,
		OperatorId: adminUserContext.Id,
		Operation:  bo.OperationCreateUser,
		Content:    fmt.Sprintf(bo.LogContentTempCreateUser, user.Name),
		UserName:   user.Name,
		UserId:     user.Id,
	}); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.CreateAdminUser] ", err.Error())
	}

	return &pb.CreateAdminUserResponse{Id: user.Id}, nil
}

func (s *AdminAuthService) DeleteAdminUser(ctx context.Context, req *pb.DeleteAdminUserRequest) (*pb.DeleteAdminUserResponse, error) {
	// get common context
	adminUserContext := middleware.GetAdminUser(ctx)

	// a user do not have right to delete it self
	if req.UserId == adminUserContext.Id {
		return nil, errorPb.ErrorNotAuthorized("not authorized")
	}
	appIdModel, err := s.adminUserBiz.GetUserAppRelation(ctx, req.UserId, req.Common.AppId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	if appIdModel == nil {
		return nil, errorPb.ErrorNoSuchUser("no such user")
	}

	// can not delete a super admin
	if appIdModel.Identity == bo.IdentitySuper {
		return nil, errorPb.ErrorNotAuthorized("not authorized")
	}

	// delete user by app id
	if err := s.adminUserBiz.DeleteAdminUserByAppId(ctx, req.UserId, req.Common.AppId); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// generate an auth log
	user, err := s.adminUserBiz.GetAdminUserById(ctx, req.UserId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminUser] ", err.Error())
		return &pb.DeleteAdminUserResponse{}, nil
	}

	// create log
	if err := s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
		AppId:      req.Common.AppId,
		OperatorId: adminUserContext.Id,
		Operation:  bo.OperationDeleteUser,
		Content:    fmt.Sprintf(bo.LogContentTempDeleteUser, user.Name),
		UserName:   user.Name,
		UserId:     user.Id,
	}); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminUser] ", err.Error())
	}

	return &pb.DeleteAdminUserResponse{}, nil
}

func (s *AdminAuthService) BanAdminUser(ctx context.Context, req *pb.BanAdminUserRequest) (*pb.BanAdminUserResponse, error) {
	// get common context
	adminUserContext := middleware.GetAdminUser(ctx)

	// TODO, validate super user
	// a user do not have right to it self
	if req.UserId == adminUserContext.Id {
		return nil, errorPb.ErrorNotAuthorized("not authorized")
	}

	appIdModel, err := s.adminUserBiz.GetUserAppRelation(ctx, req.UserId, req.Common.AppId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("databse error")
	}
	// can not delete a super admin
	if appIdModel.Identity == bo.IdentitySuper {
		return nil, errorPb.ErrorNotAuthorized("not authorized")
	}

	// validate user status
	user, err := s.adminUserBiz.GetAdminUserById(ctx, req.UserId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("databse error")
	}
	if user == nil {
		return nil, errorPb.ErrorNoSuchUser("no such user")
	}

	// update user status to banned status
	if err := s.adminUserBiz.UpdateUserAppIdStatus(ctx, req.UserId, req.Common.AppId, bo.UserAppIdBanned); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// generate an auth log
	// create log
	if err := s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
		AppId:      req.Common.AppId,
		OperatorId: adminUserContext.Id,
		Operation:  bo.OperationBanUser,
		Content:    fmt.Sprintf(bo.LogContentTempBanUser, user.Name),
		UserName:   user.Name,
		UserId:     user.Id,
	}); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.BanAdminUser] ", err.Error())
	}
	return &pb.BanAdminUserResponse{}, nil
}

func (s *AdminAuthService) RestoreAdminUser(ctx context.Context, req *pb.RestoreAdminUserRequest) (*pb.RestoreAdminUserResponse, error) {
	// get common context
	adminUserContext := middleware.GetAdminUser(ctx)

	// validate user status
	user, err := s.adminUserBiz.GetAdminUserById(ctx, req.UserId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.RestoreAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	if user == nil {
		return nil, errorPb.ErrorNoSuchUser("no such user")
	}

	appIdModel, err := s.adminUserBiz.GetUserAppRelation(ctx, req.UserId, req.Common.AppId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.RestoreAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	// can not delete a super admin
	if appIdModel.Identity == bo.IdentitySuper {
		return nil, errorPb.ErrorNotAuthorized("not authorized")
	}

	// update user status to normal
	if err := s.adminUserBiz.UpdateUserAppIdStatus(ctx, req.UserId, req.Common.AppId, bo.UserAppIdNormal); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.RestoreAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// generate an auth log
	// create log
	if err := s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
		AppId:      req.Common.AppId,
		OperatorId: adminUserContext.Id,
		Operation:  bo.OperationRestoreUser,
		Content:    fmt.Sprintf(bo.LogContentTempRestoreUser, user.Name),
		UserName:   user.Name,
		UserId:     user.Id,
	}); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.RestoreAdminUser] ", err.Error())
	}
	return &pb.RestoreAdminUserResponse{}, nil
}

func (s *AdminAuthService) ListAdminUser(ctx context.Context, req *pb.ListAdminUserRequest) (*pb.ListAdminUserResponse, error) {
	// validate page token
	var page uint32
	num, err := strconv.Atoi(req.Page)
	if err != nil {
		return nil, errorPb.ErrorBadRequest("illegal page token")
	}
	page = uint32(num)

	// get users and user count(count is not zero only when page is 1)
	users, total, err := s.adminUserBiz.ListAdminUser(ctx, req.Common.AppId, page, req.PageSize)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.ListAdminUser] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// render response
	userAggs := make([]*pb.AdminUserListAgg, 0)
	for i := 0; i < len(users); i++ {
		userAgg, err := s.PackUserListAgg(ctx, users[i], req.Common.AppId)
		if err != nil {
			s.log.WithContext(ctx).Error("[adminservice.ListAdminUser] ", err.Error())
			return nil, errorPb.ErrorInternalServerError("use case error")
		}
		userAggs = append(userAggs, userAgg)
	}

	return &pb.ListAdminUserResponse{
		TotalCount: total,
		List:       userAggs,
	}, nil
}

func (s *AdminAuthService) RetrieveAdminRoleResource(ctx context.Context, req *pb.RetrieveAdminRoleResourceRequest) (*pb.RetrieveAdminRoleResourceResponse, error) {
	role, err := s.adminUserBiz.GetAdminRole(ctx, req.RoleId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.RetrieveAdminRoleResource] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	if role == nil {
		return nil, errorPb.ErrorNoSuchRole("no such role")
	}
	rolePb := &pb.AdminRole{
		Id:           role.Id,
		Name:         role.Name,
		Description:  role.Description,
		ResourceList: role.ResourceListContent,
		Status:       role.Status,
		AppId:        role.AppId,
		CreatedAtTs:  uint32(role.CreatedAt.Unix()),
		UpdatedAtTs:  uint32(role.UpdatedAt.Unix()),
	}
	return &pb.RetrieveAdminRoleResourceResponse{Role: rolePb}, nil
}

func (s *AdminAuthService) SaveAdminRole(ctx context.Context, req *pb.SaveAdminRoleRequest) (*pb.SaveAdminRoleResponse, error) {
	// get common context
	adminUserContext := middleware.GetAdminUser(ctx)

	// create a new role
	if req.Role.Id == 0 {
		// validate if any role with same name exist
		roles, err := s.adminUserBiz.GetAdminRolesByName(ctx, req.Role.Name, req.Common.AppId)
		if err != nil {
			s.log.WithContext(ctx).Error("[adminservice.SaveAdminRole] ", err.Error())
			return nil, errorPb.ErrorInternalServerError("database error")
		}
		for i := 0; i < len(roles); i++ {
			if roles[i].Status != bo.AdminRoleDeleted {
				// role with same name exist
				return nil, errorPb.ErrorRoleAlreadyExist("role already exist")
			}
		}

		// create a new role
		resourceList, _ := json.MarshalToString(req.Role.ResourceList)
		role := &bo.AdminRole{
			Name:         req.Role.Name,
			Description:  req.Role.Description,
			ResourceList: resourceList,
			Status:       bo.AdminRoleNormal,
			AppId:        req.Role.AppId,
		}
		if role, err = s.adminUserBiz.CreateAdminRole(ctx, role); err != nil {
			s.log.WithContext(ctx).Error("[adminservice.SaveAdminRole] ", err.Error())
			return nil, errorPb.ErrorInternalServerError("database error")
		}

		// generate an auth log
		// create log
		err = s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
			AppId:      req.Common.AppId,
			OperatorId: adminUserContext.Id,
			Operation:  bo.OperationCreateRole,
			Content:    fmt.Sprintf(bo.LogContentTempCreateRole, req.Role.Name),
			RoleName:   role.Name,
			RoleId:     role.Id,
		})
		if err != nil {
			s.log.WithContext(ctx).Error("[adminservice.SaveAdminRole] ", err.Error())
		}
		return &pb.SaveAdminRoleResponse{}, nil
	}

	// validate if this role exist
	role, err := s.adminUserBiz.GetAdminRole(ctx, req.Role.Id)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.SaveAdminRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	if req.Role.Id != 0 && role == nil {
		return nil, errorPb.ErrorNoSuchRole("no such role")
	}
	// update a exsiting role, either exist one or deleted ones
	if err := s.adminUserBiz.UpdateAdminRole(ctx, req.Role.Name, req.Role.Description,
		req.Role.ResourceList, bo.AdminRoleNormal, req.Role.Id); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.SaveAdminRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// generate an auth log
	// create log
	if err := s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
		AppId:      req.Common.AppId,
		OperatorId: adminUserContext.Id,
		Operation:  bo.OperationUpdateRole,
		Content:    fmt.Sprintf(bo.LogContentTempUpdateRole, req.Role.Name),
		RoleName:   role.Name,
		RoleId:     role.Id,
	}); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.SaveAdminRole] ", err.Error())
	}
	return &pb.SaveAdminRoleResponse{}, nil
}

func (s *AdminAuthService) DeleteAdminRole(ctx context.Context, req *pb.DeleteAdminRoleRequest) (*pb.DeleteAdminRoleResponse, error) {
	// get common context
	adminUserContext := middleware.GetAdminUser(ctx)

	role, err := s.adminUserBiz.GetAdminRole(ctx, req.RoleId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	if role == nil {
		return nil, errorPb.ErrorNoSuchRole("no such role")
	}
	if err := s.adminUserBiz.DeleteAdminRole(ctx, role.Id); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("databse error")
	}

	// generate an auth log
	// create log
	if err := s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
		AppId:      req.Common.AppId,
		OperatorId: adminUserContext.Id,
		Operation:  bo.OperationDeleteRole,
		Content:    fmt.Sprintf(bo.LogContentTempDeleteRole, role.Name),
		RoleName:   role.Name,
		RoleId:     role.Id,
	}); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.DeleteAdminRole] ", err.Error())
	}
	return &pb.DeleteAdminRoleResponse{}, nil
}

func (s *AdminAuthService) SearchAdminRole(ctx context.Context, req *pb.SearchAdminRoleRequest) (*pb.SearchAdminRoleResponse, error) {
	// validate page token
	var startId uint32 = 0
	if req.Next != 0 {
		startId = req.Next
	}
	roles, err := s.adminUserBiz.SearchRole(ctx, req.Key, req.Common.AppId, startId, int(req.PageSize))
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.SearchAdminRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	// update page token and has more
	var next uint32
	if len(roles) != 0 {
		next = roles[len(roles)-1].Id
	}
	hasMore := true
	if len(roles) != int(req.PageSize) {
		hasMore = false
	}
	// render role pbs
	rolePbs := make([]*pb.AdminRole, 0)
	for i := 0; i < len(roles); i++ {
		rolePbs = append(rolePbs, &pb.AdminRole{
			Id:           roles[i].Id,
			Name:         roles[i].Name,
			Description:  roles[i].Description,
			ResourceList: roles[i].ResourceListContent,
			Status:       roles[i].Status,
			AppId:        roles[i].AppId,
			CreatedAtTs:  uint32(roles[i].CreatedAt.Unix()),
			UpdatedAtTs:  uint32(roles[i].UpdatedAt.Unix()),
		})
	}
	return &pb.SearchAdminRoleResponse{
		List:    rolePbs,
		Next:    next,
		HasMore: hasMore,
	}, nil
}

func (s *AdminAuthService) UpdateAdminUserRole(ctx context.Context, req *pb.UpdateAdminUserRoleRequest) (*pb.UpdateAdminUserRoleResponse, error) {
	// get common context
	adminUserContext := middleware.GetAdminUser(ctx)

	// validation
	// validate if this user exist
	user, err := s.adminUserBiz.GetAdminUserById(ctx, req.UserId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.UpdateAdminUserRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	if user == nil {
		return nil, errorPb.ErrorNoSuchUser("no such user")
	}
	// validate roles, roles must be in the same app id
	roles, err := s.adminUserBiz.GetAdminRolesByIds(ctx, req.RoleIds, req.Common.AppId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.UpdateAdminUserRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	// validate if all roles in request exist
	if len(roles) != len(req.RoleIds) {
		return nil, errorPb.ErrorNoSuchRole("no such role")
	}
	// validate appid
	for i := 0; i < len(roles); i++ {
		if roles[i].AppId != req.Common.AppId {
			return nil, errorPb.ErrorRoleAlreadyExist("role already exist")
		}
	}

	// diff
	// get users current roles and new roles for diff and generate auth logs
	oldRoles, err := s.adminUserBiz.ListRolesByUserId(ctx, req.UserId, req.Common.AppId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.UpdateAdminUserRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}
	newRoles, err := s.adminUserBiz.GetAdminRolesByIds(ctx, req.RoleIds, req.Common.AppId)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.UpdateAdminUserRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// diff old roles with new roles
	rolesToInsert := diffRoles(newRoles, oldRoles)
	rolesToDelete := diffRoles(oldRoles, newRoles)

	// insert those roles need to insert
	newUserRoles := make([]*bo.UserRoleRelation, 0)
	for i := 0; i < len(rolesToInsert); i++ {
		newUserRoles = append(newUserRoles, &bo.UserRoleRelation{
			UserId: req.UserId,
			RoleId: rolesToInsert[i].Id,
			AppId:  req.Common.AppId,
		})
	}
	if err := s.adminUserBiz.CreateUserRoles(ctx, newUserRoles); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.UpdateAdminUserRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// insert log
	for i := 0; i < len(rolesToInsert); i++ {
		if err := s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
			AppId:      req.Common.AppId,
			OperatorId: adminUserContext.Id,
			Operation:  bo.OperationAssignRole,
			Content:    fmt.Sprintf(bo.LogContentTempAssignRole, user.Name, rolesToInsert[i].Name),
			UserName:   user.Name,
			UserId:     user.Id,
			RoleName:   rolesToInsert[i].Name,
			RoleId:     rolesToInsert[i].Id,
		}); err != nil {
			s.log.WithContext(ctx).Error("[adminservice.UpdateAdminUserRole] ", err.Error())
		}
	}

	// delete those roles need to delete
	roleIdsToDelete := make([]uint32, 0)
	for i := 0; i < len(rolesToDelete); i++ {
		roleIdsToDelete = append(roleIdsToDelete, rolesToDelete[i].Id)
	}
	if err := s.adminUserBiz.DeleteUserRoleRelationsByUserIdAndRoleIds(ctx, req.UserId, req.Common.AppId, roleIdsToDelete); err != nil {
		s.log.WithContext(ctx).Error("[adminservice.UpdateAdminUserRole] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// delete log
	for i := 0; i < len(rolesToDelete); i++ {
		if err := s.logBiz.CreateAuthLog(ctx, &bo.AuthLog{
			AppId:      req.Common.AppId,
			OperatorId: adminUserContext.Id,
			Operation:  bo.OperationDepriveRole,
			Content:    fmt.Sprintf(bo.LogContentTempDepriveRole, user.Name, rolesToDelete[i].Name),
			UserName:   user.Name,
			UserId:     user.Id,
			RoleName:   rolesToDelete[i].Name,
			RoleId:     rolesToDelete[i].Id,
		}); err != nil {
			s.log.WithContext(ctx).Error("[adminservice.UpdateAdminUserRole] ", err.Error())
			return &pb.UpdateAdminUserRoleResponse{}, nil
		}
	}

	return &pb.UpdateAdminUserRoleResponse{}, nil
}

func (s *AdminAuthService) ListAuthLog(ctx context.Context, req *pb.ListAuthLogRequest) (*pb.ListAuthLogResponse, error) {
	// validate page token
	var startId uint32 = 0
	if req.Next != 0 {
		startId = req.Next
	}

	// query
	logs, err := s.logBiz.SearchAuthLog(ctx, biz.SearchAuthLogRequest{
		AppId:          req.Common.AppId,
		OperatorKey:    req.OperatorKey,
		ContentKey:     req.ContentKey,
		StartId:        startId,
		StartOperation: 101,
		EndOperation:   199,
		PageSize:       req.PageSize,
	})
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.listLog] ", err.Error())
		return nil, errorPb.ErrorInternalServerError("database error")
	}

	// update page token and has more
	var next uint32
	if len(logs) != 0 {
		next = logs[len(logs)-1].Id
	}
	hasMore := true
	if len(logs) != int(req.PageSize) {
		hasMore = false
	}

	// make auth log aggs
	operatorIds := make([]uint32, 0)
	for i := 0; i < len(logs); i++ {
		operatorIds = append(operatorIds, logs[i].OperatorId)
	}
	// get operators
	operators, err := s.adminUserBiz.GetAnyAdminUserByIds(ctx, operatorIds)
	if err != nil {
		s.log.WithContext(ctx).Error("[adminservice.listLog] ", err.Error())
	}
	operatorMap := make(map[uint32]*bo.AdminUser)
	for i := 0; i < len(operators); i++ {
		operatorMap[operators[i].Id] = operators[i]
	}
	logAggs := make([]*pb.AuthLogAgg, 0)
	// render response
	for i := 0; i < len(logs); i++ {
		authLog := &pb.AuthLog{
			Id:          logs[i].Id,
			OperatorId:  logs[i].OperatorId,
			Operation:   logs[i].Operation,
			Content:     logs[i].Content,
			CreatedAtTs: uint32(logs[i].CreatedAt.Unix()),
			UpdatedAtTs: uint32(logs[i].UpdatedAt.Unix()),
		}

		// get operators
		var authOperator *pb.AdminUser
		operator := operatorMap[logs[i].OperatorId]
		if operator != nil {
			authOperator = &pb.AdminUser{
				Id:          operator.Id,
				Email:       operator.Email,
				Name:        operator.Name,
				NickName:    operator.Nickname,
				CreatedAtTs: uint32(operator.CreatedAt.Unix()),
				UpdatedAtTs: uint32(operator.UpdatedAt.Unix()),
			}
		}

		// make a auth log agg
		authLogAgg := pb.AuthLogAgg{
			AuthLog:  authLog,
			Operator: authOperator,
		}
		logAggs = append(logAggs, &authLogAgg)
	}

	return &pb.ListAuthLogResponse{List: logAggs,
		Next:    next,
		HasMore: hasMore}, nil
}

func (s *AdminAuthService) PackUserAgg(ctx context.Context, user *bo.AdminUser) (*pb.AdminUserAgg, error) {
	userPb := &pb.AdminUser{
		Id:          user.Id,
		HgId:        user.HgId,
		Email:       user.Email,
		NickName:    user.Nickname,
		Name:        user.Name,
		CreatedAtTs: uint32(user.CreatedAt.Unix()),
		UpdatedAtTs: uint32(user.CreatedAt.Unix()),
	}

	// get app ids
	userApps, err := s.adminUserBiz.GetUserAppRelations(ctx, user.Id)
	if err != nil {
		return nil, err
	}
	appIdModelMap := make(map[uint32]*bo.UserAppRelation, 0)
	for i := 0; i < len(userApps); i++ {
		appIdModelMap[userApps[i].AppId] = userApps[i]
	}

	// get roles
	roleModels, err := s.adminUserBiz.ListRolesByUserId(ctx, user.Id, 0)
	if err != nil {
		return nil, err
	}

	// organize roles by app id
	appIdMap := make(map[uint32][]*pb.AdminRole)
	for i := 0; i < len(roleModels); i++ {
		appId := roleModels[i].AppId
		rolePb := &pb.AdminRole{
			Id:           roleModels[i].Id,
			Name:         roleModels[i].Name,
			Description:  roleModels[i].Description,
			ResourceList: roleModels[i].ResourceListContent,
			Status:       roleModels[i].Status,
			AppId:        roleModels[i].AppId,
			CreatedAtTs:  uint32(roleModels[i].CreatedAt.Unix()),
			UpdatedAtTs:  uint32(roleModels[i].UpdatedAt.Unix()),
		}
		roleArray, ok := appIdMap[appId]
		if ok {
			roleArray = append(roleArray, rolePb)
			appIdMap[appId] = roleArray
			continue
		}
		newRoleArray := make([]*pb.AdminRole, 1)
		newRoleArray[0] = rolePb
		appIdMap[appId] = newRoleArray
	}

	// organize resources by app id
	apps := make([]*pb.AdminAppAgg, 0)
	for i := 0; i < len(userApps); i++ {
		// if is bannd, not pack info of this app
		if userApps[i].Status == bo.UserAppIdBanned {
			continue
		}

		userApp := userApps[i]
		var appAgg pb.AdminAppAgg
		appAgg.Roles = appIdMap[userApp.AppId]
		appAgg.AppId = userApp.AppId
		resourceMap := make(map[string]bool, 0)
		resources := make([]string, 0)
		for j := 0; j < len(appIdMap[userApp.AppId]); j++ {
			role := appIdMap[userApp.AppId][j]
			roleResources := role.ResourceList
			for k := 0; k < len(roleResources); k++ {
				resourceMap[roleResources[k]] = true
			}
		}
		for resource := range resourceMap {
			resources = append(resources, resource)
		}
		appAgg.Resources = resources
		apps = append(apps, &appAgg)
		isSuper := UserIdentitySuper
		if identity := appIdModelMap[userApp.AppId].Identity; identity == bo.IdentityOrdinary {
			isSuper = UserIdentityNotSuper
		}
		if isSuper == UserIdentitySuper {
			rolePbs := appAgg.Roles
			roleList := make([]*pb.AdminRole, 0)
			roleList = append(roleList, &pb.AdminRole{
				Name:         "超级管理员",
				Description:  "超级管理员",
				ResourceList: bo.AllResources,
				Status:       bo.AdminRoleNormal,
				AppId:        appAgg.AppId,
			})
			if len(rolePbs) != 0 {
				roleList = append(roleList, rolePbs...)
			}
			appAgg.Roles = roleList
			appAgg.Resources = bo.AllResources
		}
		appAgg.IsSuper = uint32(isSuper)
	}
	return &pb.AdminUserAgg{
		User: userPb,
		Apps: apps,
	}, nil
}

func (s *AdminAuthService) PackUserListAgg(ctx context.Context, user *bo.AdminUser, appId uint32) (*pb.AdminUserListAgg, error) {
	userPb := &pb.AdminUser{
		Id:          user.Id,
		Email:       user.Email,
		NickName:    user.Nickname,
		Name:        user.Name,
		CreatedAtTs: uint32(user.CreatedAt.Unix()),
		UpdatedAtTs: uint32(user.CreatedAt.Unix()),
	}

	// get roles
	rolePbs := make([]*pb.AdminRole, 0)
	roleModels, err := s.adminUserBiz.ListRolesByUserId(ctx, user.Id, appId)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(roleModels); i++ {
		rolePb := &pb.AdminRole{
			Id:           roleModels[i].Id,
			Name:         roleModels[i].Name,
			Description:  roleModels[i].Description,
			ResourceList: roleModels[i].ResourceListContent,
			Status:       roleModels[i].Status,
			AppId:        roleModels[i].AppId,
			CreatedAtTs:  uint32(roleModels[i].CreatedAt.Unix()),
			UpdatedAtTs:  uint32(roleModels[i].UpdatedAt.Unix()),
		}
		rolePbs = append(rolePbs, rolePb)
	}

	// get resources
	var appAgg pb.AdminAppAgg
	appAgg.Roles = rolePbs
	appAgg.AppId = appId

	// see if this user is banned or is a super admin
	appIdModel, err := s.adminUserBiz.GetUserAppRelation(ctx, user.Id, appId)
	if err != nil {
		return nil, err
	}
	var (
		isBanned uint32 = UserStatusNotBanned
		isSuper  uint32 = UserIdentityNotSuper
	)
	// add is banned to this user
	if appIdModel.Status == bo.UserAppIdBanned {
		isBanned = UserStatusBanned
	}
	// add super role to user for display if the user's identity is super
	if appIdModel.Identity == bo.IdentitySuper {
		roleList := make([]*pb.AdminRole, 0)
		roleList = append(roleList, &pb.AdminRole{
			Name:         "超级管理员",
			Description:  "超级管理员",
			ResourceList: bo.AllResources,
			Status:       bo.AdminRoleNormal,
			AppId:        appId,
		})
		if len(rolePbs) != 0 {
			roleList = append(roleList, rolePbs...)
		}
		rolePbs = roleList
		isSuper = UserIdentitySuper
	}

	return &pb.AdminUserListAgg{
		User: userPb,
		App: &pb.AdminAppAgg{
			AppId:    appId,
			Roles:    rolePbs,
			IsBanned: isBanned,
			IsSuper:  isSuper,
		},
	}, nil
}

// DiffRoles
// calculate relative complement of A array  and B array, AKA A - B
// for diff new roles and old roles, find elements A have but B do not
func diffRoles(aRoles []*bo.AdminRole, bRoles []*bo.AdminRole) []*bo.AdminRole {
	bMap := make(map[uint32]bool)
	for i := 0; i < len(bRoles); i++ {
		bMap[bRoles[i].Id] = true
	}
	result := make([]*bo.AdminRole, 0)
	for i := 0; i < len(aRoles); i++ {
		_, ok := bMap[aRoles[i].Id]
		if !ok {
			result = append(result, aRoles[i])
		}
	}
	return result
}
