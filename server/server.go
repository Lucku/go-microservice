package server

import (
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

// Start initializes the server, reading in its configuration, and starts it
func Start() error {

	if err := LoadConfig(); err != nil {
		return err
	}

	router := Routes()

	// DEBUG: DELETE AT SOME POINT
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error())
	}
	// DEBUG END

	port := viper.GetString("port")

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

	return nil
}
