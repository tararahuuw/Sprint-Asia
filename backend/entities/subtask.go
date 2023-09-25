package entities

import "time"

type SubTask struct {
	Id        uint   `json:"Id"`
	Id_Task   uint   `json:"Id_Task"`
	Name      string `json:"Name"`
	Complete  bool   `json:"Complete"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
