package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SeatModel = (*customSeatModel)(nil)

type (
	// SeatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeatModel.
	SeatModel interface {
		GetSeatInfobyDateAndID(ctx context.Context, date time.Time, roomid string) ([]*Seat, error)
		ChangeSeatStatus(ctx context.Context, date time.Time, status, seat string) error
		seatModel
		withSession(session sqlx.Session) SeatModel
	}

	customSeatModel struct {
		*defaultSeatModel
	}
)

// 获取某天某座位的具体信息
func (c *customSeatModel) GetSeatInfobyDateAndID(ctx context.Context, date time.Time, roomid string) ([]*Seat, error) {
	query := fmt.Sprintf("select %s from %s where `date` = ? and `room` = ? ", seatRows, c.table)
	var seats []*Seat
	err := c.conn.QueryRowsCtx(ctx, &seats, query, date, roomid)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

// 改变座位状态
func (c *customSeatModel) ChangeSeatStatus(ctx context.Context, date time.Time, status, seat string) error {
	query := fmt.Sprintf("update %s set `status` = ? where `seat` = ? and `date` = ?", c.table)
	_,err := c.conn.ExecCtx(ctx, query, status, seat, date)
	return err
}

// NewSeatModel returns a model for the database table.
func NewSeatModel(conn sqlx.SqlConn) SeatModel {
	return &customSeatModel{
		defaultSeatModel: newSeatModel(conn),
	}
}

func (m *customSeatModel) withSession(session sqlx.Session) SeatModel {
	return NewSeatModel(sqlx.NewSqlConnFromSession(session))
}
