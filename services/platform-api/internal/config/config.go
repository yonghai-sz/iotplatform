// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Transform   zrpc.RpcClientConf
	PublicFiles string
	DataSource  string
	Auth        struct {
		AccessSecret string
		AccessExpire int64
	}
	Redis redis.RedisConf
}
