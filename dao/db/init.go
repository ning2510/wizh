package db

import (
	"fmt"
	"log"
	"os"
	"time"
	"wizh/pkg/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB     *gorm.DB
	config = viper.InitConf("dao")
)

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.Lshortfile),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Error,
			Colorful:      true,
		},
	)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Viper.GetString("mysql.username"),
		config.Viper.GetString("mysql.password"),
		config.Viper.GetString("mysql.host"),
		config.Viper.GetInt("mysql.port"),
		config.Viper.GetString("mysql.database"),
	)
	fmt.Println(dsn)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}

	// AutoMigrate 会创建表，缺失的外键，约束，列和索引。如果大小，精度，是否为空，可以更改，
	// 则 AutoMigrate 会改变列的类型。出于保护您数据的目的，它不会删除未使用的列刷新数据库
	// 的表格，使其保持最新。即如果我在旧表的基础上增加一个字段 age，那么调用 autoMigrate 后，
	// 旧表会自动多出一列 age，值为空
	// if err = DB.AutoMigrate(&User{}, &Video{}, &Comment{}, &FavoriteVideoRelation{}, &FavoriteCommentRelation{}); err != nil {
	// 	panic(err.Error())
	// }
	fmt.Println("mysql connected success!")
}
