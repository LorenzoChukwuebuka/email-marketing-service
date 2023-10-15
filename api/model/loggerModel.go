package model

import "time"

type LoggerModel struct {
	Id        int
	UUID string
	Action    string
	CreatedAt time.Time
}
