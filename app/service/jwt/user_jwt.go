package jwt

import (
	"errors"
	"gin-web/app/vars"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type UserTokenInfo struct {
	Id       int64  `json:"id"`
	RoleId   int64  `json:"role_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	jwt.StandardClaims
}

type Signed struct {
	Sign []byte
}

// CreateJwtToken 创建 token
func (s *Signed) CreateJwtToken(claims UserTokenInfo) (token string, err error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = t.SignedString(s.Sign)
	return
}

// CreateSign 签名
func CreateSign(s string) []byte {
	return []byte(s)
}

// ParseToken 解析token
func ParseToken(token string) (*UserTokenInfo, error) {
	t, err := jwt.ParseWithClaims(token, &UserTokenInfo{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(vars.YmlConfig.GetString("Token.Sign")), nil
	})
	if t == nil {
		return nil, errors.New("解析失败，Token 无效")
	}
	if err != nil {
		return nil, err
	}
	if user, ok := t.Claims.(*UserTokenInfo); ok && t.Valid {
		return user, nil
	} else {
		return nil, errors.New("解析失败，Token 无效")
	}
}

// CreateUserToken 创建用户登陆 token
func CreateUserToken(id, roleId int64, email, username, mobile string) (token string, err error) {
	userToken := UserTokenInfo{
		Id:       id,
		RoleId:   roleId,
		Username: username,
		Email:    email,
		Mobile:   mobile,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 10, // 生效开始时间
			ExpiresAt: time.Now().Unix() + vars.YmlConfig.GetInt64("Token.ExpiresAt"),
		},
	}
	sign := &Signed{
		Sign: CreateSign(vars.YmlConfig.GetString("Token.Sign")),
	}
	return sign.CreateJwtToken(userToken)
}

// ParseUserToken 解析用户 token
func ParseUserToken(token string) (info *UserTokenInfo, err error) {
	info, err = ParseToken(token)
	if err != nil {
		return nil, err
	}
	// 判断 token 在数据库中状态
	if checkTokenIsOk(info.Id, token) {
		return
	} else {
		return nil, errors.New("token expired")
	}
}

// CheckTokenIsOk 检查用户token是否有效
func checkTokenIsOk(id int64, token string) bool {
	return true
	// 在 缓存 中查询是否存在
	//suffix := utils.MD5(token)
	//if vars.DBRedis(db_cache.GetCacheKey(global.LoginToken, suffix)) {
	//	return true
	//} else {
	//}
}
