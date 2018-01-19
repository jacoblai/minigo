package engine

import (
	"github.com/jinzhu/gorm"
	"models"
	//_ "github.com/jinzhu/gorm/dialects/sqlite"
	 _ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/mysql"
	// _ "github.com/jinzhu/gorm/dialects/mssql"
)

type DbEngine struct {
	Engine       *gorm.DB   //关系型数据库引擎
	SystemConfig *ConfigYml //全局系统参数
}

func (d *DbEngine) Open(dir string) error {
	db, err := gorm.Open(d.SystemConfig.DbType, d.SystemConfig.DbConStr)
	if err != nil {
		return err
	}
	d.Engine = db
	d.Engine.AutoMigrate(&models.User{}, &models.Address{}, &models.CreditCard{}, &models.Email{}, &models.Language{})
	if d.SystemConfig.Debug == 1 {
		// 启用Logger，显示详细日志
		db.LogMode(true)
	}
	return nil
}

func (d *DbEngine) Close(){
	d.Engine.Close()
}
