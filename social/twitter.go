package social

import (
  "github.com/ChimeraCoder/anaconda"
  "os"
  "log"
  "strconv"
  "net/url"
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

func (t TwitterConnection) getSelfFriendIds() ([]int64, error) {
  strUserId := os.Getenv("TWITTER_USER_ID")
  userId, err := strconv.ParseInt(strUserId, 10, 64)
  if err != nil {
    return nil, err
  }

  v := url.Values{}
  v.Set("count", "5000")

  c, err := t.api.GetFriendsUser(userId, v)
  if err != nil {
    return nil, err
  }

  return c.Ids, nil
}
//
// func (t TwitterConnection) getFollowerIds(userId int64) ([]int64, error) {
//
// }
