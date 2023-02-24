package gorm

import (
	"fmt"
	"gin-web/app/vars"
	glog "gin-web/library/logger"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GetMysqlClient 创建 MYSQL
func GetMysqlClient() (*gorm.DB, error) {
	dsn := getMysqlDsn()
	dia := mysql.Open(dsn)
	gormDb, err := gorm.Open(dia, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger: logger.New(&GormLog{}, logger.Config{
			SlowThreshold: time.Second * time.Duration(vars.YmlConfig.GetInt("Mysql.SlowQuery")),
			LogLevel:      logger.Warn,
		}),
	})
	if err != nil {
		//gorm 数据库驱动初始化失败
		return nil, err
	}

	// 为主连接设置连接池
	if rawDb, err := gormDb.DB(); err != nil {
		return nil, err
	} else {
		rawDb.SetConnMaxIdleTime(time.Second * 30)
		rawDb.SetConnMaxLifetime(time.Duration(vars.YmlConfig.GetInt("Mysql.SetConnMaxLifetime")) * time.Second)
		rawDb.SetMaxIdleConns(vars.YmlConfig.GetInt("Mysql.SetMaxIdleConns"))
		rawDb.SetMaxOpenConns(vars.YmlConfig.GetInt("Mysql.SetMaxOpenConns"))
		return gormDb, nil
	}
}

func getMysqlDsn() string {
	User := vars.YmlConfig.GetString("Mysql.Username")
	Pass := vars.YmlConfig.GetString("Mysql.Password")
	Host := vars.YmlConfig.GetString("Mysql.Host")
	Port := vars.YmlConfig.GetInt("Mysql.Port")
	DataBase := vars.YmlConfig.GetString("Mysql.Database")
	Charset := vars.YmlConfig.GetString("Mysql.Charset")
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local", User, Pass, Host, Port, DataBase, Charset)
}

type GormLog struct{}

func (gl *GormLog) Printf(format string, v ...interface{}) {
	info := fmt.Sprintf(format, v...)
	level := logrus.InfoLevel
	if strings.HasPrefix(format, "[error]") || strings.HasPrefix(format, "[traceErr]") {
		level = logrus.ErrorLevel
	} else if strings.HasPrefix(format, "[warn]") || strings.HasPrefix(format, "[traceWarn]") {
		level = logrus.WarnLevel
	}
	glog.NewLog(level, "mysql-gorm-log").Log(logrus.Fields{}, info)
	return
}
