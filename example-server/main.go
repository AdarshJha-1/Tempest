package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	server := http.NewServeMux()
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Conn", r.URL)
		w.WriteHeader(http.StatusOK)
	})
	fmt.Println("Server is running!")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatal(err)
	}
}
