package domain

import "time"

type Event struct {
	ID        string
	Title     string
	StartTime time.Time
	EndTime   time.Time
	UpdatedAt time.Time
}
