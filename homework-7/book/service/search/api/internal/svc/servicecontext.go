package svc

import (
	"book/service/search/api/internal/config"
	"book/service/search/api/internal/middleware"
	"book/service/search/model"
	"book/service/user/rpc/userclient"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	JWT         rest.Middleware
	UserRpc     userclient.User
	SearchModel model.SearchModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	coon := sqlx.NewMysql(c.Mysql.DataSource)
	return &ServiceContext{
		Config:      c,
		JWT:         middleware.NewJWTMiddleware().Handle,
		UserRpc:     userclient.NewUser(zrpc.MustNewClient(c.UserRpc)),
		SearchModel: model.NewSearchModel(coon, c.CacheRedis),
	}
}
