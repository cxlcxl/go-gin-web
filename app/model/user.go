package model

import (
	"gin-web/app/cache"
	"gin-web/app/utils"
	"gorm.io/gorm"
)

type User struct {
	gormDb

	Id       int64  `json:"id"`
	Email    string `json:"email"`                  // 登陆邮箱
	Username string `json:"username"`               // 用户名
	Mobile   string `json:"mobile"`                 // 手机号
	State    uint8  `json:"state"`                  // 账号状态
	RoleId   int64  `json:"role_id"`                // 角色ID
	Secret   string `json:"-" gorm:"column:secret"` // 密码加密符
	Pass     string `json:"-" gorm:"column:pass"`   // 密码
	Timestamp
}

var (
	userCacheKey = "db:user"
)

func (m *User) TableName() string {
	return "users"
}

func NewUser(db *gorm.DB) *User {
	return &User{gormDb: gormDb{DB: db}}
}

func (m *User) FindUserByEmail(email string) (user *User, err error) {
	err = m.Table(m.TableName()).Where("email = ?", email).First(&user).Error
	if utils.IsRecordNotFound(err) {
		return nil, nil
	}
	return
}

func (m *User) FindUserById(id int64) (user *User, err error) {
	// SetExpire()
	err = cache.New(m.DB).QueryRow(userCacheKey, &user, id, func(db *gorm.DB, v interface{}, id interface{}) error {
		return db.Table(m.TableName()).Where("id = ?", id).First(v).Error
	})
	return
}
