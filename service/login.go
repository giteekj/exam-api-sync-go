package service

import (
	"errors"
	"exam-api-sync-go/common"
	"exam-api-sync-go/common/jwt"
	"exam-api-sync-go/common/orm"
	"exam-api-sync-go/model/request"
	"exam-api-sync-go/model/response"
	"strconv"
	"strings"

	daofxshop "exam-api-sync-go/dao/fxshop"
	modelfxshop "exam-api-sync-go/dao/model/fxshop"

	"gorm.io/gorm"
)

// LoginService 登录服务
type LoginService struct {
	DB *gorm.DB
}

// NewLoginService 创建登录服务实例
func NewLoginService() *LoginService {
	db := orm.GetDB("fxshop")
	if db != nil {
		daofxshop.SetDefault(db)
	}
	return &LoginService{
		DB: db,
	}
}

// Login 用户登录
func (s *LoginService) Login(req request.LoginRequest) (response.LoginResponse, error) {
	var resp response.LoginResponse
	// 查找用户
	u := daofxshop.Use(s.DB).SysUser
	user, err := u.Where(u.UserName.Eq(req.Username), u.UserPassword.Eq(common.MD5(req.Password)), u.UserStatus.Eq(1)).
		Or(u.UserTel.Eq(req.Username)).
		Select(u.UID, u.UserName, u.UserPassword).
		First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return resp, common.ReturnErr(common.LOGIN_FATAL)
		}
	}
	// 查找用户权限
	roles, err := s.GetUserRoles(user.UID)
	if err != nil {
		return resp, err
	}
	// 生成token
	token, err := jwt.GenerateToken(user.UID, user.UserName, roles)
	if err != nil {
		return resp, common.ReturnErr(common.GEN_TOEKN_ERROR)
	}

	// 构建响应
	resp.Token = token
	resp.Uid = user.UID
	resp.Name = user.UserName
	resp.Role = roles

	return resp, nil
}

// GetUserRoles 获取用户权限集合
func (s *LoginService) GetUserRoles(userID int) ([]string, error) {
	// 查询用户及后台用户
	useDB := daofxshop.Use(s.DB)
	u := useDB.SysUser
	ua := useDB.SysUserAdmin
	var data []response.SysUserInfo
	err := u.Where(u.UID.Eq(userID), ua.AdminStatus.Eq(1)).Join(ua, ua.UID.EqCol(u.UID)).Select(u.UID, ua.GroupIDArray).Scan(&data)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, common.ReturnErr(common.ROLE_ERROR)
		}
	}
	if len(data) == 0 {
		return nil, common.ReturnErr(common.ROLE_ERROR)
	}
	var groupIds []int
	for _, v := range data {
		if v.GroupIDArray == 0 {
			continue
		}
		groupIds = append(groupIds, v.GroupIDArray)
	}
	if len(groupIds) == 0 {
		return nil, common.ReturnErr(common.ROLE_ERROR)
	}

	// 查询用户组
	var groupData []modelfxshop.SysUserGroup
	ug := useDB.SysUserGroup
	err = ug.Where(ug.GroupID.In(groupIds...), ug.GroupStatus.Eq(1)).Select(ug.ModuleIDArray).Scan(&groupData)
	if err != nil {
		return nil, common.ReturnErr(common.ROLE_ERROR)
	}
	var (
		moduleIdSs []string
		moduleIds  []int
	)
	for _, v := range groupData {
		if v.ModuleIDArray == "" {
			continue
		}
		moduleIdList := strings.Split(v.ModuleIDArray, ",")
		for _, vv := range moduleIdList {
			moduleIdSs = append(moduleIdSs, vv)
		}
	}
	if len(moduleIdSs) == 0 {
		return nil, common.ReturnErr(common.ROLE_ERROR)
	}
	for _, v := range moduleIdSs {
		num, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		moduleIds = append(moduleIds, num)
	}

	// 查询模块
	var moduleData []modelfxshop.SysModule
	m := useDB.SysModule
	err = m.Where(m.ModuleID.In(moduleIds...)).Select(m.ModuleName).Scan(&moduleData)
	if err != nil {
		return nil, common.ReturnErr(common.ROLE_ERROR)
	}
	moduleMap := make(map[string]string)
	for _, v := range moduleData {
		if v.ModuleName == "" {
			continue
		}
		moduleMap[v.ModuleName] = v.ModuleName
	}
	var userRoles []string
	for _, v := range moduleMap {
		role, ok := common.QueryRole[v]
		if ok {
			userRoles = append(userRoles, role)
		}
	}
	if len(userRoles) == 0 {
		return nil, common.ReturnErr(common.ROLE_ERROR)
	}
	return userRoles, nil
}
