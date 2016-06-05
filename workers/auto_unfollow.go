package workers

import (
  "github.com/victorfong/twitter-auto/social"
  "github.com/victorfong/twitter-auto/model"
  "time"
  "log"
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

  log.Printf("Unfollowing %v", unfollowIds)

  if len(unfollowIds) > 0 {
    err = w.Twitter.Unfollow(unfollowIds[0])
    if err != nil {
      log.Printf("Error while unfollowing: %v", err)
      return err
    }

    err = w.Database.Unfollow(unfollowIds)
    if err != nil {
      log.Printf("Error while unfollowing: %v", err)
      return err
    }
  }

  log.Printf("Finished unfollowing %v", unfollowIds)

  return nil
}
