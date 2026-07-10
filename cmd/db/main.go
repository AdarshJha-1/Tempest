package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AdarshJha-1/Tempest/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	store, err := store.New()
	check(err)

	arg := os.Args[1]
	switch arg {
	case "st":
		jobs, err := store.ListJobs()
		check(err)
		for _, j := range jobs {
			fmt.Println(j.Name)
			fmt.Println(j.ID)
			fmt.Println(j.Status)
		}
	case "r":
		err := store.Reset()
		check(err)
	case "sr":
		results, err := store.ListResults()
		check(err)
		fmt.Println("HERE")
		fmt.Println(results)
		for _, r := range results {
			fmt.Println(r.TotalLatency)
			fmt.Println(r.TotalRequests)
			fmt.Println(r.Success2xx)
			fmt.Println(r.Client4xx)
			fmt.Println(r.Server5xx)
			fmt.Println(r.NetworkErrors)
		}
	}
}
