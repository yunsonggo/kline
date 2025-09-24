package logic

import (
	"context"

	"github.com/yunsonggo/kline/app/solana_app/internal/svc"
	"github.com/yunsonggo/kline/app/solana_app/proto/solanapb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PingLogic {
	return &PingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PingLogic) Ping(in *solanapb.PingReq) (*solanapb.PingResp, error) {
	// todo: add your logic here and delete this line

	return &solanapb.PingResp{}, nil
}
