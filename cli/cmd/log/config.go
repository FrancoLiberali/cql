package log

import (
	"github.com/ditrit/badaas-orm/cli/cmd/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
