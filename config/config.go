package config

import (
	"stndalng/utils"

	"github.com/tkanos/gonfig"
)

type PassPolicy struct {
	Fld1 string
	Fld2 string
}

type Configuration struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
	PASS_POLICY utils.JSON
}

func GetConfig() Configuration {
	configuration := Configuration{}
	gonfig.GetConf("config/config.json", &configuration)
	return configuration
}
