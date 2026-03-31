package response

// LoginResponse 登录响应参数
type LoginResponse struct {
	Token string   `json:"token"`
	Uid   int      `json:"uid"`
	Name  string   `json:"name"`
	Role  []string `json:"role"`
}

type SysUserInfo struct {
	Uid          int `json:"uid"`
	GroupIDArray int `json:"group_id_array"`
}
