package config

import "github.com/zeromicro/go-zero/core/service"

type Config struct {
	service.ServiceConf
	ListenOn string `json:",default=0.0.0.0:7611"`
}
