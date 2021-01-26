package mysql

import (
	"github.com/tinysrc/z9go/pkg/conf"
	"github.com/tinysrc/z9go/pkg/log"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Db gorm
var Db *gorm.DB

func initConfig() {
	conf.Global.SetDefault("mysql.dsn", "root:123456@tcp(127.0.0.1:3306)/z9?charset=utf8")
}

func init() {
	dsn := conf.Global.GetString("mysql.dsn")
	db, err := gorm.Open(
		mysql.New(mysql.Config{
			DSN: dsn,
		}),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatal("mysql open failed", zap.Error(err))
		return
	}
	Db = db
	log.Info("mysql init success")
}
