package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf

	Swagger struct {
		Host string `json:"Host"`
	}
}
