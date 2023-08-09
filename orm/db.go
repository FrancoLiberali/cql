package orm

import (
	"github.com/elliotchance/pie/v2"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/orm/logger"
)

// Open initialize db session based on dialector
func Open(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
	configs := pie.Filter(opts, func(opt gorm.Option) bool {
		_, isConfig := opt.(*gorm.Config)
		return isConfig
	})

	if len(configs) == 0 {
		return gorm.Open(dialector, append(opts, &gorm.Config{Logger: logger.Default})...)
	}

	lastConfig, _ := configs[len(configs)-1].(*gorm.Config)
	if lastConfig.Logger == nil {
		lastConfig.Logger = logger.Default
	}

	return gorm.Open(dialector, opts...)
}
