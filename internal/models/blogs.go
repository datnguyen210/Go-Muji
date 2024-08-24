package models

import (
	"database/sql"
	"time"
)

type Blog struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type BlogModel struct {
	DB *sql.DB
}

func (m *BlogModel) Insert(title string, content string, expired *time.Time) (int, error) {
	return 0, nil
}

func (m *BlogModel) Get(id int) (*Blog, error) {
	return nil, nil
}

func (m *BlogModel) Latest() ([]*Blog, error) {
	return nil, nil
}
