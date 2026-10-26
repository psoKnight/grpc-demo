package controller

import (
	"strings"

	"grpc-demo/generalgateway/service"
	"grpc-demo/pbdata/greeterpb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type HelloRequest struct {
	Name string `json:"name" binding:"required"`
}

func Hello(c *gin.Context) {
	var req HelloRequest
	if err := c.ShouldBindJSON(&req); err != nil || strings.TrimSpace(req.Name) == "" {
		WriteCodeMsgResponse(c, ParamError, "name 不能为空")
		return
	}

	res, err := service.SayHello(c.Request.Context(), req.Name)
	if err != nil {
		if status.Code(err) == codes.InvalidArgument {
			WriteCodeMsgResponse(c, ParamError, status.Convert(err).Message())
			return
		}
		WriteCodeMsgResponse(c, InternalError, err.Error())
		return
	}

	WriteCodeResponse(c, Success, &greeterpb.HelloReply{Message: res.GetMessage()})
}
