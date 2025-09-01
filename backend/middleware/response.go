package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Total     int64       `json:"total,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

type ErrorCode int

const (
	Success ErrorCode = iota
	InvalidParams
	DatabaseError
	NotFound
	InternalError
	ValidationError
)

var codeMessages = map[ErrorCode]string{
	Success:         "操作成功",
	InvalidParams:   "参数错误",
	DatabaseError:   "数据库操作失败",
	NotFound:        "资源未找到",
	InternalError:   "服务器内部错误",
	ValidationError: "数据验证失败",
}

func SuccessResponse(c *gin.Context, data interface{}) {
	response := Response{
		Code:      int(Success),
		Message:   codeMessages[Success],
		Data:      data,
		Timestamp: time.Now().Unix(),
	}
	c.JSON(http.StatusOK, response)
}

func SuccessWithTotal(c *gin.Context, data interface{}, total int64) {
	response := Response{
		Code:      int(Success),
		Message:   codeMessages[Success],
		Data:      data,
		Total:     total,
		Timestamp: time.Now().Unix(),
	}
	c.JSON(http.StatusOK, response)
}

func ErrorResponse(c *gin.Context, code ErrorCode, message string) {
	if message == "" {
		message = codeMessages[code]
	}
	
	response := Response{
		Code:      int(code),
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	
	var httpStatus int
	switch code {
	case NotFound:
		httpStatus = http.StatusNotFound
	case InvalidParams, ValidationError:
		httpStatus = http.StatusBadRequest
	default:
		httpStatus = http.StatusInternalServerError
	}
	
	c.JSON(httpStatus, response)
}

func ErrorHandler() gin.HandlerFunc {
	return gin.Recovery()
}