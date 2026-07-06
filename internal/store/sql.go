package store

const (
	CREATE_TABLE_STMT = `
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

	SELECT_JOB_BY_ID_STMT = `
		SELECT config
		FROM Jobs
		WHERE id = ?
	`
)
