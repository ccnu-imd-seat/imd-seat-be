package config

import (
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	MySQL struct {
		Name     string
		User     string
		Password string
		Host     string
		Port     int
	}

	Redis struct {
		Host     string
		Port     int
		Password string
	}
}
