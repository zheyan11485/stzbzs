package model

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var Conn *gorm.DB

func InitDB(databaseName string) {
	// 配置 GORM 日志级别，避免记录 "record not found" 等信息级日志
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,  // 慢查询阈值
			LogLevel:                  logger.Error, // 设置日志级别，只记录错误级别的日志
			IgnoreRecordNotFoundError: true,         // 忽略 RecordNotFound 错误
			Colorful:                  true,         // 启用彩色打印
		},
	)

	db, err := gorm.Open(sqlite.Open(databaseName+".db?cache=shared&mode=rwc"), &gorm.Config{
		Logger: newLogger, // 使用自定义日志配置
	})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&TeamUser{}, &Task{}, &Report{}, &BattleReport{}, &WuHistoryWeek{})
	if err != nil {
		return
	}

	Conn = db
}
