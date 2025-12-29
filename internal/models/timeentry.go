package models

import "time"

type TimeEntry struct {
	ID        string    `json:"id"`
	StartTime time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	TaskName  string    `json:"task_name"`
	Tags      []string  `json:"tags"`
}

func (t *TimeEntry) Duration() time.Duration {
	if t.EndTime == nil {
		return time.Since(t.StartTime)
	}
	return t.EndTime.Sub(t.StartTime)
}

func (t *TimeEntry) IsRunning() bool {
	return t.EndTime == nil
}
