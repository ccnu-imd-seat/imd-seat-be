package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SeatModel = (*customSeatModel)(nil)

type (
	// SeatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeatModel.
	SeatModel interface {
		seatModel
	}

	customSeatModel struct {
		*defaultSeatModel
	}
)

// NewSeatModel returns a model for the database table.
func NewSeatModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SeatModel {
	return &customSeatModel{
		defaultSeatModel: newSeatModel(conn, c, opts...),
	}
}
