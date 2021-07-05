package di

import (
	"fmt"
	"gorm.io/driver/postgres"
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
		enabled := c.viper.GetBool("logging.db.enabled")

		var (
			logger gormlogger.Interface
			db     *gorm.DB
			err    error
		)

		slowQuery := c.viper.GetDuration("logging.db.slow_query")
		loggingLevel := c.viper.GetString("logging.db.level")
		colors := c.viper.GetBool("logging.db.color")

		log.Trace().
			Bool("enabled", enabled).
			Bool("colors", colors).
			Str("logging_level", loggingLevel).
			Str("slow_query", slowQuery.String()).
			Msg("Database logging")

		if enabled {
			logger = gormlogger.New(
				stdlog.New(os.Stdout, "\r\n", stdlog.LstdFlags),
				gormlogger.Config{
					SlowThreshold:             slowQuery,
					LogLevel:                  logging.ParseDBLogLevel(loggingLevel),
					IgnoreRecordNotFoundError: true,
					Colorful:                  colors,
				},
			)

			log.Trace().Msg("GORM Logger created")
		} else {
			logger = nil
		}

		driver := c.viper.GetString("database.driver")

		switch driver {
		case "sqlite":
			db, err = createSqliteDriver(c, logger)

			if err != nil {
				log.Fatal().
					Err(err).
					Msg("error connecting to SQLite3 Database")
			}
		case "postgres":
			db, err = createPostgresDriver(c, logger)

			if err != nil {
				log.Fatal().
					Err(err).
					Msg("error connecting to PostgreSQL database")
			}
		default:
			log.Fatal().Str("driver", driver).Msg("Not supported database driver")
		}

		maxIdleConnections,
		maxOpenConnections,
		connMaxLifetime,
		maxIdleLifeTime := getSqlConnectionSettings(c, "sqlite")

		log.Trace().
			Str("driver", driver).
			Int("max_idle_connections", maxIdleConnections).
			Int("max_open_connections", maxOpenConnections).
			Str("conn_max_lifetime", connMaxLifetime.String()).
			Str("max_idle_lifetime", maxIdleLifeTime.String()).
			Msg("Creating GORM Instance")

		sqlDb, err := db.DB()

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("error connecting to Database")
		}

		sqlDb.SetMaxIdleConns(maxIdleConnections)
		sqlDb.SetMaxIdleConns(maxOpenConnections)
		sqlDb.SetConnMaxLifetime(connMaxLifetime)
		sqlDb.SetConnMaxIdleTime(maxIdleLifeTime)

		log.Debug().Msg("Migrating the database")
		err = db.AutoMigrate(models.GetModels()...)

		if err != nil {
			log.Fatal().
				Err(err).
				Msg("Failed to migrate DATABASE")
		}

		log.Debug().Msg("Database migrated")

		c.db = db
	}

	return c.db
}

func createSqliteDriver(c *container, logger gormlogger.Interface) (*gorm.DB, error) {
	path := c.viper.GetString("database.sqlite.path")
	dbPath, err := utils.GetAbsolutePath(path)

	if err != nil {
		log.Error().Err(err).Str("path", path).Msg("Error while getting absolute path")
		return nil, err
	}

	log.Debug().Str("path", dbPath).Msg("SQLite3 Database path")

	_, err = utils.CreatePath(dbPath, 0744)

	if err != nil {
		log.Error().Err(err).Msg("Error while creating database folder")
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), getGormConfig(logger))

	if err != nil {
		log.Error().Err(err).Msg("Error while creating GORM Instance")
		return nil, err
	}

	return db, nil
}

func createPostgresDriver(c *container, logger gormlogger.Interface) (*gorm.DB, error) {
	dsn := c.viper.GetString("database.postgres.dsn")

	log.Debug().Str("dsn", dsn).Msg("Postgres DSN")

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: false,
		WithoutReturning:     false,
	}), getGormConfig(logger))

	if err != nil {
		log.Error().
			Err(err).
			Msg("Error while creating GORM Instance")

		return nil, err
	}

	return db, nil
}

func getGormConfig(logger gormlogger.Interface) *gorm.Config {
	return &gorm.Config{
		NowFunc:                                  time.Now().UTC,
		Logger:                                   logger,
		PrepareStmt:                              true,
		DisableForeignKeyConstraintWhenMigrating: false,
		DisableNestedTransaction:                 true,
	}
}

func getSqlConnectionSettings(c *container, driver string) (int, int, time.Duration, time.Duration) {
	maxIdleConnections := c.viper.GetInt(fmt.Sprintf("database.%s.max_idle_connections", driver))
	maxOpenConnections := c.viper.GetInt(fmt.Sprintf("database.%s.max_open_connections", driver))
	connMaxLifetime := c.viper.GetDuration(fmt.Sprintf("database.%s.conn_max_lifetime", driver))
	maxIdleLifeTime := c.viper.GetDuration(fmt.Sprintf("database.%s.max_idle_lifetime", driver))

	return maxIdleConnections, maxOpenConnections, connMaxLifetime, maxIdleLifeTime
}
