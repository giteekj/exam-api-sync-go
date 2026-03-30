package model

// AgentRole 代理角色表
type AgentRole struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	RoleCode    string `json:"role_code" gorm:"size:50;not null;index"`
	RoleName    string `json:"role_name" gorm:"size:50;not null;index"`
	RoleStatus  int8   `json:"role_status" gorm:"size:2;not null;default:1"` // 0: 禁用, 1: 启用
	CreatedTime int    `json:"created_time" gorm:"not null"`
	UpdatedTime *int   `json:"updated_time"`
}

// TableName 指定表名
func (AgentRole) TableName() string {
	return "agent_role"
}

// AgentUserRole 用户代理角色关联表
type AgentUserRole struct {
	ID          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Uid         int    `json:"uid" gorm:"not null;uniqueIndex:idx_uid_role_id"` // 代理用户ID
	RoleCode    string `json:"role_code" gorm:"size:50;not null;index"`
	RoleID      int    `json:"role_id" gorm:"not null;uniqueIndex:idx_uid_role_id"` // 角色ID
	RoleStatus  int8   `json:"role_status" gorm:"size:2;not null;default:1"`        // 0: 禁用, 1: 启用
	CreatedTime int    `json:"created_time" gorm:"not null"`
	UpdatedTime *int   `json:"updated_time"`
}

// TableName 指定表名
func (AgentUserRole) TableName() string {
	return "agent_user_role"
}

// AgentInventoryRecord 代理库存记录
type AgentInventoryRecord struct {
	ID                      int     `json:"id" gorm:"primaryKey;autoIncrement"`                                        // 主键ID
	Uid                     int     `json:"uid" gorm:"not null;index"`                                                 // 代理用户ID
	UserName                string  `json:"user_name" gorm:"size:50;not null;index"`                                   // 代理姓名
	Phone                   string  `json:"phone" gorm:"size:20;not null;index"`                                       // 手机号
	PurchasePrice           float64 `json:"purchase_price" gorm:"type:decimal(10,2);not null;default:0.00"`            // 进货价(元)
	PurchaseQuantity        float64 `json:"purchase_quantity" gorm:"type:decimal(10,2);not null;default:0.00"`         // 进货量
	SelfPurchaseInventory   float64 `json:"self_purchase_inventory" gorm:"type:decimal(10,2);not null;default:0.00"`   // 自购库存
	UpperAllocation         float64 `json:"upper_allocation" gorm:"type:decimal(10,2);not null;default:0.00"`          // 上级分配
	RedeemedInventory       float64 `json:"redeemed_inventory" gorm:"type:decimal(10,2);not null;default:0.00"`        // 已兑库存
	AllocatedToSubordinates float64 `json:"allocated_to_subordinates" gorm:"type:decimal(10,2);not null;default:0.00"` // 分配下级
	RemainingInventory      float64 `json:"remaining_inventory" gorm:"type:decimal(10,2);not null;default:0.00"`       // 剩余库存
	ActiveCodeConsumption   float64 `json:"active_code_consumption" gorm:"type:decimal(10,2);not null;default:0.00"`   // 激活码消耗
	FullCodeConsumption     float64 `json:"full_code_consumption" gorm:"type:decimal(10,2);not null;default:0.00"`     // 全科码消耗
	SingleCodeConsumption   float64 `json:"single_code_consumption" gorm:"type:decimal(10,2);not null;default:0.00"`   // 单科码消耗
	UpperUid                int     `json:"upper_uid" gorm:"not null;default:0"`                                       // 上级ID
	TopUid                  int     `json:"top_uid" gorm:"not null;default:0"`                                         // 顶代ID
	Province                string  `json:"province" gorm:"size:50"`                                                   // 省份
	City                    string  `json:"city" gorm:"size:50"`                                                       // 城市
	District                string  `json:"district" gorm:"size:50"`                                                   // 区县
	BusinessFollower        string  `json:"business_follower" gorm:"size:50"`                                          // 业务更进人
	OperationFollower       string  `json:"operation_follower" gorm:"size:50"`                                         // 运营更进人
	StoreTime               int     `json:"store_time" gorm:"not null;index"`                                          // 入库时间
	CreatedTime             int     `json:"created_time" gorm:"not null"`                                              // 创建时间
	UpdatedTime             int     `json:"updated_time"`                                                              // 更新时间
}

// TableName 指定表名
func (AgentInventoryRecord) TableName() string {
	return "agent_inventory_record"
}
