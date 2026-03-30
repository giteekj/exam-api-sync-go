package response

// InventoryQueryResponse 库存查询响应
type InventoryQueryResponse struct {
	Total                 int64                 `json:"total"`                   // 总记录数
	TotalPurchaseQuantity float64               `json:"total_purchase_quantity"` // 累计进货量
	TotalConsumeInventory float64               `json:"total_consume_inventory"` // 累计消耗库存
	TotalRemainInventory  float64               `json:"total_remain_inventory"`  // 剩余库存
	List                  []InventoryRecordItem `json:"list"`                    // 记录列表
}

// InventoryRecordItem 库存记录项
type InventoryRecordItem struct {
	ID                      int     `json:"id"`
	UserName                string  `json:"user_name"`
	Phone                   string  `json:"phone"`
	PurchasePrice           float64 `json:"purchase_price"`
	PurchaseQuantity        float64 `json:"purchase_quantity"`
	SelfPurchaseInventory   float64 `json:"self_purchase_inventory"`
	UpperAllocation         float64 `json:"upper_allocation"`
	RedeemedInventory       float64 `json:"redeemed_inventory"`
	AllocatedToSubordinates float64 `json:"allocated_to_subordinates"`
	RemainingInventory      float64 `json:"remaining_inventory"`
	ActiveCodeConsumption   float64 `json:"active_code_consumption"`
	FullCodeConsumption     float64 `json:"full_code_consumption"`
	SingleCodeConsumption   float64 `json:"single_code_consumption"`
	Province                string  `json:"province"`
	City                    string  `json:"city"`
	District                string  `json:"district"`
	BusinessFollower        string  `json:"business_follower"`
	OperationFollower       string  `json:"operation_follower"`
	StoreTime               int     `json:"store_time"`
}
