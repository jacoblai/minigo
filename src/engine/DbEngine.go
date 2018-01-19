package engine

import (
	"github.com/jinzhu/gorm"
	"models"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DbEngine struct {
	Engine       *gorm.DB   //关系型数据库引擎
	SystemConfig *ConfigYml //全局系统参数
}

func (d *DbEngine) Open(dir string) error {
	connectionString := dir + "/test.db"
	db, err := gorm.Open("sqlite3", connectionString)
	if err != nil {
		return err
	}
	d.Engine = db
	d.Engine.AutoMigrate(&models.User{})
	if d.SystemConfig.Debug == 1 {
		// 启用Logger，显示详细日志
		db.LogMode(true)
	}
	return nil
}
