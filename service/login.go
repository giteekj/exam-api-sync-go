package service

import (
	"exam-api-sync-go/common"

	"exam-api-sync-go/common/jwt"
	"exam-api-sync-go/common/orm"
	"exam-api-sync-go/model"
	"exam-api-sync-go/model/request"
	"exam-api-sync-go/model/response"

	"gorm.io/gorm"
)

// LoginService 登录服务
type LoginService struct {
	DB *gorm.DB
}

// NewLoginService 创建登录服务实例
func NewLoginService() *LoginService {
	return &LoginService{
		DB: orm.GetDB("fxshop"),
	}
}

// Login 用户登录
func (s *LoginService) Login(req request.LoginRequest) (response.LoginResponse, error) {
	var response response.LoginResponse

	// 查找用户
	var user model.SysUser
	result := s.DB.Where("user_name = ? AND user_password = ?", req.Username, common.MD5(req.Password)).First(&user)
	if result.Error != nil {
		return response, common.ReturnErr(common.LOGIN_FATAL)
	}

	// 查找用户权限
	syncService := NewInventorySyncService()
	roles, err := syncService.GetUserRoles(user.Uid)
	if err != nil {
		return response, err
	}

	// 生成token
	token, err := jwt.GenerateToken(user.Uid, user.UserName, roles)
	if err != nil {
		return response, common.ReturnErr(common.GEN_TOEKN_ERROR)
	}

	// 构建响应
	response.Token = token
	response.Uid = user.Uid
	response.Name = user.UserName
	response.Role = roles

	return response, nil
}
