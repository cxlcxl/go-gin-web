package gorm

import (
	"fmt"
	"gin-web/app/vars"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GetMysqlClient 创建 MYSQL
func GetMysqlClient() (*gorm.DB, error) {
	dsn := getMysqlDsn()
	dia := mysql.Open(dsn)
	gormDb, err := gorm.Open(dia, &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		//Logger:                 redefineLog(sqlType), // 拦截、接管 gorm v2 自带日志
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
