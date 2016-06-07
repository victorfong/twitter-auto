package workers

import (
  "github.com/victorfong/twitter-auto/social"
  "github.com/victorfong/twitter-auto/model"
  "time"
  "log"
  "math/rand"
)

type AutoUnfollowWorker struct {
  SleepTime time.Duration
  Twitter social.Twitter
  Database model.Database
}

func (w AutoUnfollowWorker) Start() {

  for true {
    go w.unfollow()
    time.Sleep(w.SleepTime)
  }

}

func (w AutoUnfollowWorker) unfollow() error{

  unfollowIds, err := w.Database.GetUnfollowList()
  if err != nil {
    log.Printf("Error while unfollowing: %v", err)
    return err
  }

  if len(unfollowIds) > 50 {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    unfollowerId := unfollowIds[r.Intn(len(unfollowIds))]

    log.Printf("Unfollowing %v", unfollowerId)
    err = w.Twitter.Unfollow()
    if err != nil {
      log.Printf("Error while unfollowing: %v", err)
      return err
    }

    t := make([]int64, 1)
    t[0] = unfollowerId
    err = w.Database.Unfollow(t)
    if err != nil {
      log.Printf("Error while unfollowing: %v", err)
      return err
    }
  } else {
    log.Printf("Not enough people to unfollow\n")
  }

  log.Printf("Finished unfollowing %v", unfollowIds)

  return nil
}
