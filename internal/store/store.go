package store

import (
	"context"
	"database/sql"

	"log"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/config"
	"github.com/AdarshJha-1/Tempest/internal/job"
	"github.com/AdarshJha-1/Tempest/internal/metrics"
	"github.com/google/uuid"
)

type Store interface {
	Ping() error
	Init() error
	Close() error
	Insert(name string, configByte []byte) (string, error)
	GetJobConfigById(jobId string) ([]byte, error)
	UpdateJobStatusById(jobId string, status string) error
	UpdateJobFinishTimeById(jobId string) error
	ListAllJob() ([]job.Job, error)
	ListAllJobConfig() ([]config.Config, error)

	InsertResult(jobId string, result *metrics.Result) error
	ListAllResult() ([]*metrics.Result, error)

	Clean() error
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

func (s *store) Insert(name string, configByte []byte) (string, error) {

	newUUID := uuid.NewString()
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", err
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
		log.Println("Transaction failed on step 1:", err)
		return "", err
	}
	_, err = result.LastInsertId()
	if err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		log.Println("Commit failed")
		return "", err
	}
	return newUUID, nil
}

func (s *store) GetJobConfigById(jobId string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	var configData []byte
	err := s.db.QueryRowContext(ctx, SELECT_JOB_CONFIG_BY_ID_STMT,
		jobId,
	).Scan(&configData)
	return configData, err
}

func (s *store) UpdateJobStatusById(jobId string, status string) error {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, UPDATE_JOB_STATUS_BY_ID_STMT,
		status,
		time.Now(),
		jobId,
	)
	if err != nil {
		log.Printf("%q: %s\n", err, UPDATE_JOB_STATUS_BY_ID_STMT)
		return err
	}
	return nil
}

func (s *store) UpdateJobFinishTimeById(jobId string) error {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, UPDATE_JOB_FINISH_TIME_BY_ID_STMT,
		time.Now(),
		jobId,
	)
	if err != nil {
		log.Printf("%q: %s\n", err, UPDATE_JOB_FINISH_TIME_BY_ID_STMT)
		return err
	}
	return nil
}

func (s *store) ListAllJobConfig() ([]config.Config, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, GET_ALL_JOB_CONFIG)

	if err != nil {
		log.Printf("%q: %s\n", err, UPDATE_JOB_STATUS_BY_ID_STMT)
		return nil, err
	}
	defer rows.Close()

	var configs []config.Config
	for rows.Next() {
		var config config.Config
		err := rows.Scan(&config)
		if err != nil {
			return nil, err
		}

		configs = append(configs, config)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return configs, nil
}

func (s *store) ListAllJob() ([]job.Job, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, GET_ALL_JOB)

	if err != nil {
		log.Printf("%q: %s\n", err, UPDATE_JOB_STATUS_BY_ID_STMT)
		return nil, err
	}
	defer rows.Close()

	var jobs []job.Job
	for rows.Next() {
		var j job.Job
		err := rows.Scan(
			&j.ID,
			&j.Name,
			&j.Status,
			&j.ConfigJSON,
			&j.CreatedAt,
			&j.StartedAt,
			&j.FinishedAt,
		)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, j)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (s *store) Clean() error {

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, CLEAN_DB)
	return err
}

func (s *store) InsertResult(jobId string, result *metrics.Result) error {

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	_, err := s.db.ExecContext(ctx, INSERT_RESULT_STMT,
		jobId,
		result.TotalRequests,
		result.Success2xx,
		result.Client4xx,
		result.Server5xx,
		result.NetworkErrors,
		result.TotalLatency,
	)
	if err != nil {
		log.Println("failed in insert result", err)
	}
	return err
}

func (s *store) ListAllResult() ([]*metrics.Result, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, SELECT_ALL_RESULT_STMT)

	if err != nil {
		log.Printf("%q: %s\n", err, SELECT_ALL_RESULT_STMT)
		return nil, err
	}
	defer rows.Close()

	// this whole for temp db viewing thing so i will remove it eventually
	var jobId string
	var results []*metrics.Result
	for rows.Next() {
		r := &metrics.Result{}
		err := rows.Scan(
			&jobId,
			&r.TotalRequests,
			&r.Success2xx,
			&r.Client4xx,
			&r.Server5xx,
			&r.NetworkErrors,
			&r.TotalLatency,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, r)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
