package service

import (
	"fmt"
	"log"
	"time"

	"exam-api-sync-go/common/orm"
	"exam-api-sync-go/model"

	"gorm.io/gorm"
)

// InventorySyncService 库存同步服务
type InventorySyncService struct {
	FxshopDB     *gorm.DB
	FxshopSyncDB *gorm.DB
}

// NewInventorySyncService 创建库存同步服务实例
func NewInventorySyncService() *InventorySyncService {
	return &InventorySyncService{
		FxshopDB:     orm.GetDB("fxshop"),
		FxshopSyncDB: orm.GetDB("fxshop_sync"),
	}
}

// SyncInventory 同步库存数据
func (s *InventorySyncService) SyncInventory() error {
	// 计算昨天的开始和结束时间
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	startOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 0, 0, 0, 0, yesterday.Location())
	endOfDay := time.Date(yesterday.Year(), yesterday.Month(), yesterday.Day(), 23, 59, 59, 999999999, yesterday.Location())

	log.Printf("开始同步 %s 到 %s 的库存数据", startOfDay.Format("2006-01-02 15:04:05"), endOfDay.Format("2006-01-02 15:04:05"))

	// 定义库存数据结构体
	type InventoryData struct {
		PurchasePrice    float64 `json:"purchase_price"`
		PurchaseQuantity float64 `json:"purchase_quantity"`
		Uid              int     `json:"uid"`
		UserName         string  `json:"user_name"`
		Phone            string  `json:"phone"`
		Balance          float64 `json:"balance"`
		Type             int     `json:"type"`
		Gift             float64 `json:"gift"`
		Consumed         float64 `json:"consumed"`
		Buy              float64 `json:"buy"`
		Receive          float64 `json:"receive"`
		Send             float64 `json:"send"`
		Back             float64 `json:"back"`
		//ActiveCodeConsumption float64 `json:"active_code_consumption"`
		//FullCodeConsumption   float64 `json:"full_code_consumption"`
		//SingleCodeConsumption float64 `json:"single_code_consumption"`
		Province          string `json:"province"`
		City              string `json:"city"`
		District          string `json:"district"`
		BusinessFollower  string `json:"business_follower"`
		OperationFollower string `json:"operation_follower"`
		UpperUid          int    `json:"upper_uid"`
		TopUid            int    `json:"top_uid"`
		CreateTime        int    `json:"create_time"`
	}

	// 查询库存数据（移除时间条件）
	var inventoryData []InventoryData
	query := `
SELECT 
su.uid,
su.user_name,
su.user_tel as phone,
nm.balance,
nm.type,
nm.gift,
nm.consumed,
nm.buy,
nm.receive,
nm.send,
nm.back,
COALESCE(var.province, '') as province,
COALESCE(var.city, '') as city,
COALESCE(var.district, '') as district,
COALESCE(me.tracker, '') as business_follower,
COALESCE(me.opt_tracker, '') as operation_follower,
COALESCE(me.partner_parent_id, 0) as upper_uid,
COALESCE(me.top_partner_id, 0) as top_uid,
nm.create_time,
COALESCE(acl.price, 0) as purchase_price,
COALESCE(acl.inventory, 0) as purchase_quantity
FROM 
sys_user_zgk_inventory nm
LEFT JOIN 
sys_user su ON nm.uid = su.uid
LEFT JOIN 
vsl_member me ON nm.uid = me.uid
LEFT JOIN 
vsl_agent_region var ON nm.uid = var.uid AND var.status < 2 AND var.delete_at = 0 AND var.store_status < 7 
LEFT JOIN
agent_contract_log acl ON nm.uid = acl.uid
`
	if err := s.FxshopDB.Raw(query).Scan(&inventoryData).Error; err != nil {
		return fmt.Errorf("查询库存数据失败: %v", err)
	}

	if len(inventoryData) == 0 {
		log.Println("没有需要同步的库存数据")
		return nil
	}

	// 判断全科码/单科码
	var (
		activeCodeConsumption float64 // 激活码消耗
		fullCodeConsumption   float64 // 全科码消耗
		singleCodeConsumption float64 // 单科码消耗
	)

	// 批量插入数据（使用UPSERT避免重复）
	var records []model.AgentInventoryRecord
	nowUnix := time.Now().Unix()

	for _, data := range inventoryData {
		if data.Type == 1 { // 1 中高考库存(单科码)
			singleCodeConsumption = data.Balance
		} else if data.Type == 2 { // 2 限时库存(全科码)
			fullCodeConsumption = data.Balance
		}
		activeCodeConsumption = activeCodeConsumption + singleCodeConsumption
		record := model.AgentInventoryRecord{
			Uid:                     data.Uid,
			UserName:                data.UserName,
			Phone:                   data.Phone,
			PurchasePrice:           data.PurchasePrice,
			PurchaseQuantity:        data.PurchaseQuantity,
			SelfPurchaseInventory:   data.Buy,
			UpperAllocation:         data.Receive,
			AllocatedToSubordinates: data.Send,
			RemainingInventory:      data.Balance,
			ActiveCodeConsumption:   activeCodeConsumption,
			FullCodeConsumption:     fullCodeConsumption,
			SingleCodeConsumption:   singleCodeConsumption,
			UpperUid:                data.UpperUid,
			RedeemedInventory:       data.Consumed,
			TopUid:                  data.TopUid,
			Province:                data.Province,
			City:                    data.City,
			District:                data.District,
			BusinessFollower:        data.BusinessFollower,
			OperationFollower:       data.OperationFollower,
			StoreTime:               data.CreateTime,
			CreatedTime:             int(nowUnix),
			UpdatedTime:             int(nowUnix),
		}
		records = append(records, record)
	}

	// 批量插入（使用事务确保数据一致性）
	err := s.FxshopSyncDB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(records); i += 100 {
			end := i + 100
			if end > len(records) {
				end = len(records)
			}
			batch := records[i:end]
			if err := tx.CreateInBatches(batch, 100).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("插入数据失败: %v", err)
	}

	log.Printf("库存同步完成，共同步 %d 条记录", len(records))
	return nil
}
