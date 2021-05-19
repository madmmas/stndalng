package config

import (
	"github.com/tkanos/gonfig"
)

type Configuration struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_PORT     string
	DB_HOST     string
	DB_NAME     string
}

type PassPolicy struct {
	PASS_SIZE          int
	PASS_UPPER         bool
	PASS_LETTER        bool
	PASS_NUMBER        bool
	PASS_HISTORY       int
	PASS_SPECIAL       bool
	IF_PASS_EXPIRE     bool
	LOCKOUT_DURATION   int
	DAYS_TOBE_EXPIRED  int
	LOCKOUT_THRESHOLD  int
	TOKEN_TOBE_EXPIRED int
	TOKEN_CRYPTO_KEY   string
}

func GetConfig() Configuration {
	configuration := Configuration{}
	gonfig.GetConf("config/config.json", &configuration)
	return configuration
}

func GetPassPolicy() PassPolicy {
	passPolicy := PassPolicy{}
	gonfig.GetConf("config/pass_policy.json", &passPolicy)
	return passPolicy
}
