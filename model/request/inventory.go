package request

// InventoryQueryRequest 库存查询请求
type InventoryQueryRequest struct {
	StartTime         string `json:"start_time" form:"start_time"`                 // 开始时间
	EndTime           string `json:"end_time" form:"end_time"`                     // 结束时间
	Province          string `json:"province" form:"province"`                     // 省份
	City              string `json:"city" form:"city"`                             // 城市
	District          string `json:"district" form:"district"`                     // 区县
	UserName          string `json:"user_name" form:"user_name"`                   // 代理姓名
	Phone             string `json:"phone" form:"phone"`                           // 手机号
	UpperUid          string `json:"upper_uid" form:"upper_uid"`                   // 顶代ID
	TopUid            string `json:"top_uid" form:"top_uid"`                       // 上级ID
	BusinessFollower  string `json:"business_follower" form:"business_follower"`   // 业务跟进人
	OperationFollower string `json:"operation_follower" form:"operation_follower"` // 运营跟进人
	Page              int    `json:"page" form:"page"`                             // 页码
	PageSize          int    `json:"page_size" form:"page_size"`                   // 每页大小
	Sort              string `json:"sort" form:"sort"`                             // 排序
}
