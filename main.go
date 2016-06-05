package main

import (
	"fmt"
	// "html"
	"log"
	"net/http"
	// "net/url"
	"os"
	"github.com/victorfong/twitter-auto/model"
	"github.com/victorfong/twitter-auto/social"
	"github.com/victorfong/twitter-auto/workers"
	"time"
	// "github.com/ChimeraCoder/anaconda"
)



func startWorkers() {

	twitter := social.NewTwitter()

	database := model.DatabaseConnection{}
	model.InitDatabase(&database)

	dbSyncWorker := workers.DatabaseSyncWorker{
		SleepTime: time.Hour,
		Twitter: twitter,
		Database: database,
	}

	go dbSyncWorker.Start()
}

func main() {
	go startWorkers()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
