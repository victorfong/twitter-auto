package main

import (
	"fmt"
	"time"
	// "html"
	"log"
	"net/http"
	// "net/url"
	"os"

	"github.com/victorfong/twitter-auto/model"
	"github.com/victorfong/twitter-auto/social"
	"github.com/victorfong/twitter-auto/workers"
	// "github.com/ChimeraCoder/anaconda"
)

func startWorkers() {

	twitter := social.NewTwitter()

	database := model.DatabaseConnection{}
	err := model.InitDatabase(&database)
	if err != nil {
		log.Printf("Error while initializing database: %v\n", err)
	}

	dbSyncWorker := workers.DatabaseSyncWorker{
		SleepTime: time.Hour,
		Twitter:   twitter,
		Database:  database,
	}

	go dbSyncWorker.Start()

	unfollowWorker := workers.AutoUnfollowWorker{
		SleepTime: time.Minute,
		Twitter: twitter,
		Database: database,
	}

	go unfollowWorker.Start()

	followWorker := workers.AutoFollowWorker{
		SleepTime: time.Minute,
		Twitter:   twitter,
		Database:  database,
	}

	go followWorker.Start()
}

func main() {
	go startWorkers()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
