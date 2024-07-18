package model

import "time"

type Logger struct {
	Id        int
	UUID string
	Action    string
	CreatedAt time.Time
}
