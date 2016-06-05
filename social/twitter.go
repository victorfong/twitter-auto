package social

import (
  "github.com/ChimeraCoder/anaconda"
  "os"
  "log"
)

type Twitter interface {
  Init()
}

type TwitterConnection struct {
  api *anaconda.TwitterApi
}

func (t *TwitterConnection) Init() {
  log.Printf("Initializing Anaconda endpoint")
	anaconda.SetConsumerKey(os.Getenv("TWITTER_CONSUMER_KEY"))
	anaconda.SetConsumerSecret(os.Getenv("TWITTER_CONSUMER_SECRET"))
	t.api = anaconda.NewTwitterApi(os.Getenv("TWITTER_TOKEN"), os.Getenv("TWITTER_TOKEN_SECRET"))
	log.Printf("Finished Initializing")
}
