package configuration

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Hold the configuration values for the pagination
type PaginationConfiguration interface {
	ConfigurationHolder
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

func (paginationConfiguration *paginationConfigurationImpl) Log() {
	zap.L().Info("Pagination configuration",
		zap.Uint("maxelemPerPage", paginationConfiguration.pagesNb),
	)
}

// Reload pagination configuration
func (paginationConfiguration *paginationConfigurationImpl) Reload() {
	paginationConfiguration.pagesNb = viper.GetUint("server.pagination.page.max")
}
