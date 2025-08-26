package store

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Question struct {
	ID           int64
	Prompt       string
	Options      []string
	CorrectIndex int
	Category     *string
	Difficulty   *string
	CreatedAt    time.Time
}

func (q *Question) Validate() error {

	p := strings.TrimSpace(q.Prompt)
	if p == "" {
		return errors.New("MODEL -> VALIDATE | Prompt required")
	}

	if len(q.Options) < 2 || len(q.Options) > 8 {
		return errors.New("MODEL -> VALIDATE | options must be between 2..8")
	}

	for i, opt := range q.Options {
		if strings.TrimSpace(opt) == "" {
			return errors.New("MODEL -> VALDIATE | option " + strconv.Itoa(i) + " empty")
		}
	}

	if q.CorrectIndex < 0 || q.CorrectIndex >= len(q.Options) {
		return errors.New("MODEL -> VALIDATE | answer out of range")
	}

	return nil
}
