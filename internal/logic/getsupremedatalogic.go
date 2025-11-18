package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/pkg/response"
	"imd-seat-be/internal/svc"
	"imd-seat-be/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSupremeDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSupremeDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSupremeDataLogic {
	return &GetSupremeDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSupremeDataLogic) GetSupremeData() (resp *types.NormalResp, err error) {
	reservations, err := l.svcCtx.ReservationModel.GetAllReservations(l.ctx)
	if err != nil {
		return nil, errorx.WrapError(errorx.FetchErr, err)
	}

	var data []types.ReservationDetails

	// 结构体的转换
	for _, res := range reservations {
		data = append(data, types.ReservationDetails{
			ID:     int(res.Id),
			StuID:  res.StudentId,
			Type:   res.Type,
			Date:   res.Date.Format(time.DateOnly),
			Room:   res.Room,
			SeatID: res.Seat,
			Status: res.Status,
		})
	}

	fileName := fmt.Sprintf("export_%d.json", time.Now().UnixMilli())
	filePath := filepath.Join("files", fileName)

	err = os.MkdirAll("files", 0755)
	if err != nil {
		return nil, errorx.WrapError(errorx.SaveErr, err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return nil, errorx.WrapError(errorx.SaveErr, err)
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(data)
	if err != nil {
		return nil, errorx.WrapError(errorx.SaveErr, err)
	}

	downloadURL := fmt.Sprintf("%s/download?filename=%s", l.svcCtx.Config.Domain, fileName)

	resp = &types.NormalResp{
		Base: response.Success(),
		Data: downloadURL,
	}

	return resp, nil
}

func (l *GetSupremeDataLogic) GetSupremeList() (resp *types.SupremeList, err error) {
	resp = &types.SupremeList{
		Base: response.Success(),
		Data: struct {
			Admins []string `json:"admins"`
		}{
			l.svcCtx.Config.Admin.Id,
		},
	}

	return resp, nil
}

func (l *GetSupremeDataLogic) Download(filename string) (*os.File, error) {
	filePath := "./files/" + filename

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errorx.WrapError(errorx.NotFound, err)
	}

	return file, nil
}
