package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/AdarshJha-1/Tempest/internal/config"
	"github.com/AdarshJha-1/Tempest/internal/queue"
	"github.com/AdarshJha-1/Tempest/internal/store"
	"github.com/AdarshJha-1/Tempest/testdata"
	"github.com/goccy/go-yaml"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

// TODO have to do better :/
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// ugly aah code :)
func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from system environment")
	}

	var usrConfig config.Config
	err := yaml.Unmarshal([]byte(testdata.YmlData), &usrConfig)
	check(err)

	que := queue.New()
	resStr, err := que.Ping()
	check(err)
	fmt.Println("PING REDIS OKK", resStr)

	store, err := store.New(que)
	check(err)
	defer store.Close()

	err = store.Ping()
	check(err)
	fmt.Println("PING DB OKK")

	err = store.Init()
	check(err)

	configByte, err := json.Marshal(usrConfig)
	check(err)

	err = store.Insert(usrConfig.Name, configByte)
	check(err)

	time.Sleep(10 * time.Second)
	jobs, err := store.ListAllJob()
	check(err)
	for _, job := range jobs {
		fmt.Println(job.Name)
		fmt.Println(job.Status)
		fmt.Println(job.FinishedAt)
	}
	// TODO causing error OMG
	// currJobID, err := que.GetJobID()
	// check(err)

	// configData, err := store.GetJobConfigByID(currJobID)
	// check(err)

	// var cfg config.Config
	// err = json.Unmarshal(configData, &cfg)
	// check(err)

	// pkg.PrettyPrintJSON(cfg)
}
