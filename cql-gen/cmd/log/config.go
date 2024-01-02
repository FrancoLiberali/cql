package log

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/FrancoLiberali/cql-gen/cmd/version"
)

var Logger = log.WithField("version", version.Version)

const VerboseKey = "verbose"

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
}

func SetLevel() {
	verbose := viper.GetBool(VerboseKey)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}
