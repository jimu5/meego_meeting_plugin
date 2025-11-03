package dal

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"

	"meego_meeting_plugin/config"
	"meego_meeting_plugin/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)

	b, err := gorm.Open(getDialectorFromYamlConfig(), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	err = b.AutoMigrate(&model.PendingTask{})
	if err != nil {
		panic(err)
	}
	err = b.AutoMigrate(&model.CalendarBind{})
	if err != nil {
		panic(err)
	}
	err = b.AutoMigrate(&model.VCMeeting{})
	if err != nil {
		panic(err)
	}
	err = b.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	err = b.AutoMigrate(&model.VCMeetingUnBind{})
	if err != nil {
		panic(err)
	}
	err = b.AutoMigrate(&model.JoinChatRecord{})
	if err != nil {
		panic(err)
	}
	db = b
	return b
}

func getDialectorFromYamlConfig() gorm.Dialector {
	dbConfig := config.Config.Database
	var dsn string
	switch config.Config.Database.Type {
	case "sqlite":
		dsn = dbConfig.DBName + ".db"
		return sqlite.Open(dsn)
	case "mysql", "tidb":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
		return mysql.Open(dsn)
	case "postgresql":
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbConfig.Host, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.Port)
		return postgres.Open(dsn)
	case "sqlserver":
		dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
		return sqlserver.Open(dsn)
	default:
		// 未配置情况下使用 sqlite
		dsn = dbConfig.DBName + ".db"
		return sqlite.Open(dsn)
	}
}
