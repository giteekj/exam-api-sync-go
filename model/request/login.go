package request

// LoginRequest 登录请求参数
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	//DeviceID string `json:"device_id" binding:"required"`
	//Client   string `json:"client" binding:"required"`
}
