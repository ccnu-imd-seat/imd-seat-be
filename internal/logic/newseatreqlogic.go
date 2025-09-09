package logic

import (
	"context"
	"fmt"
	"net/url"

	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/skip2/go-qrcode"
	"github.com/zeromicro/go-zero/core/logx"
)

type NewSeatReqLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewNewSeatReqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NewSeatReqLogic {
	return &NewSeatReqLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *NewSeatReqLogic) NewSeatReq(req *types.NewSeatReq) (resp *types.NewSeatRes, err error) {
	var URLs []string
	for _, seatID := range req.Seat_id {
		qrContent := fmt.Sprintf("http://%s:%d/api/v1/checkin?seat_id=%s",
			l.svcCtx.Config.Host,
			l.svcCtx.Config.Port,
			url.QueryEscape(seatID),
		)
		filename := seatID + ".png"
		png, err := qrcode.Encode(qrContent, qrcode.High, 256)
		if err != nil {
			return nil, errorx.WrapError(errorx.DefaultErr, err)
		}

		//上传到图床
		url, err := l.svcCtx.ImageUploader.Upload(l.ctx, png, "seatQRs/"+filename, filename)
		if err != nil {
			return nil, errorx.WrapError(errorx.DefaultErr, err)
		}
		URLs = append(URLs, url)
	}
	resp = &types.NewSeatRes{
		Base: response.Success(),
		Data: types.ImageData{
			ImageURL: URLs,
		},
	}
	return
}
