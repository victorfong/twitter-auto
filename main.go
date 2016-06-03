package main

import (
	"fmt"
	// "html"
	"log"
	"net/http"
	// "net/url"
	"os"
	"github.com/victorfong/twitter-auto/model"
	// "time"
	// "github.com/ChimeraCoder/anaconda"
)



func heartbeat() {
	// log.Printf("Initializing Anaconda endpoint")
	// anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	// anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	// api := anaconda.NewTwitterApi(os.Getenv("TWITTER_TOKEN"), os.Getenv("TWITTER_TOKEN_SECRET"))
	// log.Printf("Finished Initializing")
	//
	// log.Printf("Getting Follower List")
	// pages := api.GetFollowersListAll(nil)
	// log.Printf("Finished Getting Follower List")
	// i := 0
	// for page := range pages {
	// 	if page.Error != nil {
	// 		log.Printf("Error caught in follower page")
	// 		fmt.Println(page.Error)
	// 	}
	//
	//
	// 	for _, user := range page.Followers {
	//
	//
	// 		log.Printf("Follower: %d", i)
	// 		i = i + 1
	// 		fmt.Printf("ID: %d\n", user.Id)
	// 		fmt.Printf("Name: %s\n", user.Name)
	// 		fmt.Printf("Following: %t\n", user.Following)
	// 	}
	// }
	//
	// i = 0
	// friendPages := api.GetFriendsListAll(nil)
	// for friendPage := range friendPages {
	// 	if friendPage.Error != nil {
	// 		log.Printf("Error caught in follower friendPage")
	// 		fmt.Println(friendPage.Error)
	// 	}
	//
	// 	for _, friend := range friendPage.Friends {
	// 		log.Printf("Friend: %d", i)
	// 		i = i + 1
	// 		fmt.Printf("ID: %d\n", friend.Id)
	// 		fmt.Printf("Name: %s\n", friend.Name)
	// 		fmt.Printf("Following: %t\n", friend.Following)
	// 	}
	// }
	// log.Printf("Finished Twitter Calls")
	database := model.DatabaseConnection{}
	database.Init("Something URI")
	//
	//
	// for true {
	// 	log.Printf("Hello World!")
	// 	time.Sleep(1 * time.Second)
	// }
}

func main() {
	go heartbeat()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello")
	})
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
