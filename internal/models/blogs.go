package models

import (
	"database/sql"
	"errors"
	"time"
)

type BlogModelInterface interface {
	Insert(title, content string, expires int) (int, error)
	Get(id int) (*Blog, error)
	Latest() ([]*Blog, error)
}

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

	b := &Blog{}

	err := m.DB.QueryRow(stmt, id).Scan(&b.ID, &b.Title, &b.Content, &b.Created, &b.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return b, nil
}

func (m *BlogModel) Latest() ([]*Blog, error) {
	stmt := `SELECT id, title, content, created, expires FROM blogs
	WHERE expires > UTC_TIMESTAMP() ORDER BY id DESC LIMIT 10`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	// Closing a resultset with defer rows.Close() is
	// critical in the code above. As long as a resultset is open it will
	// keep the underlying database connection open… so if
	// something goes wrong in this method and the resultset isn’t
	// closed, it can rapidly lead to all the connections in your pool
	// being used up.
	defer rows.Close()

	blogs := []*Blog{}
	for(rows.Next()) {
		b := &Blog{}

		err := rows.Scan(&b.ID, &b.Title, &b.Content, &b.Created, &b.Expires)
		if err != nil {
			return nil, err
		}
		blogs = append(blogs, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return blogs, nil
}

