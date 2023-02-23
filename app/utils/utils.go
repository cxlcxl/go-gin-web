package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"strings"
	"time"
)

func IsRecordNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func MD5(params string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(params))
	return hex.EncodeToString(md5Ctx.Sum(nil))
}

// Base64 base64...
func Base64(params string) string {
	return base64.StdEncoding.EncodeToString([]byte(params))
}

// Shuffle 打乱数组原因顺序
func Shuffle(slice []interface{}) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

// GenerateSecret 生成密码加密串
func GenerateSecret(n int) string {
	if n == 0 {
		rand.Seed(time.Now().UnixNano())
		n = rand.Intn(15)
		if n < 3 {
			n = 8
		}
	}
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	bytes := []byte(str)
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		result[i] = bytes[rand.Intn(len(bytes))]
	}
	return string(result)
}

// Password 登陆密码
func Password(pass, secret string) string {
	return strings.ToUpper(MD5(base64.StdEncoding.EncodeToString([]byte(secret + pass + secret))))
}
