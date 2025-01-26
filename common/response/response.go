package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ResponseBase struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrCode struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (err *ErrCode) ReplaceMsg(message string) *ErrCode {
	err.Message = message
	return err
}

var (
	ParamError       = &ErrCode{Status: 6001, Message: "参数错误"}
	ParamParserError = &ErrCode{Status: 6002, Message: "参数转义失败"}
	SystemError      = &ErrCode{Status: 6003, Message: "内部系统错误"}
	NotData          = &ErrCode{Status: 6004, Message: "无可用数据"}
	DataError        = &ErrCode{Status: 6005, Message: "数据错误"}

	DataRepeatError       = &ErrCode{Status: 6006, Message: "数据重复错误"}
	ContractAppExistError = &ErrCode{Status: 6007, Message: "合同已经关联了大区"}

	UserAuthError = &ErrCode{Status: 6008, Message: "用户权限不足"}
)

//type SuccessResponse struct {
//	Status  int         `json:"status"`
//	Message string      `json:"message"`
//	Data    interface{} `json:"data"`
//}

func Response(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, ResponseBase{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// Success
// 成功的请求
func Success(ctx *gin.Context, data interface{}) {
	Response(ctx, 200, "success", data)
}

// Fail
// 失败的请求
func Fail(ctx *gin.Context, errcode *ErrCode) {
	Response(ctx, errcode.Status, errcode.Message, nil)
}
func FailWithData(ctx *gin.Context, errcode *ErrCode, data interface{}) {
	Response(ctx, errcode.Status, errcode.Message, data)
}
