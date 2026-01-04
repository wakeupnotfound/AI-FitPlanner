package errors

// 业务状态码定义
const (
	// 成功
	Success = 200

	// 客户端错误 (4000系列)
	ErrBadRequest       = 4000 // 请求错误
	ErrInvalidParam     = 4001 // 参数无效
	ErrUnauthorized     = 4010 // 未认证
	ErrForbidden        = 4030 // 无权限
	ErrNotFound         = 4040 // 资源不存在
	ErrMethodNotAllowed = 4050 // 方法不允许
	ErrConflict         = 4090 // 冲突

	// 服务器错误 (5000系列)
	ErrInternalServer  = 5000 // 内部错误
	ErrExternalService = 5001 // 外部服务错误
	ErrDatabase        = 5002 // 数据库错误
	ErrCache           = 5003 // 缓存错误

	// 业务错误 (6000系列)
	ErrUserExists         = 6001 // 用户已存在
	ErrUserNotFound       = 6002 // 用户不存在
	ErrWrongPassword      = 6003 // 密码错误
	ErrTokenExpired       = 6004 // Token过期
	ErrPlanNotFound       = 6005 // 计划不存在
	ErrAiApiNotConfigured = 6006 // AI API未配置
	ErrApiLimitExceeded   = 6007 // API调用超限
	ErrInvalidCredentials = 6008 // 无效的凭证
)
