package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoomModel = (*customRoomModel)(nil)

type (
	// RoomModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoomModel.
	RoomModel interface {
		GetAvailableRoom(ctx context.Context, status string) ([]*Room, error)
		roomModel
	}

	customRoomModel struct {
		*defaultRoomModel
	}
)

// 获取可预约的房间
func (c *customRoomModel) GetAvailableRoom(ctx context.Context, status string) ([]*Room, error) {
	query := fmt.Sprintf("select %s from %s where `status` = ?  ", roomRows, c.table)
	var rooms []*Room
	err := c.conn.QueryRowsCtx(ctx,&rooms, query, status)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

// NewRoomModel returns a model for the database table.
func NewRoomModel(conn sqlx.SqlConn, opts ...cache.Option) RoomModel {
	return &customRoomModel{
		defaultRoomModel: newRoomModel(conn),
	}
}
