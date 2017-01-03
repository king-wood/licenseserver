package main

import (
	"licenseserver/models"
	_ "licenseserver/routers"

	"github.com/astaxie/beego"
	log "github.com/cihub/seelog"
)

func init() {
	models.InitDB(beego.AppConfig.String("dbType"), beego.AppConfig.String("dbConnStr"))

	logger, err := log.LoggerFromConfigAsFile("conf/log.xml")
	if err != nil {
		panic(err)
	}
	log.ReplaceLogger(logger)
}

func main() {
	beego.Run()
}
