package store

import (
	"context"
	"database/sql"
)

type SQLiteQuestions struct {
	db *sql.DB
}

type QuestionsRepo interface {
	Count(ctx context.Context) (int, error)
	GetRandom(ctx context.Context, limit int, cat *string, diff *string) ([]Question, error)

	Insert(ctx context.Context, q Question) (int64, error)
	GetByID(ctx context.Context, id int64) (Question, error)
	UpdateByID(ctx context.Context, id int64, q Question) error
	DeleteByID(ctx context.Context, id int64) error
}

func NewQuestionsRepo(db *sql.DB) QuestionsRepo {
	return &SQLiteQuestions{db: db}
}

func (r *SQLiteQuestions) Count(ctx context.Context) (int, error) {
	return 1, nil
}

func (r *SQLiteQuestions) GetRandom(ctx context.Context, limit int, cat *string, diff *string) ([]Question, error) {
	return []Question{}, nil
}

func (r *SQLiteQuestions) Insert(ctx context.Context, q Question) (int64, error) {
	return 1, nil
}

func (r *SQLiteQuestions) GetByID(ctx context.Context, id int64) (Question, error) {
	return Question{}, nil
}

func (r *SQLiteQuestions) UpdateByID(Ctx context.Context, id int64, q Question) error {
	return nil
}

func (r *SQLiteQuestions) DeleteByID(ctx context.Context, id int64) error {
	return nil
}
