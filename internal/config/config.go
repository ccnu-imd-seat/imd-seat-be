package config

import (
	"fmt"

	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	Domain string

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

	Qiniu struct {
		AccessKey string
		SecretKey string
		Domain    string
		Bucket    string
	}

	Admin struct {
		Id []string
	}
}

func (c Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		c.MySQL.User,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.Name,
	)
}
