package model

import (
	"context"
	"fmt"
	"imd-seat-be/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ ReservationModel = (*customReservationModel)(nil)

type (
	// ReservationModel is an interface to be customized, add more methods here,
	// and implement the added methods in customReservationModel.
	ReservationModel interface {
		UpdateReservstionMessage(ctx context.Context, id int64, status string) error
		GetReservationByStatus(ctx context.Context, date time.Time, status string) ([]*Reservation, error)
		GetReservationByStudentId(ctx context.Context, studentId string) ([]*Reservation, error)
		GetTodayReservationByStudentId(ctx context.Context, studentId string, seat string) (*Reservation, error)
		GetReservationByDate(ctx context.Context, date time.Time) ([]*Reservation, error)
		GetAnydayReservationByStudentId(ctx context.Context, studentId string, seat string, date string) (*Reservation, error)
		GetAllReservations(ctx context.Context) ([]*Reservation, error)
		HasBookedReservationInSelectedDay(ctx context.Context, studentId string, date string) (bool, error)
		reservationModel
		withSession(session sqlx.Session) ReservationModel
	}

	customReservationModel struct {
		*defaultReservationModel
	}
)

func (c *customReservationModel) HasBookedReservationInSelectedDay(ctx context.Context, studentId string, date string) (bool, error) {
	query := fmt.Sprintf(
		"SELECT COUNT(1) FROM %s WHERE `student_id` = ? AND `status` = ? AND `date` = ?",
		c.table,
	)

	var count int
	err := c.conn.QueryRowCtx(ctx, &count, query,
		studentId,
		types.BookedStatus,
		date,
	)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// 更新预约状态
func (c *customReservationModel) UpdateReservstionMessage(ctx context.Context, id int64, status string) error {
	query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, status, id)
	return err
}

// 根据状态查找
func (c *customReservationModel) GetReservationByStatus(ctx context.Context, date time.Time, status string) ([]*Reservation, error) {
	query := fmt.Sprintf("select %s from %s where `status` = ? and `date` = ?", reservationRows, c.table)
	var reservations []*Reservation
	err := c.conn.QueryRowsCtx(ctx, &reservations, query, status, date)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// 根据学号查找
func (c *customReservationModel) GetReservationByStudentId(ctx context.Context, studentId string) ([]*Reservation, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `student_id` = ?", reservationRows, c.table)
	var reservations []*Reservation
	err := c.conn.QueryRowsCtx(ctx, &reservations, query, studentId)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// 根据根据学号查找当天预约进行签到处理
func (c *customReservationModel) GetTodayReservationByStudentId(ctx context.Context, studentId string, seat string) (*Reservation, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `student_id` = ? and `date` = ? and `seat` = ?", reservationRows, c.table)
	var reservation Reservation
	date := time.Now()
	dateStr := date.Format(time.DateOnly)
	err := c.conn.QueryRowCtx(ctx, &reservation, query, studentId, dateStr, seat)
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

// 根据根据学号和日期查找当天预约进行签到处理
func (c *customReservationModel) GetAnydayReservationByStudentId(ctx context.Context, studentId string, seat string, date string) (*Reservation, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `student_id` = ? and `date` = ? and `seat` = ?", reservationRows, c.table)
	var reservation Reservation
	err := c.conn.QueryRowCtx(ctx, &reservation, query, studentId, date, seat)
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

// 根据日期查找
func (c *customReservationModel) GetReservationByDate(ctx context.Context, date time.Time) ([]*Reservation, error) {
	query := fmt.Sprintf("SELECT %s FROM %s WHERE `date` = ?", reservationRows, c.table)
	var reservations []*Reservation
	err := c.conn.QueryRowsCtx(ctx, &reservations, query, date)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// GetAllReservations 查找所有预约记录
func (c *customReservationModel) GetAllReservations(ctx context.Context) ([]*Reservation, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", reservationRows, c.table)
	var reservations []*Reservation
	err := c.conn.QueryRowsCtx(ctx, &reservations, query)
	if err != nil {
		return nil, err
	}
	return reservations, nil
}

// NewReservationModel returns a model for the database table.
func NewReservationModel(conn sqlx.SqlConn) ReservationModel {
	return &customReservationModel{
		defaultReservationModel: newReservationModel(conn),
	}
}

func (m *customReservationModel) withSession(session sqlx.Session) ReservationModel {
	return NewReservationModel(sqlx.NewSqlConnFromSession(session))
}
