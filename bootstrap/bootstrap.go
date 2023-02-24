package bootstrap

import (
	validator2 "gin-web/app/validator"
	"gin-web/app/vars"
	"gin-web/library/config"
	mg "gin-web/library/gorm"
	"gin-web/library/logger"
	"gin-web/library/redis"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"log"
	"os"
	"time"
)

func init() {
	checkConfigFiles()

	// 初始化 WEB 配置文件「配置文件在最开始加载，很多初始化依赖配置信息」
	vars.YmlConfig = config.CreateYamlFactory()
	vars.YmlConfig.ConfigFileChangeListen()

	initDatabase()
	initRedis()
	initCasbin()

	// 创建软连接、更好的管理静态资源
	initFoldersLink()

	// 初始日志
	logger.NewGLog()

	// 初始化验证器语言
	if err := validator2.LoadValidatorLocal(); err != nil {
		log.Fatal(err.Error())
	}
	// 注册自定义验证器
	validator2.RegisterValidators()
}

// 检查必要的配置文件
func checkConfigFiles() {
	if _, err := os.Stat(vars.BasePath + "/config/web.yaml"); err != nil {
		log.Fatal("请检查 WEB 配置文件是否存在：", err)
		return
	}
}

// 初始化数据库连接
func initDatabase() {
	if dbMysql, err := mg.GetMysqlClient(); err != nil {
		log.Fatal("数据库启动失败：", err)
	} else {
		vars.DBMysql = dbMysql
	}
}

// 初始化 redis 缓存
func initRedis() {
	host := vars.YmlConfig.GetString("Redis.Host")
	pass := vars.YmlConfig.GetString("Redis.Password")
	prefix := vars.YmlConfig.GetString("Redis.KeyPrefix")
	db := vars.YmlConfig.GetInt("Redis.Db")
	if dbRedis, err := redis.GetRedisInstance(host, pass, prefix, db); err != nil {
		log.Fatal("Redis 缓存实例创建失败：", err.Error())
	} else {
		vars.DBRedis = dbRedis
	}
}

func initCasbin() {
	if _, err := os.Stat(vars.BasePath + "/config/casbin.conf"); err != nil {
		log.Fatal("Casbin 认证配置文件不存在：", vars.BasePath+"/config/casbin.conf")
		return
	}
	tableName := vars.YmlConfig.GetString("Casbin.TableName")
	adapter, err := gormadapter.NewAdapterByDBUseTableName(vars.DBMysql, "", tableName)
	if err != nil {
		log.Fatal("初始化 Casbin 适配器（Adapter）失败：", err)
		return
	}

	enforcer, err := casbin.NewSyncedEnforcer(vars.BasePath+"/config/casbin.conf", adapter)
	if err != nil {
		log.Fatal("初始化 Casbin 执行器（Enforcer）失败：", err)
		return
	}

	if err = enforcer.LoadPolicy(); err != nil {
		log.Fatal("Casbin LoadPolicy 失败：", err)
	}
	// 每 20s 刷新验证规则
	enforcer.StartAutoLoadPolicy(time.Second * time.Duration(vars.YmlConfig.GetInt("Casbin.RefreshTime")))
	vars.Casbin = enforcer
}

func initFoldersLink() {
	// 创建软连接、更好的管理静态资源
	_ = os.Remove(vars.BasePath + "/web/storage")
	if err := os.Symlink(vars.BasePath+"/storage/app", vars.BasePath+"/web/storage"); err != nil {
		log.Fatal(err.Error())
	}
}
