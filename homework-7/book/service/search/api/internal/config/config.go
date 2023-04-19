package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret  string
		RefreshSecret string
	}
	UserRpc zrpc.RpcClientConf
	Mysql   struct {
		DataSource string
	}
	CacheRedis cache.CacheConf
}
