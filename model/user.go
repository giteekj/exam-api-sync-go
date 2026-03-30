package model

// SysUser 用户表
type SysUser struct {
	Uid        int    `json:"uid" gorm:"primaryKey"`
	UserName   string `json:"user_name" gorm:"size:50;not null;index"`
	UserTel    string `json:"user_tel" gorm:"size:20;not null;index"`
	UserPwd    string `json:"user_pwd" gorm:"size:50;not null"`
	ProvinceID int    `json:"province_id" gorm:"not null;index"`
	CityID     int    `json:"city_id" gorm:"not null;index"`
	DistrictID int    `json:"district_id" gorm:"not null;index"`
	Status     int    `json:"status" gorm:"not null;default:1"`
	CreatedAt  int64  `json:"created_at" gorm:"not null"`
	UpdatedAt  int64  `json:"updated_at" gorm:"not null"`
}

// TableName 指定表名
func (SysUser) TableName() string {
	return "sys_user"
}

// SysUserAdmin 后台用户表
type SysUserAdmin struct {
	Uid          int    `json:"uid" gorm:"primaryKey"`
	GroupIdArray string `json:"group_id_array" gorm:"size:255;not null"`
	Role         string `json:"role" gorm:"size:50;not null;index"`
}

// TableName 指定表名
func (SysUserAdmin) TableName() string {
	return "sys_user_admin"
}

// SysUserGroup 用户组表
type SysUserGroup struct {
	GroupId       int    `json:"group_id" gorm:"primaryKey"`
	GroupName     string `json:"group_name" gorm:"size:50;not null;index"`
	ModuleIdArray string `json:"module_id_array" gorm:"size:255;not null"`
	Status        int    `json:"status" gorm:"not null;default:1"`
}

// TableName 指定表名
func (SysUserGroup) TableName() string {
	return "sys_user_group"
}

// SysUserZgkInventory 库存表
type SysUserZgkInventory struct {
	ID          int     `json:"id" gorm:"primaryKey;autoIncrement"`
	Uid         int     `json:"uid" gorm:"not null;index"`
	Balance     float64 `json:"balance" gorm:"type:decimal(10,2);not null"`
	Gift        float64 `json:"gift" gorm:"type:decimal(10,2);not null"`
	Consumed    float64 `json:"consumed" gorm:"type:decimal(10,2);not null"`
	Buy         float64 `json:"buy" gorm:"type:decimal(10,2);not null"`
	Receive     float64 `json:"receive" gorm:"type:decimal(10,2);not null"`
	Send        float64 `json:"send" gorm:"type:decimal(10,2);not null"`
	Back        float64 `json:"back" gorm:"type:decimal(10,2);not null"`
	CreatedTime int     `json:"created_time" gorm:"not null;index"`
}

// TableName 指定表名
func (SysUserZgkInventory) TableName() string {
	return "sys_user_zgk_inventory"
}
