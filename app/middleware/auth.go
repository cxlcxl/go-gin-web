package middleware

import (
	"gin-web/app/response"
	"gin-web/app/service/jwt"
	"gin-web/app/vars"
	"github.com/gin-gonic/gin"
	"strings"
)

type Headers struct {
	Authorization string `header:"Authorization"`
}

func CheckAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		headers := Headers{}
		if err := ctx.ShouldBindHeader(&headers); err != nil {
			response.Fail(ctx, "Token 信息读取失败")
			return
		}
		if headers.Authorization == "" {
			response.Fail(ctx, "No token information detected!")
			return
		}
		token := strings.Split(headers.Authorization, " ")
		if len(token) != 2 || len(token[1]) < 10 {
			response.Fail(ctx, "Token 信息有误")
			return
		}
		user, err := jwt.ParseUserToken(token[1])
		if err != nil {
			response.TokenExpired(ctx)
			return
		}
		loginUser := &vars.LoginUser{
			UserId:   user.Id,
			RoleId:   user.RoleId,
			Username: user.Username,
			Email:    user.Email,
			Mobile:   user.Mobile,
		}
		// 用户登陆信息，在控制器可以直接 get 获取
		ctx.Set(vars.LoginUserKey, loginUser)
		ctx.Next()
	}
}
