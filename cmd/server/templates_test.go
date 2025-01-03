package main

import (
	"testing"
	"time"

	"github.com/datnguyen210/go-blog/internal/assert"
)

func TestFormatDate(t *testing.T) {
	tests := []struct {
		name string 
		tm time.Time
		want string
	} {
		{
			name: "UTC",
			tm : time.Date(2021, 12, 31, 12, 30, 0, 0, time.UTC),
			want: "2021-12-31 at 12:30",
		},
		{
			name: "Empty",
			tm: time.Time{},
			want: "",
		},
		{
			name: "HKT",
			tm: time.Date(2044, 1, 1, 1, 0, 0, 0, time.FixedZone("HKT", 8*60*60)),
			want: "2043-12-31 at 17:00",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hd := formatDate(test.tm)

			assert.Equal(t, hd, test.want)
		})
	}
}