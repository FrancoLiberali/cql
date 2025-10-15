package cql

import (
	"context"

	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/logger"
)

type DB struct {
	GormDB                *gorm.DB
	withLoggerFromContext *LoggerFromContext
}

// gormDBWithContext return a gormdb with changed context to ctx
func (db *DB) gormDBWithContext(ctx context.Context) *gorm.DB {
	withContext := db.GormDB.WithContext(ctx)

	if db.withLoggerFromContext != nil {
		newLogger := db.withLoggerFromContext.getLoggerFunc(ctx)
		if newLogger != nil {
			withContext.Logger = newLogger
		}
	}

	return withContext
}

type LoggerFromContext struct {
	getLoggerFunc func(context.Context) logger.Interface
}

func (*LoggerFromContext) Apply(*Config) error { return nil }

func (*LoggerFromContext) AfterInitialize(*gorm.DB) error { return nil }

// WithLoggerFromContext allows to set a way cql does all the logs from a logger that is in the context
func WithLoggerFromContext(getLoggerFunc func(context.Context) logger.Interface) *LoggerFromContext {
	return &LoggerFromContext{
		getLoggerFunc: getLoggerFunc,
	}
}

type (
	Dialector = gorm.Dialector
	Option    = gorm.Option
	Config    = gorm.Config
)

// Open initialize db session based on dialector
//
// For details see https://compiledquerylenguage.readthedocs.io/en/latest/cql/connecting_to_a_database.html
func Open(dialector Dialector, opts ...Option) (*DB, error) {
	configs := pie.Filter(opts, func(opt Option) bool {
		_, isConfig := opt.(*Config)
		return isConfig
	})

	var gormDB *gorm.DB

	var err error

	if len(configs) == 0 {
		gormDB, err = gorm.Open(dialector, append(opts, &Config{Logger: logger.Default})...)
		if err != nil {
			return nil, err
		}
	} else {
		lastConfig, _ := configs[len(configs)-1].(*Config)
		if lastConfig.Logger == nil {
			lastConfig.Logger = logger.Default
		}

		gormDB, err = gorm.Open(dialector, opts...)
		if err != nil {
			return nil, err
		}
	}

	db := &DB{
		GormDB: gormDB,
	}

	withLoggerFromContextIndex := pie.FindFirstUsing(opts, func(opt Option) bool {
		_, isWithLoggerFromContext := opt.(*LoggerFromContext)
		return isWithLoggerFromContext
	})
	if withLoggerFromContextIndex != -1 {
		db.withLoggerFromContext = opts[withLoggerFromContextIndex].(*LoggerFromContext)
	}

	return db, nil
}
