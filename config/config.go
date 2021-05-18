package config

import (
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	EncryptionKey string
	DB_USERNAME   string
	DB_PASSWORD   string
	DB_PORT       string
	DB_HOST       string
	DB_NAME       string
	PASS_POLICY   PassPolicy
}

type PassPolicy struct {
	PASS_SIZE         int
	PASS_UPPER        bool
	PASS_LETTER       bool
	PASS_NUMBER       bool
	PASS_HISTORY      int
	PASS_SPECIAL      bool
	IF_PASS_EXPIRE    bool
	LOCKOUT_DURATION  int
	DAYS_TOBE_EXPIRED int
	LOCKOUT_THRESHOLD int
}

func GetConfig() Configuration {
	configuration := Configuration{}
	gonfig.GetConf("config/config.json", &configuration)
	return configuration
}
