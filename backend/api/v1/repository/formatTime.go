package repository

import "time"

func FormatTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	formattedTime := t.Format(time.RFC3339)
	return &formattedTime
}
