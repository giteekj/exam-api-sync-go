package api

import (
	"exam-api-sync-go/common"
	"exam-api-sync-go/common/response"
	"exam-api-sync-go/model/request"
	"exam-api-sync-go/service"

	"github.com/gin-gonic/gin"
)

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
		response.ErrorResponse(c, common.FATAL, nil)
		return
	}
	req.StartTime = body.GetString("start_time", c)
	req.EndTime = body.GetString("end_time", c)
	req.Province = body.GetString("province", c)
	req.City = body.GetString("city", c)
	req.District = body.GetString("district", c)
	req.UserName = body.GetString("user_name", c)
	req.Phone = body.GetString("phone", c)
	req.BusinessFollower = body.GetString("business_follower", c)
	req.OperationFollower = body.GetString("operation_follower", c)

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
	result, err := inventoryService.GetInventoryList(req, userID.(int), roles)
	if err != nil {
		response.ErrorResponse(c, common.SYSTEM_ERROR, nil)
		return
	}

	response.JsonResponse(c, result)
}
