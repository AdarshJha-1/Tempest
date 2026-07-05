package queue

import "time"

type Job struct {
	ID         string
	Name       string
	Status     string
	ConfigJSON []byte
	CreatedAt  time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
}
