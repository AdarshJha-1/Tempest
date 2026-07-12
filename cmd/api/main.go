package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AdarshJha-1/Tempest/internal/config"
	"github.com/AdarshJha-1/Tempest/internal/queue"
	"github.com/AdarshJha-1/Tempest/internal/store"
	"github.com/AdarshJha-1/Tempest/testdata"
	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	usrConfig := &config.Config{}
	err := yaml.Unmarshal([]byte(testdata.YmlData), usrConfig)
	check(err)

	errs := config.Validate(usrConfig)
	check(errs)

	que := queue.New()
	resStr, err := que.Ping()
	check(err)
	fmt.Println("PING REDIS OKK", resStr)

	store, err := store.New()
	check(err)
	defer store.Close()

	err = store.Ping()
	check(err)
	fmt.Println("PING DB OKK")

	err = store.Init()
	check(err)

	configByte, err := json.Marshal(usrConfig)
	check(err)

	jobId, err := store.CreateJob(usrConfig.Name, configByte)
	check(err)

	err = que.PushJobID(jobId)
	check(err)
}
