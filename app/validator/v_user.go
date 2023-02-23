package validator

import (
	"gin-web/app/handlers"
	"gin-web/app/validator/types"
	"github.com/gin-gonic/gin"
)

func (v BsValidator) VLogin(ctx *gin.Context) {
	var params types.Login
	bindData(ctx, &params, (handlers.HUser{}).Login)
}
