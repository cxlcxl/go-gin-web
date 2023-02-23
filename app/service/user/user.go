package serviceuser

import (
	"gin-web/app/errs"
	"gin-web/app/model"
	"gin-web/app/service/jwt"
	"gin-web/app/utils"
	"gin-web/app/validator/types"
	"gin-web/app/vars"
)

func Login(l *types.Login) (string, error) {
	user, err := model.NewUser(vars.DBMysql).FindUserByEmail(l.Email)
	if err != nil {
		return "", err
	}
	if user.State != 1 {
		return "", errs.LoginUserVoid
	}
	if user.Pass != utils.Password(l.Pass, user.Secret) {
		return "", errs.LoginPassErr
	}
	token, err := jwt.CreateUserToken(user.Id, user.RoleId, user.Email, user.Username, user.Mobile)
	if err != nil {
		return "", errs.LoginFailToken
	}
	if _, err = jwt.ParseUserToken(token); err != nil {
		return "", errs.LoginFailToken
	}
	return token, nil
}
