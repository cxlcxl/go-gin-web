package vars

import (
	ci "gin-web/library/config/interface"
	"gin-web/library/redis"
	"github.com/casbin/casbin/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

var (
	BasePath  string
	YmlConfig ci.YamlConfigInterface
	Casbin    *casbin.SyncedEnforcer
	GLog      *logrus.Logger
	DBMysql   *gorm.DB
	DBRedis   *redis.DBRedis
)

const (
	ApiPrefix       = "/api/v1"
	LoginUserKey    = "__gin_web_user__"
	ConfigKeyPrefix = "go-gin-web"
	DateFormat      = "2006-01-02"
	DateTimeFormat  = "2006-01-02 15:04:05"
)

func init() {
	if dir, err := os.Getwd(); err != nil {
		log.Fatal("文件目录获取失败")
		return
	} else {
		// 路径进行处理，兼容单元测试程序程序启动时的奇怪路径
		if len(os.Args) > 1 && strings.HasPrefix(os.Args[1], "-test") {
			BasePath = strings.Replace(strings.Replace(dir, `\test`, "", 1), `/test`, "", 1)
		} else {
			BasePath = dir
		}
	}
}
