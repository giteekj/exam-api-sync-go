package api

import (
	"exam-api-sync-go/common"
	"exam-api-sync-go/common/response"
	"exam-api-sync-go/common/setting"
	"exam-api-sync-go/model/request"
	"exam-api-sync-go/service"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// InventoryManualSync 库存记录手动同步接口
func InventoryManualSync(c *gin.Context) {
	// 解析请求参数
	var req request.InventorySyncQueryRequest
	body := common.GetRequestBody(c)
	if body == nil {
		response.ErrorResponse(c, common.FATAL, nil)
		return
	}
	req.Token = body.GetString("token", c)
	req.StartTime = body.GetString("start_time", c)
	req.EndTime = body.GetString("end_time", c)
	if req.StartTime == "" || req.EndTime == "" {
		response.ErrorResponse(c, common.LACK_OF_PARAMETER, nil)
		return
	}
	// 判断手动同步Token是否正确
	if req.Token != setting.Sync.Token {
		response.ErrorResponse(c, common.SYNC_TOKEN_ERROR, nil)
	}
	// 定义时间格式
	layout := "2006-01-02"
	loc, err := time.LoadLocation("Local")
	if err != nil {
		response.ErrorResponse(c, common.GET_TIME_ERROR, nil)
		return
	}
	parsedStartTime, err := time.ParseInLocation(layout, req.StartTime, loc)
	if err != nil {
		response.ErrorResponse(c, common.TIMEFORMAT_ERROR, nil)
		return
	}
	endTime, err := time.ParseInLocation(layout, req.EndTime, loc)
	if err != nil {
		response.ErrorResponse(c, common.TIMEFORMAT_ERROR, nil)
		return
	}
	parsedEndTime := endTime.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
	syncService := service.NewInventorySyncService()
	if msg, err := syncService.DoSyncInventory(parsedStartTime, parsedEndTime, "manual"); err != nil {
		log.Printf("库存手动同步失败: %v", err)
		response.ErrorResponse(c, common.FATAL, nil, err.Error())
	} else {
		log.Println("库存手动同步成功")
		data := map[string]string{
			"message": msg,
		}
		response.ErrorResponse(c, common.SUCCESS, data)
	}
}

// InventoryQuery 库存查询接口
func InventoryQuery(c *gin.Context) {
	// 解析请求参数
	var req request.InventoryQueryRequest
	body := common.GetRequestBody(c)
	if body == nil {
		response.ErrorResponse(c, common.FATAL, nil)
		return
	}
	req.Page = body.GetInt("page", c)
	req.PageSize = body.GetInt("page_size", c)
	req.Sort = body.GetString("sort", c)
	if req.Page < 0 || req.PageSize < 0 || req.Sort == "" {
		response.ErrorResponse(c, common.LACK_OF_PARAMETER, nil)
		return
	}
	req.StartTime = body.GetString("start_time", c)
	req.EndTime = body.GetString("end_time", c)
	if req.StartTime != "" && req.EndTime == "" {
		response.ErrorResponse(c, common.TIMEFORMAT_ERROR, nil)
		return
	}
	if req.StartTime == "" && req.EndTime != "" {
		response.ErrorResponse(c, common.TIMEFORMAT_ERROR, nil)
		return
	}
	req.Province = body.GetString("province", c)
	req.City = body.GetString("city", c)
	req.District = body.GetString("district", c)
	req.UserName = body.GetString("user_name", c)
	req.Phone = body.GetString("phone", c)
	req.BusinessFollower = body.GetString("business_follower", c)
	req.OperationFollower = body.GetString("operation_follower", c)
	req.UpperUid = body.GetInt("upper_uid", c)
	req.TopUid = body.GetInt("top_uid", c)

	// 获取用户信息
	userID, exists := c.Get("user_id")
	if !exists {
		response.ErrorResponse(c, common.LOGIN_EXPIRE, nil)
		return
	}
	// 获取用户权限信息
	role, exists := c.Get("roles")
	if !exists {
		response.ErrorResponse(c, common.ROLE_ERROR, nil)
	}
	roles, ok := role.([]string)
	if !ok {
		response.ErrorResponse(c, common.ROLE_ERROR, nil)
	}
	// 调用服务层处理查询
	inventoryService := service.NewInventoryService()
	if len(roles) == 0 {
		response.ErrorResponse(c, common.ROLE_ERROR, nil)
		return
	}
	roles = []string{
		"advancedQuery",
	}
	result, err := inventoryService.GetInventoryList(req, userID.(int), roles)
	if err != nil {
		response.ErrorResponse(c, common.SYSTEM_ERROR, nil)
		return
	}

	response.JsonResponse(c, result)
}
