package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Hold the configuration values for the pagination
type PaginationConfiguration interface {
	Holder
	GetMaxElemPerPage() uint
}

// Concrete implementation of the PaginationConfiguration interface
type paginationConfigurationImpl struct {
	pagesNb uint
}

// Instantiate a new configuration holder for the pagination
func NewPaginationConfiguration() PaginationConfiguration {
	paginationConfiguration := new(paginationConfigurationImpl)
	paginationConfiguration.Reload()

	return paginationConfiguration
}

// Return the maximum number of element returned per page
func (paginationConfiguration *paginationConfigurationImpl) GetMaxElemPerPage() uint {
	return paginationConfiguration.pagesNb
}

// Reload pagination configuration
func (paginationConfiguration *paginationConfigurationImpl) Reload() {
	paginationConfiguration.pagesNb = viper.GetUint(ServerPaginationMaxElemPerPage)
}

// Log the values provided by the configuration holder
func (paginationConfiguration *paginationConfigurationImpl) Log(logger *zap.Logger) {
	logger.Info("Pagination configuration",
		zap.Uint("maxelemPerPage", paginationConfiguration.pagesNb),
	)
}
