package errors

import "fmt"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("code=%d, message=%s, error=%v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("code=%d, message=%s", e.Code, e.Message)
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func Wrap(err error, code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// 常用错误
var (
	ErrUsernameExists   = New(ErrUserExists, "用户名已存在")
	ErrEmailExists      = New(ErrUserExists, "邮箱已存在")
	ErrInvalidUsername  = New(ErrInvalidParam, "用户名格式不正确")
	ErrInvalidEmail     = New(ErrInvalidParam, "邮箱格式不正确")
	ErrInvalidPassword  = New(ErrInvalidParam, "密码格式不正确")
	ErrPasswordMismatch = New(ErrInvalidParam, "两次输入的密码不一致")
	ErrUserDisabled     = New(ErrForbidden, "用户已被禁用")
	ErrTokenInvalid     = New(ErrUnauthorized, "无效的token")
	ErrSessionNotFound  = New(ErrUnauthorized, "会话不存在或已过期")
	ErrPermissionDenied = New(ErrForbidden, "权限不足")
	ErrResourceNotFound = New(ErrNotFound, "请求的资源不存在")
	ErrNoDefaultAIAPI   = New(ErrAiApiNotConfigured, "未设置默认的AI API")
	ErrAIApiTestFailed  = New(ErrExternalService, "AI API测试失败")
	ErrPlanGeneration   = New(ErrExternalService, "计划生成失败")
	ErrDuplicateRecord  = New(ErrConflict, "记录已存在")
)
