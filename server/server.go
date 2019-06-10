package server

import (
	"net/http"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/viper"
)

// Start initializes the server, reading in its configuration, and starts it
func Start() error {

	if err := LoadConfig(); err != nil {
		return err
	}

	InitLogging()

	router := Routes()

	port := viper.GetString("port")

	log.Infof("Server started on port %s", port)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		return err
	}

	return nil
}

// LoadConfig loads the configuration of the server from a config file inside the directory of execution
// The config file contains entries to set the access key for the consumed API and a port definition. The
// former is a mandatory parameter, while the latter, if absent, defaults to port 8080
func LoadConfig() error {

	viper.SetConfigName("config")
	viper.AddConfigPath("./")

	viper.SetEnvPrefix("API")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	viper.SetDefault("port", 8080)
	viper.SetDefault("logging-level", "info")

	return nil
}

func InitLogging() {

	logLevel := viper.GetString("logging-level")

	log.SetLevelFromString(logLevel)
}
