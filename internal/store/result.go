package store

import (
	"context"
	"log"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/metrics"
)

func (s *store) CreateResult(jobId string, result *metrics.Result) error {

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

func (s *store) ListResults() ([]metrics.Result, error) {
	ctx, cancel := context.WithTimeout(s.ctx, 2*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, SELECT_ALL_RESULT_STMT)

	if err != nil {
		log.Printf("%q: %s\n", err, SELECT_ALL_RESULT_STMT)
		return nil, err
	}
	defer rows.Close()

	// TODO i have to handle jobId idk how tho
	var jobId string
	var results []metrics.Result
	for rows.Next() {
		r := metrics.Result{}
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
