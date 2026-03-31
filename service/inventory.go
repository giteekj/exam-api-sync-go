package service

import (
	"exam-api-sync-go/common"
	"exam-api-sync-go/common/orm"
	daofxshopSync "exam-api-sync-go/dao/fxshop_sync"
	req "exam-api-sync-go/model/request"
	resp "exam-api-sync-go/model/response"
	"time"

	"gorm.io/gorm"
)

// InventoryService 库存服务
type InventoryService struct {
	DB *gorm.DB
}

// NewInventoryService 创建库存服务实例
func NewInventoryService() *InventoryService {
	db := orm.GetDB("fxshop_sync")
	if db != nil {
		daofxshopSync.SetDefault(db)
	}
	return &InventoryService{
		DB: db,
	}
}

// GetInventoryList 获取库存列表
func (s *InventoryService) GetInventoryList(req req.InventoryQueryRequest, userID int, userRole []string) (resp.InventoryQueryResponse, error) {
	var total int64

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
	air := daofxshopSync.Use(s.DB).AgentInventoryRecord
	query := air.Where(air.Status.Eq(1))
	if req.Province != "" {
		query.Where(air.Province.Eq(req.Province))
	} else if req.City != "" {
		query.Where(air.City.Eq(req.City))
	} else if req.District != "" {
		query.Where(air.District.Eq(req.District))
	} else if req.StartTime != "" && req.EndTime != "" {
		// 解析日期时间字符串为time.Time对象，默认时间是00:00:00
		tStart, err := time.Parse("2006-01-02", req.StartTime)
		if err != nil {
			return resp.InventoryQueryResponse{}, err
		}
		startTimestamp := tStart.Unix()
		tEnd, err := time.Parse("2006-01-02", req.EndTime)
		if err != nil {
			return resp.InventoryQueryResponse{}, err
		}
		// 添加23小时59分钟59秒
		tEnd = tEnd.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
		endTimestamp := tEnd.Unix()
		query.Where(air.StoreTime.Gte(int(startTimestamp)), air.StoreTime.Lte(int(endTimestamp)))
	}
	// 优化权限逻辑
	if !nationalDataQuery {
		// 非全国权限用户只能查询自己的数据
		query = query.Where(air.UID.Eq(userID))
	} else {
		// 全国权限用户可以查询全国数据
		query.Or(air.UID.Eq(userID)).Or(air.UpperUID.Eq(userID)).Or(air.TopUID.Eq(userID))
	}
	if advancedQuery || nationalDataQuery {
		if req.UserName != "" {
			query.Where(air.UserName.Eq(req.UserName))
		} else if req.Phone != "" {
			query.Where(air.Phone.Eq(req.Phone))
		} else if req.UpperUid != 0 {
			query.Where(air.UpperUID.Eq(req.UpperUid))
		} else if req.TopUid != 0 {
			query.Where(air.TopUID.Eq(req.TopUid))
		} else if req.BusinessFollower != "" {
			query.Where(air.BusinessFollower.Eq(req.BusinessFollower))
		} else if req.OperationFollower != "" {
			query.Where(air.OperationFollower.Eq(req.OperationFollower))
		}
	}
	// 总数
	total, err := query.Count()
	if err != nil {
		return resp.InventoryQueryResponse{}, err
	}
	if req.Sort == common.CONSUMEINVENTORYDESC { // 按消耗量由高到低
		query.Order(air.FullCodeConsumption.Desc())
	} else if req.Sort == common.PURCHASEQUANTITYDESC { // 按进货量由高到低
		query.Order(air.PurchaseQuantity.Desc())
	} else if req.Sort == common.REMAININVENTORY { // 按剩余库存由高到低
		query.Order(air.RemainingInventory.Desc())
	}
	// 分页
	offset := (req.Page - 1) * req.PageSize
	var list []resp.InventoryRecordItem
	err = query.Offset(offset).Limit(req.PageSize).Scan(&list)
	if err != nil {
		return resp.InventoryQueryResponse{}, err
	}
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
