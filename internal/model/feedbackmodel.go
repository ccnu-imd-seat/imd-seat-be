package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FeedbackModel = (*customFeedbackModel)(nil)

type (
	// FeedbackModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFeedbackModel.
	FeedbackModel interface {
		feedbackModel
		withSession(session sqlx.Session) FeedbackModel
		FindAll(ctx context.Context) ([]*Feedback, error)
	}

	customFeedbackModel struct {
		*defaultFeedbackModel
	}
)

// NewFeedbackModel returns a model for the database table.
func NewFeedbackModel(conn sqlx.SqlConn) FeedbackModel {
	return &customFeedbackModel{
		defaultFeedbackModel: newFeedbackModel(conn),
	}
}

func (m *customFeedbackModel) FindAll(ctx context.Context) ([]*Feedback, error) {
	query := fmt.Sprintf("SELECT %s FROM %s", feedbackRows, m.table)
	var feedback []*Feedback
	err := m.conn.QueryRowsCtx(ctx, &feedback, query)
	if err != nil {
		return nil, err
	}
	return feedback, nil
}

func (m *customFeedbackModel) withSession(session sqlx.Session) FeedbackModel {
	return NewFeedbackModel(sqlx.NewSqlConnFromSession(session))
}
