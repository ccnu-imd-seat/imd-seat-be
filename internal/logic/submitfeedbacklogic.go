package logic

import (
	"context"
	"errors"

	"imd-seat-be/internal/model"
	"imd-seat-be/internal/pkg/contextx"
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SubmitFeedbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSubmitFeedbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitFeedbackLogic {
	return &SubmitFeedbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SubmitFeedbackLogic) SubmitFeedback(req *types.FeedbackReq) (resp *types.GeneralRes, err error) {
	studentID, ok := contextx.GetStudentID(l.ctx)
	if !ok {
		return nil, errorx.WrapError(errorx.JWTError, errors.New("token读取学号失败"))
	}

	Feedback := model.Feedback{
		StudentId: studentID,
		Content:   req.Content,
	}

	_, err = l.svcCtx.FeedbackModel.Insert(l.ctx, &Feedback)
	if err != nil {
		return nil, errorx.WrapError(errorx.CreateErr, err)
	}

	resp = response.GeneralRes()

	return resp, nil
}
