package models

import (
	"database/sql"
	"errors"
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

func (m *BlogModel) Insert(title string, content string, expires int) (int, error) {
	stmt := `INSERT INTO blogs (title, content, created, expires) 
	VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *BlogModel) Get(id int) (*Blog, error) {
	stmt := `SELECT id, title, content, created, expires FROM blogs 
	WHERE expires > UTC_TIMESTAMP() AND id =?`

	s := &Blog{}

	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (m *BlogModel) Latest() ([]*Blog, error) {
	stmt := `SELECT id, title, content, created, expires FROM blogs
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	blogs := []*Blog{}

	for rows.Next() {
		blog := &Blog{}
		err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.Created, &blog.Expires)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, blog)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return blogs, nil
}
