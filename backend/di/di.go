package di

import (
	stdlog "log"
	"os"
	"time"

	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	gormlogger "gorm.io/gorm/logger"

	"github.com/BrosSquad/ts-1-chat-app/backend/logging"
	"github.com/BrosSquad/ts-1-chat-app/backend/models"
	"github.com/BrosSquad/ts-1-chat-app/backend/utils"
)

type (
	Config struct {
		LogsRoot     string
		DbPath       string
		LogToConsole bool
	}

	Container interface {
		GetDatabase() *gorm.DB

		GetDebugLogger() *logging.Debug
		GetErrorLogger() *logging.Error
		GetInfoLogger() *logging.Info
	}

	container struct {
		Config

		db          *gorm.DB
		errorLogger *logging.Error
		debugLogger *logging.Debug
		infoLogger  *logging.Info
	}
)

func (c *container) GetDebugLogger() *logging.Debug {
	if c.debugLogger == nil {
		c.debugLogger = logging.NewDebugLogger(c.LogsRoot)
	}

	return c.debugLogger
}

func (c *container) GetErrorLogger() *logging.Error {
	if c.errorLogger == nil {
		c.errorLogger = logging.NewErrorLogger(c.LogsRoot, c.LogToConsole)
	}

	return c.errorLogger

}

func (c *container) GetInfoLogger() *logging.Info {
	if c.infoLogger == nil {
		c.infoLogger = logging.NewInfoLogger(c.LogsRoot, c.LogToConsole)
	}

	return c.infoLogger
}

func (c *container) GetDatabase() *gorm.DB {
	if c.db == nil {
		_, err := utils.CreateLogFile(c.DbPath, 0744)

		if err != nil {
			log.Fatal().Err(err).Msg("Error while creating database folder")
		}

		db, err := gorm.Open(sqlite.Open(c.DbPath), &gorm.Config{
			NowFunc: time.Now().UTC,
			Logger: gormlogger.New(
				stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags), // TODO: Write custom Database Logger
				gormlogger.Config{
					SlowThreshold:             time.Second,
					LogLevel:                  logger.Info,
					IgnoreRecordNotFoundError: true,
					Colorful:                  true,
				},
			),
			PrepareStmt:                              true,
			DisableForeignKeyConstraintWhenMigrating: false,
			DisableNestedTransaction:                 true,
		})

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error connecting to SQLITE3 Database")
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

func New(config Config) Container {
	return &container{
		Config: config,
	}
}
