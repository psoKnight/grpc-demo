package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	Success       = 0
	ParamError    = 1300
	InternalError = 1301
)

var errorText = map[int]string{
	Success:       "成功",
	ParamError:    "参数错误",
	InternalError: "内部错误",
}

type CodeResponse struct {
	ErrorNo  int         `json:"error_no"`
	ErrorMsg string      `json:"error_msg"`
	Data     interface{} `json:"data,omitempty"`
}

func WriteCodeResponse(c *gin.Context, code int, data interface{}) {
	res, _ := json.Marshal(CodeResponse{
		ErrorNo:  code,
		ErrorMsg: errorText[code],
		Data:     data,
	})
	c.String(http.StatusOK, string(res))
}

func WriteCodeMsgResponse(c *gin.Context, code int, msg string) {
	res, _ := json.Marshal(CodeResponse{
		ErrorNo:  code,
		ErrorMsg: msg,
	})
	c.String(http.StatusOK, string(res))
}
