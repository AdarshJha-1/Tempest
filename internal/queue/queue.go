package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Queue interface {
	Ping() (string, error)
	PushJobID(jobID string) error
	GetJobID() (string, error)
}

type queue struct {
	rdb *redis.Client
	ctx context.Context
}

func (q *queue) Ping() (string, error) {
	return q.rdb.Ping(q.ctx).Result()
}

func (q *queue) PushJobID(jobID string) error {

	_, err := q.rdb.RPush(q.ctx, "jobs", jobID).Result()

	return err
}

// FIFO manner
func (q *queue) GetJobID() (string, error) {
	jobID, err := q.rdb.LPop(context.Background(), "jobs").Result()
	if err != nil {
		return "", err
	}
	return jobID, nil
}

func NewQueue() Queue {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	return &queue{
		rdb: rdb,
		ctx: ctx,
	}
}
