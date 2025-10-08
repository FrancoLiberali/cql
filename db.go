package cql

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/FrancoLiberali/cql/logger"
)

type DB struct {
	GormDB *gorm.DB
}

// Open initialize db session based on dialector
//
// For details see https://compiledquerylenguage.readthedocs.io/en/latest/cql/connecting_to_a_database.html
func Open(dialector gorm.Dialector, opts ...gorm.Option) (*DB, error) {
	configs := pie.Filter(opts, func(opt gorm.Option) bool {
		_, isConfig := opt.(*gorm.Config)
		return isConfig
	})

	if len(configs) == 0 {
		gormDB, err := gorm.Open(dialector, append(opts, &gorm.Config{Logger: logger.Default})...)
		if err != nil {
			return nil, err
		}

		return &DB{
			GormDB: gormDB,
		}, nil
	}

	lastConfig, _ := configs[len(configs)-1].(*gorm.Config)
	if lastConfig.Logger == nil {
		lastConfig.Logger = logger.Default
	}

	gormDB, err := gorm.Open(dialector, opts...)
	if err != nil {
		return nil, err
	}

	return &DB{
		GormDB: gormDB,
	}, nil
}
