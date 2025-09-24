package svc

import (
	"github.com/yunsonggo/kline/api/solana_api/internal/config"
	"github.com/yunsonggo/kline/api/solana_api/internal/middleware"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config          config.Config
	CostInterceptor rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:          c,
		CostInterceptor: middleware.NewCostInterceptorMiddleware().Handle,
	}
}
