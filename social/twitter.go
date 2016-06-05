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

  GetSelfFriendIds() ([]int64, error)
  GetFriendIds(userId int64) ([]int64, error)

  GetSelfFollowerIds() ([]int64, error)
  GetFollowerIds(userId int64) ([]int64, error)
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

func (t TwitterConnection) GetSelfFriendIds() ([]int64, error) {
  userId, err := getCurrentUserId()
  if err != nil {
    return nil, err
  }
  return t.GetFriendIds(userId)
}

func (t TwitterConnection) GetFriendIds(userId int64) ([]int64, error) {
  strUserId := strconv.FormatInt(userId, 10)

  v := url.Values{}
  v.Set("user_id", strUserId)
  v.Set("count", "5000")

  c, err := t.api.GetFriendsIds(v)
  if err != nil {
    return nil, err
  }

  return c.Ids, nil
}

func (t TwitterConnection) GetSelfFollowerIds() ([]int64, error) {
  userId, err := getCurrentUserId()
  if err != nil {
    return nil, err
  }
  return t.GetFollowerIds(userId)
}

func (t TwitterConnection) GetFollowerIds(userId int64) ([]int64, error) {
  strUserId := strconv.FormatInt(userId, 10)

  v := url.Values{}
  v.Set("user_id", strUserId)
  v.Set("count", "5000")

  c, err := t.api.GetFollowersIds(v)
  if err != nil {
    return nil, err
  }

  return c.Ids, nil
}

func getCurrentUserId() (int64, error) {
  strUserId := os.Getenv("TWITTER_USER_ID")
  return strconv.ParseInt(strUserId, 10, 64)
}
