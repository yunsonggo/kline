package config

import (
	"github.com/yunsonggo/kline/pkg/xconfig"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	GinConf xconfig.GinConf
}
