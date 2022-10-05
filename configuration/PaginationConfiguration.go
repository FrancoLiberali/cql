package configuration

import "github.com/spf13/viper"

type PaginationConfiguration interface {
	GetMaxElemPerPage() uint
}

// Hold the configuration values for the pagination
type paginationConfigurationImpl struct{}

// Instantiate a new configuration holder for the pagination
func NewPaginationConfiguration() PaginationConfiguration {
	return &paginationConfigurationImpl{}
}

// Return the maximum number of element returned per page
func (lc *paginationConfigurationImpl) GetMaxElemPerPage() uint {
	return viper.GetUint("server.pagination.page.max")
}
