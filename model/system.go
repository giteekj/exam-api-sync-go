package model

// AgentInventoryRecordSync 代理库存同步记录表
type AgentInventoryRecordSync struct {
	ID             int    `json:"id" gorm:"primaryKey;autoIncrement"`
	LastSyncTime   int    `json:"last_sync_time" gorm:"size:11;not null;index"`
	SyncCount      int    `json:"sync_count" gorm:"size:11;not null;index"`
	SyncType       string `json:"sync_type" gorm:"size:20;not null;index"`
	SyncStatus     string `json:"sync_status" gorm:"size:30;not null"`
	QuerySql       string `json:"query_sql" gorm:"not null"`
	QueryStartTime int    `json:"query_start_time" gorm:"size:11;not null;index"`
	QueryEndTime   int    `json:"query_end_time" gorm:"size:11;not null;index"`
	CreatedTime    int    `json:"created_time" gorm:"size:11;not null"`
	UpdatedTime    int    `json:"updated_time"`
}

// TableName 指定表名
func (AgentInventoryRecordSync) TableName() string {
	return "agent_inventory_record_sync"
}
