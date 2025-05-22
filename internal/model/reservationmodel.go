package model

import "github.com/zeromicro/go-zero/core/stores/sqlx"

var _ ReservationModel = (*customReservationModel)(nil)

type (
	// ReservationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReservationModel.
	ReservationModel interface {
		reservationModel
		withSession(session sqlx.Session) ReservationModel
	}

	customReservationModel struct {
		*defaultReservationModel
	}
)

// NewReservationModel returns a model for the database table.
func NewReservationModel(conn sqlx.SqlConn) ReservationModel {
	return &customReservationModel{
		defaultReservationModel: newReservationModel(conn),
	}
}

func (m *customReservationModel) withSession(session sqlx.Session) ReservationModel {
	return NewReservationModel(sqlx.NewSqlConnFromSession(session))
}
