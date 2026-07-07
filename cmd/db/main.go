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
	if arg == "st" { // show table
		jobs, err := store.ListAllJob()
		check(err)
		for _, j := range jobs {
			fmt.Println(j.Name)
			fmt.Println(j.ID)
			fmt.Println(j.Status)
		}
	} else if arg == "cl" { // remove all
		err := store.Clean()
		check(err)
	}

}
