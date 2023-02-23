package handlers

import (
	"gin-web/app/response"
	serviceuser "gin-web/app/service/user"
	"gin-web/app/validator/types"
	"github.com/gin-gonic/gin"
)

type HUser struct{}

func (h HUser) Login(ctx *gin.Context, v interface{}) {
	params := v.(*types.Login)
	if token, err := serviceuser.Login(params); err != nil {
		response.Fail(ctx, err.Error())
	} else {
		response.Success(ctx, token)
	}
}
