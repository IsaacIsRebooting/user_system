package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"my_user_system/conf"
	"sync"
	"time"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func openDB() {
	mysqlConf := conf.GetGlobalConfig().DbConfig
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", mysqlConf.User,
		mysqlConf.Password, mysqlConf.Host, mysqlConf.Port, mysqlConf.Dbname)
	log.Info("mdb addr:" + connArgs)

	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("fetch db connecion err:" + err.Error())
	}

	sqlDB.SetMaxIdleConns(mysqlConf.MaxIdleConn)
	sqlDB.SetMaxOpenConns(mysqlConf.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Duration(mysqlConf.MaxIdleTime))
}

// 获取数据库连接
func GetDB() *gorm.DB {
	dbOnce.Do(openDB)
	return db
}
