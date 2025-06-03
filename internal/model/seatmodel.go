package model

import (
	"context"
	"errors"
	"fmt"
	"imd-seat-be/internal/pkg/errorx"
	"imd-seat-be/internal/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SeatModel = (*customSeatModel)(nil)

type (
	// SeatModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSeatModel.
	SeatModel interface {
		GetAvaliabledays(ctx context.Context) ([]time.Time, error)
		GetSeatInfobyDateAndID(ctx context.Context, date time.Time, roomid string) ([]*Seat, error)
		ChangeSeatStatusByType(ctx context.Context, date time.Time, status, seat, typing string) error
		FindOneBySeatRoomDate(ctx context.Context, seat string, room string, date time.Time) (*Seat, error)
		InsertSeatsForDateRange(ctx context.Context, room string, seats []string, startDate, endDate string) error
		DeleteSeatsBeforeDate(ctx context.Context, date string) error
		seatModel
		withSession(session sqlx.Session) SeatModel
	}

	customSeatModel struct {
		*defaultSeatModel
	}
	daterow struct {
		Date time.Time `db:"date"`
	}
)

// 获取可预约日期
func (c *customSeatModel) GetAvaliabledays(ctx context.Context) ([]time.Time, error) {
	query := fmt.Sprintf("select distinct date as date from %s order by date", c.table)
	var rows []daterow
	err := c.conn.QueryRowsCtx(ctx, &rows, query)
	//整合结果
	dates := make([]time.Time, 0, len(rows))
	for _, r := range rows {
		dates = append(dates, r.Date)
	}
	return dates, err
}

// 获取某天某座位的具体信息
func (c *customSeatModel) GetSeatInfobyDateAndID(ctx context.Context, date time.Time, roomid string) ([]*Seat, error) {
	logx.Infof("查询座位信息: %v, 房间: %s", date, roomid)
	query := fmt.Sprintf("select %s from %s where `date` = ? and `room` = ? ", seatRows, c.table)
	var seats []*Seat
	err := c.conn.QueryRowsCtx(ctx, &seats, query, date, roomid)
	if err != nil {
		return nil, err
	}
	return seats, nil
}

// 改变座位状态
func (c *customSeatModel) ChangeSeatStatusByType(ctx context.Context, date time.Time, status, seat, typing string) error {
	if typing == "day" {
		query := fmt.Sprintf("update %s set `status` = ? where `seat` = ? and `date` = ?", c.table)
		_, err := c.conn.ExecCtx(ctx, query, status, seat, date)
		return err
	} else if typing == "week" {
		for i := 0; i < 7; i++ {
			currentDate := date.AddDate(0, 0, i)
			query := fmt.Sprintf("UPDATE %s SET `status` = ? WHERE `seat` = ? AND `date` = ?", c.table)
			if _, err := c.conn.ExecCtx(ctx, query, status, seat, currentDate); err != nil {
				return err
			}
		}
	} else {
		return errors.New("type 格式错误")
	}
	return nil
}

// 查看座位状态
func (c *customSeatModel) FindOneBySeatRoomDate(ctx context.Context, seat string, room string, date time.Time) (*Seat, error) {
	query := fmt.Sprintf("select %s from %s where `seat` = ? and `room` = ? and `date` = ? limit 1", seatRows, c.table)
	var seatInfo Seat
	err := c.conn.QueryRowCtx(ctx, &seatInfo, query, seat, room, date.Format(time.DateOnly))
	if err != nil {
		return nil, err
	}
	return &seatInfo, nil
}

// 指定日期之间生成
func (c *customSeatModel) InsertSeatsForDateRange(ctx context.Context, room string, seats []string, startDate, endDate string) error {
	sDate, err := time.Parse(time.DateOnly, startDate)
	if err != nil {
		return errors.New("startDate 参数不符合格式")
	}

	eDate, err := time.Parse(time.DateOnly, endDate)
	if err != nil {
		return errors.New("endDate 参数不符合格式")
	}

	for date := sDate; !date.After(eDate); date = date.AddDate(0, 0, 1) {
		for _, seatID := range seats {
			seat := &Seat{
				Seat:   seatID,
				Room:   room,
				Date:   date,
				Status: types.AvaliableStatus,
			}
			_, err := c.Insert(ctx, seat)
			if err != nil {
				return errorx.WrapError(errorx.CreateErr, fmt.Errorf("failed to insert seat %s on %v: %w", seatID, date, err))
			}

		}

	}
	return nil
}

func (c *customSeatModel) DeleteSeatsBeforeDate(ctx context.Context, date string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE `date` < ?", c.table)
	_, err := c.conn.ExecCtx(ctx, query, date)
	if err != nil {
		return errorx.WrapError(errorx.DeleteErr, fmt.Errorf("删除座位信息失败: %w", err))
	}
	return nil
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
