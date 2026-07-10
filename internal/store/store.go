package store

import (
	"context"
	"database/sql"

	"log"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/job"
	"github.com/AdarshJha-1/Tempest/internal/metrics"
)

type Store interface {
	Ping() error
	Init() error
	Close() error
	CreateJob(name string, configByte []byte) (string, error)
	GetJobConfig(jobId string) ([]byte, error)
	UpdateJobStatus(jobId string, status string) error
	UpdateJobFinishTime(jobId string) error
	ListJobs() ([]job.Job, error)

	CreateResult(jobId string, result *metrics.Result) error
	ListResults() ([]metrics.Result, error)

	Reset() error
}

type store struct {
	ctx context.Context
	db  *sql.DB
}

func New() (Store, error) {
	newDB, err := sql.Open("sqlite3", "store.db")
	if err != nil {
		return nil, err
	}
	return &store{
		ctx: context.Background(),
		db:  newDB,
	}, nil
}

func (s *store) Ping() error {

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()
	return s.db.PingContext(ctx)
}

func (s *store) Init() error {

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx,
		CREATE_JOBS_TABLE_STMT,
	)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx,
		CREATE_RESULT_TABLE_STMT,
	)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println("Commit failed")
		return err
	}
	return nil
}

func (s *store) Close() error {
	return s.db.Close()
}

func (s *store) Reset() error {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, DROP_JOBS_TABLE)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, DROP_RESULT_TABLE)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Println("Commit failed")
		return err
	}
	return nil
}
