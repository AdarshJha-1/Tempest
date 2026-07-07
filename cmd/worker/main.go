package main

import (
	"fmt"
	"log"

	"github.com/AdarshJha-1/Tempest/internal/executor"
	"github.com/AdarshJha-1/Tempest/internal/queue"
	"github.com/AdarshJha-1/Tempest/internal/store"
	"github.com/AdarshJha-1/Tempest/internal/worker"
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

	executor := executor.New(30)

	worker := worker.New(que, store, executor, 2)

	worker.Start()

}
