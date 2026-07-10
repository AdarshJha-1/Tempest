package store

const (
	CREATE_JOBS_TABLE_STMT = `
		CREATE TABLE IF NOT EXISTS Jobs (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			status TEXT NOT NULL,
			config BLOB NOT NULL,
			created_at DATETIME,
			started_at DATETIME,
			finished_at DATETIME
		)
	`

	INSERT_JOB_STMT = `
		INSERT INTO Jobs (
			id,
			name, 
			status,
			config, 
			created_at
		) VALUES (?, ?, ?, ?, ?)
	`

	SELECT_JOB_CONFIG_BY_ID_STMT = `
		SELECT config
		FROM Jobs
		WHERE id = ?
	`

	UPDATE_JOB_STATUS_BY_ID_STMT = `
		UPDATE Jobs
		SET 
			status = ?,
			started_at = ?
		WHERE id = ?
	`

	UPDATE_JOB_FINISH_TIME_BY_ID_STMT = `
		UPDATE Jobs
		SET finished_at = ?
		WHERE id = ?
	`

	GET_ALL_JOB = `
		SELECT *
		FROM Jobs
	`
	GET_ALL_JOB_CONFIG = `
		SELECT config
		FROM Jobs
	`
	CLEAN_DB = `
		DELETE
		FROM Jobs
	`
)

const (
	CREATE_RESULT_TABLE_STMT = `
		CREATE TABLE IF NOT EXISTS Results (
			job_id TEXT PRIMARY KEY NOT NULL,
			total_requests INTEGER,
			success_2xx INTEGER,
			client_4xx INTEGER,
			server_5xx INTEGER,
			network_errors INTEGER,
			total_latency INTEGER,
			
			FOREIGN KEY (job_id) REFERENCES Jobs(id) ON DELETE CASCADE
		)
	`

	INSERT_RESULT_STMT = `
		INSERT INTO Results (
			job_id,
			total_requests,
			success_2xx,
			client_4xx,
			server_5xx,
			network_errors,
			total_latency
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	SELECT_ALL_RESULT_STMT = `
		SELECT *
		FROM Results
	`
)
