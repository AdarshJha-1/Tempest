package queue

import (
	"context"
	"time"

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
	ctx, cancel := context.WithTimeout(q.ctx, 2*time.Second)
	defer cancel()
	return q.rdb.Ping(ctx).Result()
}

func (q *queue) PushJobID(jobID string) error {
	ctx, cancel := context.WithTimeout(q.ctx, 2*time.Second)
	defer cancel()
	_, err := q.rdb.RPush(ctx, "jobs", jobID).Result()
	return err
}

// FIFO manner
func (q *queue) GetJobID() (string, error) {

	// for later use
	// ctx, cancel := context.WithTimeout(q.ctx, 2*time.Second)
	// defer cancel()

	// TODO change to stream thing idk like BLop something
	jobID, err := q.rdb.LPop(context.Background(), "jobs").Result()
	if err != nil {
		return "", err
	}
	return jobID, nil
}

func New() Queue {
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
