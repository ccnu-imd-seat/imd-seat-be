package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SeatModel = (*customSeatModel)(nil)

type (
	// SeatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeatModel.
	SeatModel interface {
		GetSeatInfobyDateAndID(ctx context.Context, date time.Time, roomid string) ([]*Seat, error)
		UpdateSeatsStatusAndDate(ctx context.Context, date time.Time) error
		seatModel
	}

	customSeatModel struct {
		*defaultSeatModel
	}
)

// 获取某天某座位的具体信息
func (c *customSeatModel) GetSeatInfobyDateAndID(ctx context.Context, date time.Time, roomid string) ([]*Seat, error) {
	query := fmt.Sprintf("select %s from %s where `date` = ? and `seat` = ? ", seatRows, c.table)
	var seats []*Seat
	err := c.QueryRowsNoCacheCtx(ctx, &seats, query, date, roomid)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

// 更新所有座位的日期和状态信息
func (c *customSeatModel) UpdateSeatsStatusAndDate(ctx context.Context, date time.Time) error {
	query := fmt.Sprintf("update %s set `status` = ? and `date` = ?", c.table)
	_, err := c.ExecNoCacheCtx(ctx, query, "available", date)
	return err
}

// NewSeatModel returns a model for the database table.
func NewSeatModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SeatModel {
	return &customSeatModel{
		defaultSeatModel: newSeatModel(conn, c, opts...),
	}
}
