package dal

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"meego_meeting_plugin/model"
	"os"
	"time"
)

var db = initDB()

func initDB() *gorm.DB {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)
	b, err := gorm.Open(sqlite.Open("plugin.db"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	b.AutoMigrate(&model.CalendarBind{})
	b.AutoMigrate(&model.VCMeeting{})
	b.AutoMigrate(&model.User{})
	b.AutoMigrate(&model.VCMeetingUnBind{})
	b.AutoMigrate(&model.JoinChatRecord{})
	return b
}
