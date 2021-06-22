package di

import (
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/utils"
)

func (c *container) GetDatabase() *gorm.DB {
	if c.db == nil {
		_, err := utils.CreateLogFile(c.DbPath, 0744)

		if err != nil {
			log.Fatal().Err(err).Msg("Error while creating database folder")
		}

		enabled := c.viper.GetBool("logging.db.enabled")

		var logger gormlogger.Interface

		if enabled {
			logger = gormlogger.New(
				stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
				gormlogger.Config{
					SlowThreshold:             c.viper.GetDuration("logging.db.slow_query"),
					LogLevel:                  logging.ParseDBLogLevel(c.viper.GetString("logging.db.level")),
					IgnoreRecordNotFoundError: true,
					Colorful:                  c.viper.GetBool("logging.db.color"),
				},
			)
		} else {
			logger = nil
		}

		db, err := gorm.Open(sqlite.Open(c.DbPath), &gorm.Config{
			NowFunc:                                  time.Now().UTC,
			Logger:                                   logger,
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: false,
			DisableNestedTransaction:                 true,
		})

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error connecting to SQLite3 Database")
		}

		err = db.AutoMigrate(&models.User{}, &models.Message{})

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to migrate DATABASE")
		}

		c.db = db
	}

	return c.db
}
