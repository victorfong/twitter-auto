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
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  for true {
    go w.unfollow()
    time.Sleep(time.Duration(r.Intn(15) + 1) * w.SleepTime)
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
    err = w.Twitter.Unfollow(unfollowerId)
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

    log.Printf("Finished unfollowing %v", unfollowerId)
  } else {
    log.Printf("Not enough people to unfollow\n")
  }

  return nil
}
