package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ SeatModel = (*customSeatModel)(nil)

type (
	// SeatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeatModel.
	SeatModel interface {
		seatModel
		withSession(session sqlx.Session) SeatModel
	}

	customSeatModel struct {
		*defaultSeatModel
	}
)

// NewSeatModel returns a model for the database table.
func NewSeatModel(conn sqlx.SqlConn) SeatModel {
	return &customSeatModel{
		defaultSeatModel: newSeatModel(conn),
	}
}

func (m *customSeatModel) withSession(session sqlx.Session) SeatModel {
	return NewSeatModel(sqlx.NewSqlConnFromSession(session))
}
