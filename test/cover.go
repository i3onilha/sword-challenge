package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting server on port 6767")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file := "./test/coverage.html"
		http.ServeFile(w, r, file)
	})
	log.Fatal(http.ListenAndServe(":6767", nil))
}
