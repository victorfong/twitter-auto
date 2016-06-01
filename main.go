package main

import (
	"fmt"
	// "html"
	"log"
	"net/http"
	// "net/url"
	"os"
	"time"
)

func heartbeat() {
	for true {
		log.Printf("Hello World!")

		time.Sleep(1 * time.Second)
	}
}

func main() {
	go heartbeat()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
