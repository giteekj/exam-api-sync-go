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

// GetUserRoles 获取用户权限集合
func (s *InventoryService) GetUserRoles(userID int) ([]string, error) {
	var userRoles []model.AgentUserRole
	result := s.DB.Select("uid, role_id, agent_user_role.role_status, ar.role_code").
		Where("uid = ? AND agent_user_role.role_status = 1 AND ar.role_status = 1", userID).
		Joins("LEFT JOIN agent_role ar ON agent_user_role.role_id = ar.id").
		Find(&userRoles)
	if result.Error != nil {
		return nil, common.ReturnErr(common.ROLE_ERROR)
	}
	var roleCodes []string
	for _, userRole := range userRoles {
		roleCodes = append(roleCodes, userRole.RoleCode)
	}
	return roleCodes, nil
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
	// 构建查询
	query := s.DB.Model(&model.AgentInventoryRecord{})
	if req.Province != "" {
		query.Where("province = ?", req.Province)
	} else if req.City != "" {
		query.Where("city = ?", req.City)
	} else if req.District != "" {
		query.Where("district = ?", req.District)
	} else if req.StartTime != "" && req.EndTime != "" {
		query.Where("start_time >= ? AND start_time <= ?", req.StartTime, req.EndTime)
	}
	if realTimeQuery || advancedQuery {
		query.Where("uid = ?", userID)
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
