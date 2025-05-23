package logic

import (
	"context"

	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetViolationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetViolationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetViolationLogic {
	return &GetViolationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetViolationLogic) GetViolation() (resp *types.GeneralRes, err error) {
	

	return
}
