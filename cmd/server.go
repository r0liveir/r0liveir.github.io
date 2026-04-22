// Simple tool for local server.

package main 

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("../"))
	http.Handle("/", fs)

	log.Println("Server started on http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
