package models

import "time"

type TimeEntry struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time,omitempty"`
	Duration  int64     `json:"duration,omitempty"`
	TaskName  string    `json:"task_name"`
	Tags      []string  `json:"tags"`
}
