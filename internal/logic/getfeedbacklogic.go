package logic

import (
	"context"

	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFeedbackLogic {
	return &GetFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetFeedbackLogic) GetFeedback() (resp *types.FeedbackRes, err error) {
	reservations, err := l.svcCtx.FeedbackModel.FindAll(l.ctx)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}

	var data []types.Feedback

	// 结构体的转换
	for _, res := range reservations {
		data = append(data, types.Feedback{
			StudentId: res.StudentId,
			Content:   res.Content,
		})
	}

	resp = &types.FeedbackRes{
		Base: response.Success(),
		Data: data,
	}

	return resp, nil
}
