package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReservationModel = (*customReservationModel)(nil)

type (
	// ReservationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReservationModel.
	ReservationModel interface {
		UpdateReservstionMessage(ctx context.Context, id int64, status string) error
		GetReservationByStatus(ctx context.Context, date time.Time,status string) ([]*Reservation, error)
		reservationModel
	}

	customReservationModel struct {
		*defaultReservationModel
	}
)

// 更新预约状态
func (c *customReservationModel) UpdateReservstionMessage(ctx context.Context, id int64, status string) error {
	query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", c.table)
	_, err := c.ExecNoCacheCtx(ctx, query, status, id)
	return err
}

// 根据状态查找
func (c *customReservationModel) GetReservationByStatus(ctx context.Context, date time.Time, status string) ([]*Reservation, error) {
	query := fmt.Sprintf("select %s from %s where `status` = ? and `date` = ?", seatRows, c.table)
	var reservations []*Reservation
	err := c.QueryRowsNoCacheCtx(ctx, &reservations, query, status, date)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// NewReservationModel returns a model for the database table.
func NewReservationModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ReservationModel {
	return &customReservationModel{
		defaultReservationModel: newReservationModel(conn, c, opts...),
	}
}
