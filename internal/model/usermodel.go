package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserModel = (*customUserModel)(nil)

type (
	// UserModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserModel.
	UserModel interface {
		userModel
		withSession(session sqlx.Session) UserModel
		FindScoreByID(ctx context.Context, studentID string) (int, error)
		UpdateScore(ctx context.Context, studentID string, score int) error
		RenewScore(ctx context.Context) error
	}

	customUserModel struct {
		*defaultUserModel
	}
)

// NewUserModel returns a model for the database table.
func NewUserModel(conn sqlx.SqlConn) UserModel {
	return &customUserModel{
		defaultUserModel: newUserModel(conn),
	}
}

// 获取信誉分
func (m *customUserModel) FindScoreByID(ctx context.Context, studentID string) (int, error) {
	query := fmt.Sprintf("select %s from %s where `student_id` = ?", userRows, m.table)
	var user User
	err := m.conn.QueryRowCtx(ctx, &user, query, studentID)
	if err != nil {
		return 0, err
	}
	return int(user.Score), err
}

// 更新信誉分
func (m *customUserModel) UpdateScore(ctx context.Context, studentID string, score int) error {
	query := fmt.Sprintf("update %s set `score` = ? where `student_id` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, score, studentID)
	return err
}

// 恢复信誉分
func (m *customUserModel) RenewScore(ctx context.Context) error {
	query := fmt.Sprintf("update %s set `score` = ?", m.table)
	_, err := m.conn.ExecCtx(ctx, query, 100)
	return err
}

func (m *customUserModel) withSession(session sqlx.Session) UserModel {
	return NewUserModel(sqlx.NewSqlConnFromSession(session))
}
