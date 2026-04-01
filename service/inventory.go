package service

import (
	"exam-api-sync-go/common"
	"exam-api-sync-go/common/orm"
	daofxshopSync "exam-api-sync-go/dao/fxshop_sync"
	req "exam-api-sync-go/model/request"
	resp "exam-api-sync-go/model/response"
	"fmt"
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
		//daofxshopSync.SetDefault(db)
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
		query.Where(air.Province.Like(fmt.Sprintf("%%%s%%", req.Province)))
	}
	if req.City != "" {
		query.Where(air.City.Like(fmt.Sprintf("%%%s%%", req.City)))
	}
	if req.District != "" {
		query.Where(air.District.Like(fmt.Sprintf("%%%s%%", req.District)))
	}
	if req.StartTime != "" && req.EndTime != "" {
		// 定义时间格式
		layout := "2006-01-02"
		// 解析日期时间字符串为time.Time对象，默认时间是00:00:00
		loc, err := time.LoadLocation("Local")
		if err != nil {
			return resp.InventoryQueryResponse{}, err
		}
		tStart, err := time.ParseInLocation(layout, req.StartTime, loc)
		if err != nil {
			return resp.InventoryQueryResponse{}, err
		}
		startTimestamp := tStart.Unix()
		tEnd, err := time.ParseInLocation(layout, req.EndTime, loc)
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
		// 非全国权限用户只能查询自己名下的数据(用户ID是自己、业务跟进人用户ID\运营跟进人用户ID是自己)
		query.Where(air.Where(air.UID.Eq(userID)).Or(air.BusinessFollowerUID.Eq(userID)).Or(air.OperationFollowerUID.Eq(userID)))
	}
	// 全国用户可查询所有数据
	if advancedQuery || nationalDataQuery {
		if req.UserName != "" {
			query.Where(air.UserName.Like(fmt.Sprintf("%%%s%%", req.UserName)))
		}
		if req.Phone != "" {
			query.Where(air.Phone.Eq(req.Phone))
		}
		if req.UpperUid != 0 {
			query.Where(air.UpperUID.Eq(req.UpperUid))
		}
		if req.TopUid != 0 {
			query.Where(air.TopUID.Eq(req.TopUid))
		}
		if req.BusinessFollower != "" {
			query.Where(air.BusinessFollower.Like(fmt.Sprintf("%%%s%%", req.BusinessFollower)))
		}
		if req.OperationFollower != "" {
			query.Where(air.OperationFollower.Like(fmt.Sprintf("%%%s%%", req.OperationFollower)))
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
	err = query.Debug().Offset(offset).Limit(req.PageSize).Scan(&list)
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
