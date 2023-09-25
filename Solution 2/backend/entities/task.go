package entities

import "time"

type Task struct {
	Id        uint
	Name      string
	Deadline  time.Time
	Complete  bool
	Progress  float64
	Expired   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
