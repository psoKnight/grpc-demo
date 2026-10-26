package controller

import (
	"encoding/json"
	"strconv"
	"strings"

	"grpc-demo/generalgateway/service"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CalculateRequest struct {
	Operation string          `json:"operation"`
	Left      json.RawMessage `json:"left"`
	Right     json.RawMessage `json:"right"`
}

func Calculate(c *gin.Context) {
	var req CalculateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		WriteCodeMsgResponse(c, ParamError, "请求体格式错误")
		return
	}

	operation := strings.ToLower(strings.TrimSpace(req.Operation))
	if operation == "" {
		WriteCodeMsgResponse(c, ParamError, "operation 不能为空")
		return
	}
	if operation != service.OperationAdd && operation != service.OperationMultiply {
		WriteCodeMsgResponse(c, ParamError, "operation 仅支持 add 或 multiply")
		return
	}

	left, err := parseInt64Field(req.Left)
	if err != nil {
		WriteCodeMsgResponse(c, ParamError, "left 必须是 int64 整数")
		return
	}
	right, err := parseInt64Field(req.Right)
	if err != nil {
		WriteCodeMsgResponse(c, ParamError, "right 必须是 int64 整数")
		return
	}

	res, err := service.Calculate(c.Request.Context(), operation, left, right)
	if err != nil {
		switch status.Code(err) {
		case codes.InvalidArgument, codes.OutOfRange:
			WriteCodeMsgResponse(c, ParamError, status.Convert(err).Message())
		default:
			WriteCodeMsgResponse(c, InternalError, err.Error())
		}
		return
	}

	WriteCodeResponse(c, Success, res)
}

func parseInt64Field(raw json.RawMessage) (int64, error) {
	value := strings.TrimSpace(string(raw))
	if value == "" {
		return 0, strconv.ErrSyntax
	}
	if value[0] == '"' {
		if err := json.Unmarshal(raw, &value); err != nil {
			return 0, err
		}
	}
	return strconv.ParseInt(strings.TrimSpace(value), 10, 64)
}
