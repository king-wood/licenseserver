package main

import (
	"licenseserver/models"
	_ "licenseserver/routers"

	"net/http"

	log "github.com/cihub/seelog"
	"github.com/spf13/viper"
)

const (
	SERVER_CONFIG_FILE = "conf/app.toml"
	LOG_CONFIG_FILE    = "conf/log.xml"
)

func init() {
	viper.Reset()
	viper.SetConfigType("toml")
	viper.SetConfigFile(SERVER_CONFIG_FILE)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	models.InitDB(viper.GetString("database.dbType"), viper.GetString("database.dbConnStr"))
	logger, err := log.LoggerFromConfigAsFile(LOG_CONFIG_FILE)
	if err != nil {
		panic(err)
	}
	log.ReplaceLogger(logger)
}

func main() {
	log.Debug(http.ListenAndServe(":"+viper.GetString("server.port"), nil))
}
