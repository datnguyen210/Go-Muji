package mocks

import (
	"time"

	"github.com/datnguyen210/go-blog/internal/models"
)

var mockBlog = &models.Blog{
	ID: 1,
	Title: "An old friend",
	Content: "Old but gold",
	Created: time.Now(),
	Expires: time.Now(),
}

type BlogModel struct{}

func (m *BlogModel) Insert(title, content string, expires int) (int, error){
	return 2, nil
}

func (m *BlogModel) Get (id int) (*models.Blog, error) {
	switch id {
	case 1:
		return mockBlog, nil
	default: 
		return nil, models.ErrNoRecord
	}
}

func (m *BlogModel) Latest() ([]*models.Blog, error) {
	return []*models.Blog{mockBlog}, nil
}