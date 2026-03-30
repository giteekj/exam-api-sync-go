package common

var (
	NOT_EXISTS             = 1
	SUCCESS                = 1
	FATAL                  = -1
	REDIRECT               = 4
	LOGIN_EXPIRE           = -1000
	UPDATA_FAIL            = -1001
	DELETE_FAIL            = -1002
	LACK_OF_PARAMETER      = -1006
	SYSTEM_ERROR           = -1007
	ADD_FAIL               = -1008
	LOGIN_FATAL            = -2001
	USER_LOCK              = -2002
	USER_NBUND             = -2003
	PARAMETER_ERROR        = -1011
	ROLE_ERROR             = -1012
	INVALID_MOBILE         = 3001
	LOGIN_OUT_OF_LIMIT     = 3002
	SMS_CODE_OUT_OF_LIMIT  = 3003
	SMS_CODE_FATAL         = 3004
	TOO_MANY_REQUESTS      = 3005
	PLEASE_TRY_AGAIN_LATER = 3006

	PREFIX_OSS_FILE_PATH = "https://imagine-exam.oss-cn-hangzhou.aliyuncs.com/"
	PAGE_SIZE            = 20
	EMPTY_ARRAY          = []map[string]any{}
	EMPTY_MAP            = map[string]any{}
	EMPTY_STRING_ARRAY   = []string{}

	message = map[int]string{
		SUCCESS:                "ok",
		FATAL:                  "fatal",
		LACK_OF_PARAMETER:      "缺少参数",
		SYSTEM_ERROR:           "系统错误",
		ROLE_ERROR:             "未配置权限",
		LOGIN_FATAL:            "账号或者密码错误",
		USER_LOCK:              "用户被锁定",
		PARAMETER_ERROR:        "参数错误",
		LOGIN_EXPIRE:           "登录信息已过期，请重新登录",
		ADD_FAIL:               "添加失败",
		UPDATA_FAIL:            "修改失败",
		DELETE_FAIL:            "删除失败",
		INVALID_MOBILE:         "手机号格式不正确",
		LOGIN_OUT_OF_LIMIT:     "登录失败次数超出上限!",
		SMS_CODE_OUT_OF_LIMIT:  "验证码发送次数超出上限!",
		SMS_CODE_FATAL:         "验证码错误!",
		TOO_MANY_REQUESTS:      "请求过于频繁!",
		PLEASE_TRY_AGAIN_LATER: "请稍后再试!",
	}

	// 库存明细排序编号
	CONSUMEINVENTORYDESC = "consumeInventoryDesc" // 按消耗量由高到低
	PURCHASEQUANTITYDESC = "purchaseQuantityDesc" // 按进货量由高到低
	REMAININVENTORY      = "remainInventory"      // 按剩余库存由高到低

	// 库存查询权限
	REALTIMEQUERY     = "realTimeQuery"     // 实时库存查询权限
	ADVANCEDQUERY     = "advancedQuery"     // 高级库存查询权限
	NATIONALDATAQUERY = "nationalDataQuery" // 全国库存查询权限
)
