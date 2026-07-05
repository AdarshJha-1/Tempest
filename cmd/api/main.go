package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/config"
	"github.com/AdarshJha-1/Tempest/internal/queue"
	"github.com/goccy/go-yaml"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

// TODO have to do better :/
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func init() {

}

func main() {

	// os.Remove("./testdata/test.db")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	data, err := os.ReadFile("./testdata/test.yml")
	check(err)

	var usrConfig config.Config
	err = yaml.Unmarshal(data, &usrConfig)
	check(err)

	db, err := sql.Open("sqlite3", "./testdata/test.db")
	check(err)
	defer db.Close()

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS Jobs (
			id TEXT NOT NULL PRIMARY KEY,
			name TEXT NOT NULL,
			status TEXT NOT NULL,
			config BLOB NOT NULL,
			created_at DATETIME,
			started_at DATETIME,
			finished_at DATETIME
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	que := queue.NewQueue()
	str, err := que.Ping()
	check(err)
	fmt.Println("PING ->", str)

	configJSON, err := json.Marshal(usrConfig)
	check(err)

	_, err = db.Exec(`
		INSERT INTO Jobs (
			id,
			name, 
			status,
			config, 
			created_at
		) VALUES (?, ?, ?, ?, ?)
	`,
		uuid.NewString(),
		usrConfig.Name,
		"pending",
		configJSON,
		time.Now(),
	)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	var configData []byte

	err = db.QueryRow(`
		SELECT config
		FROM Jobs
	`).Scan(&configData)
	check(err)

	var cfg config.Config
	err = json.Unmarshal(configData, &cfg)
	check(err)

	// pkg.PrettyPrintJSON(cfg)
}
