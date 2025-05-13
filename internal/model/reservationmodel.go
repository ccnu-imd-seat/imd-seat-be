package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReservationModel = (*customReservationModel)(nil)

type (
	// ReservationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReservationModel.
	ReservationModel interface {
		reservationModel
	}

	customReservationModel struct {
		*defaultReservationModel
	}
)

// NewReservationModel returns a model for the database table.
func NewReservationModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) ReservationModel {
	return &customReservationModel{
		defaultReservationModel: newReservationModel(conn, c, opts...),
	}
}
