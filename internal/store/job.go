package store

import (
	"context"
	"log"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/job"
	"github.com/google/uuid"
)

func (s *store) CreateJob(name string, configByte []byte) (string, error) {

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

func (s *store) GetJobConfig(jobId string) ([]byte, error) {

	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	var configData []byte
	err := s.db.QueryRowContext(ctx, SELECT_JOB_CONFIG_BY_ID_STMT,
		jobId,
	).Scan(&configData)
	return configData, err
}

func (s *store) UpdateJobStatus(jobId string, status string) error {
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

func (s *store) UpdateJobFinishTime(jobId string) error {
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

func (s *store) ListJobs() ([]job.Job, error) {
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
