package repository

import (
	"crypto/tls"
	"time"

	driverMysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"qr-file-hub/backend/internal/config"
	"qr-file-hub/backend/internal/model"
)

var DB *gorm.DB

func InitDB(cfg *config.Config) (*gorm.DB, error) {
	_ = driverMysql.RegisterTLSConfig("tidb", &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	})

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(mysql.Open(cfg.DatabaseDSN), gormConfig)
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := db.AutoMigrate(&model.FileAsset{}, &model.Bridge{}, &model.AccessLog{}); err != nil {
		return nil, err
	}

	var bridgeCount int64
	db.Model(&model.Bridge{}).Count(&bridgeCount)
	if bridgeCount == 0 {
		now := time.Now()
		defaultBridges := []model.Bridge{
			{
				Domain:        "https://kakashi-ace.github.io",
				IsActive:      true,
				IsPrimary:     true,
				IsCustom:      false,
				WeChatStatus:  "normal",
				QQStatus:      "normal",
				LatencyMs:     68,
				LastCheckedAt: &now,
			},
			{
				Domain:        "https://kakashi-ace.github.io/bridge-v4zahz/",
				IsActive:      true,
				IsPrimary:     false,
				IsCustom:      false,
				WeChatStatus:  "normal",
				QQStatus:      "normal",
				LatencyMs:     75,
				LastCheckedAt: &now,
			},
		}
		for _, b := range defaultBridges {
			db.Create(&b)
		}
	}

	DB = db
	return db, nil
}
