package logic

import (
	"context"
	"errors"

	"imd-seat-be/internal/pkg/contextx"
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetScoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetScoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetScoreLogic {
	return &GetScoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetScoreLogic) GetScore() (resp *types.MyScoreRes, err error) {
	studentID, ok := contextx.GetStudentID(l.ctx)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}

	score, err := l.svcCtx.UserModel.FindScoreByID(l.ctx, studentID)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}

	resp = &types.MyScoreRes{
		Base: response.Success(),
		Data: types.ScoreData{
			Score: score,
		},
	}

	return resp, nil
}
