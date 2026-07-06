package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/queue"
	"github.com/google/uuid"
)

type Store interface {
	Ping() error
	Init() error
	Close() error
	Insert(name string, configByte []byte) error
	GetByID(jobId string) ([]byte, error)
}

type store struct {
	ctx   context.Context
	db    *sql.DB
	queue queue.Queue
}

func New(queue queue.Queue) (Store, error) {
	newDB, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, err
	}
	return &store{
		ctx:   context.Background(),
		db:    newDB,
		queue: queue,
	}, nil
}

func (s *store) Ping() error {

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()
	return s.db.PingContext(ctx)
}

func (s *store) Init() error {
	_, err := s.db.Exec(CREATE_TABLE_STMT)
	if err != nil {
		log.Printf("%q: %s\n", err, CREATE_TABLE_STMT)
		return err
	}
	return nil
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) Insert(name string, configByte []byte) error {

	newUUID := uuid.NewString()
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(s.ctx,
		INSERT_JOB_STMT,
		newUUID,
		name,
		"pending",
		configByte,
		time.Now(),
	)

	if err != nil {
		fmt.Println("Transaction failed on step 1:", err)
		return err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		log.Println("Commit failed")
		return err
	}

	err = s.queue.PushJobID(newUUID)
	if err != nil {
		log.Println("Redis -> failed to push job")
		return err
	}
	return nil
}

func (s *store) GetByID(jobId string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	var configData []byte
	err := s.db.QueryRowContext(ctx, SELECT_JOB_BY_ID_STMT,
		jobId,
	).Scan(&configData)
	return configData, err
}
