package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"strings"
	"time"
)

type SQLiteQuestions struct {
	db *sql.DB
}

type QuestionsRepo interface {
	Count(ctx context.Context) (int, error)
	GetRandom(ctx context.Context, limit int, category *string, difficulty *string) ([]Question, error)

	Insert(ctx context.Context, q Question) (int64, error)
	GetByID(ctx context.Context, id int64) (Question, error)
	UpdateByID(ctx context.Context, id int64, q Question) error
	DeleteByID(ctx context.Context, id int64) error
}

func NewQuestionsRepo(db *sql.DB) QuestionsRepo {
	return &SQLiteQuestions{db: db}
}

func (r *SQLiteQuestions) Count(ctx context.Context) (int, error) {
	var n int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM questions`).Scan(&n)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func (r *SQLiteQuestions) GetRandom(ctx context.Context, limit int, category *string, difficulty *string) ([]Question, error) {
	if limit <= 0 {
		limit = 1
	}
	if limit > 100 {
		limit = 100
	}

	q := "SELECT id, prompt, options_json, correct_index, category, difficulty, created_at FROM questions"
	var where []string
	var args []any

	if category != nil && strings.TrimSpace(*category) != "" {
		where = append(where, "category = ?")
		args = append(args, strings.TrimSpace(*category))
	}

	if difficulty != nil && strings.TrimSpace(*difficulty) != "" {
		where = append(where, "difficulty = ?")
		args = append(args, strings.TrimSpace(*difficulty))
	}

	if len(where) > 0 {
		q += " WHERE " + strings.Join(where, " AND ")
	}

	q += " ORDER BY RANDOM() LIMIT ?"

	args = append(args, limit)

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		log.Printf("REPO -> GETRND | error running query")
		return nil, err
	}
	defer rows.Close()

	out := make([]Question, 0, limit)
	for rows.Next() {
		var id int64
		var prompt, optJSON string
		var correct int
		var cat, diff sql.NullString
		var created time.Time

		if err := rows.Scan(&id, &prompt, &optJSON, &correct, &cat, &diff, &created); err != nil {
			return nil, err
		}

		var opts []string
		if err := json.Unmarshal([]byte(optJSON), &opts); err != nil {
			return nil, err
		}

		item := Question{
			ID:           id,
			Prompt:       prompt,
			Options:      opts,
			CorrectIndex: correct,
			Category:     nsp(cat),
			Difficulty:   nsp(diff),
			CreatedAt:    created,
		}

		out = append(out, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *SQLiteQuestions) Insert(ctx context.Context, q Question) (int64, error) {
	if err := q.Validate(); err != nil {
		log.Printf("REPO -> INSERT | Invalid Question")
		return 0, err
	}

	optsJSON, err := json.Marshal(q.Options)
	if err != nil {
		log.Printf("REPO -> INSERT | Error marshaling options")
		return 0, err
	}

	query := `
		INSERT INTO questions (prompt, options_json, correct_index, category, difficulty)
		VALUES (?, ?, ?, ?, ?)`

	res, err := r.db.ExecContext(ctx, query, q.Prompt, string(optsJSON), q.CorrectIndex, valNil(q.Category), valNil(q.Difficulty))
	if err != nil {
		log.Printf("REPO -> INSERT | Error executing query")
		return 0, err
	}

	id, err := res.LastInsertId()

	return id, err
}

func (r *SQLiteQuestions) GetByID(ctx context.Context, id int64) (Question, error) {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	var rowID int64
	var prompt, optJSON string
	var correct int
	var cat, diff sql.NullString
	var created time.Time

	query := `
		SELECT id, prompt, options_json, correct_index, category, difficulty, created_at
		FROM questions WHERE id = ?`

	err := r.db.QueryRowContext(ctx, query, id).Scan(&rowID, &prompt, &optJSON, &correct, &cat, &diff, &created)
	if err != nil {
		log.Printf("REPO -> GETBYID | Error executing query")
		return Question{}, err
	}

	var opts []string
	if err := json.Unmarshal([]byte(optJSON), &opts); err != nil {
		log.Printf("REPO -> GETBYID | Error unmarshalling options")
		return Question{}, err
	}

	ques := Question{
		ID:           rowID,
		Prompt:       prompt,
		Options:      opts,
		CorrectIndex: correct,
		Category:     nsp(cat),
		Difficulty:   nsp(diff),
		CreatedAt:    created,
	}
	return ques, nil
}

func (r *SQLiteQuestions) UpdateByID(ctx context.Context, id int64, q Question) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	err := q.Validate()
	if err != nil {
		log.Printf("REPO -> UPDATEBYID | Invalid Question object")
		return err
	}

	optsJSON, err := json.Marshal(q.Options)
	if err != nil {
		log.Printf("REPO -> UPDATEBYID | Error marshaling options")
		return err
	}

	query := `
		UPDATE questions
		SET prompt = ?, options_json = ?, correct_index = ?, category = ?, difficulty = ?
		WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, q.Prompt, string(optsJSON), q.CorrectIndex, valNil(q.Category), valNil(q.Difficulty), id)
	if err != nil {
		log.Printf("REPO -> UPDATEBYID | Error executing query")
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("REPO -> UPDATEBYID | Error rows affected")
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *SQLiteQuestions) DeleteByID(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	query := `
		DELETE FROM questions WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		log.Printf("REPO -> DELETEBYID | Error executing query")
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		log.Printf("REPO -> DELETEBYID | Error rows affected")
		return err
	}

	if affected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func nsp(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func valNil(s *string) any {
	if s == nil {
		return nil
	}
	return *s
}
