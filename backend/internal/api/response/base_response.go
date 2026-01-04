package response

import (
	"time"
)

type BaseResponse struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

func Success(data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:      200,
		Message:   "success",
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}

func Error(code int, message string) *BaseResponse {
	return &BaseResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
}

// 常用响应函数
func UnauthorizedError(message string) *BaseResponse {
	return Error(4010, message)
}

func BadRequestError(message string) *BaseResponse {
	return Error(4000, message)
}

func NotFoundError(message string) *BaseResponse {
	return Error(4040, message)
}

func InternalServerError(message string) *BaseResponse {
	return Error(5000, message)
}

func ForbiddenError(message string) *BaseResponse {
	return Error(4030, message)
}

// 带数据的错误响应
func ErrorWithData(code int, message string, data interface{}) *BaseResponse {
	return &BaseResponse{
		Code:      code,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
}
