package service

import (
	"exam-api-sync-go/common"
	"exam-api-sync-go/common/orm"
	"exam-api-sync-go/model"
	req "exam-api-sync-go/model/request"
	resp "exam-api-sync-go/model/response"

	"gorm.io/gorm"
)

// InventoryService 库存服务
type InventoryService struct {
	DB *gorm.DB
}

// NewInventoryService 创建库存服务实例
func NewInventoryService() *InventoryService {
	return &InventoryService{
		DB: orm.GetDB("fxshop_sync"),
	}
}

// GetInventoryList 获取库存列表
func (s *InventoryService) GetInventoryList(req req.InventoryQueryRequest, userID int, userRole []string) (resp.InventoryQueryResponse, error) {
	var total int64
	var list []resp.InventoryRecordItem

	// 判断用户是否具有实时库存查询权限
	realTimeQuery := common.StrContains(userRole, common.REALTIMEQUERY)
	// 判断用户是否具有高级库存查询权限
	advancedQuery := common.StrContains(userRole, common.ADVANCEDQUERY)
	// 判断用户是否具有全国库存数据查询权限
	nationalDataQuery := common.StrContains(userRole, common.NATIONALDATAQUERY)
	// 未配置权限
	if !realTimeQuery && !advancedQuery && !nationalDataQuery {
		return resp.InventoryQueryResponse{}, common.ReturnErr(common.ROLE_ERROR)
	}
	// 构建查询
	query := s.DB.Model(&model.AgentInventoryRecord{}).Where("status = ?", 1)
	if req.Province != "" {
		query.Where("province = ?", req.Province)
	} else if req.City != "" {
		query.Where("city = ?", req.City)
	} else if req.District != "" {
		query.Where("district = ?", req.District)
	} else if req.StartTime != "" && req.EndTime != "" {
		query.Where("start_time >= ? AND start_time <= ?", req.StartTime, req.EndTime)
	}
	// 优化权限逻辑
	if !nationalDataQuery {
		// 非全国权限用户只能查询自己的数据
		query = query.Where("uid = ?", userID)
	} else {
		// 全国权限用户可以查询全国数据
		query.Or("uid = ?", userID).Or("upper_uid = ?", userID).Or("top_uid = ?", userID)
	}
	if advancedQuery || nationalDataQuery {
		if req.UserName != "" {
			query.Where("user_name = ?", req.UserName)
		} else if req.Phone != "" {
			query.Where("phone = ?", req.Phone)
		} else if req.UpperUid != "" {
			query.Where("upper_uid = ?", req.UpperUid)
		} else if req.TopUid != "" {
			query.Where("top_uid = ?", req.TopUid)
		} else if req.BusinessFollower != "" {
			query.Where("business_follower = ?", req.BusinessFollower)
		} else if req.OperationFollower != "" {
			query.Where("operation_follower = ?", req.OperationFollower)
		}
	}
	// 总数
	query.Count(&total)
	if req.Sort == common.CONSUMEINVENTORYDESC { // 按消耗量由高到低
		query.Order("full_code_consumption desc")
	} else if req.Sort == common.PURCHASEQUANTITYDESC { // 按进货量由高到低
		query.Order("purchase_quantity desc")
	} else if req.Sort == common.REMAININVENTORY { // 按剩余库存由高到低
		query.Order("remaining_inventory desc")
	}
	// 分页
	offset := (req.Page - 1) * req.PageSize
	query.Offset(offset).Limit(req.PageSize).Find(&list)
	var (
		totalPurchaseQuantity float64
		totalConsumeInventory float64
		totalRemainInventory  float64
	)
	for _, v := range list {
		totalPurchaseQuantity += v.PurchaseQuantity
		totalConsumeInventory += v.FullCodeConsumption
		totalRemainInventory += v.RemainingInventory
	}
	return resp.InventoryQueryResponse{
		Total:                 total,
		TotalPurchaseQuantity: totalPurchaseQuantity,
		TotalConsumeInventory: totalConsumeInventory,
		TotalRemainInventory:  totalRemainInventory,
		List:                  list,
	}, nil
}
