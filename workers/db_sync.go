package workers

import (
  "github.com/victorfong/twitter-auto/social"
  "github.com/victorfong/twitter-auto/model"
  "time"
  "log"
)

type DatabaseSyncWorker struct {
  SleepTime time.Duration
  Twitter social.Twitter
  Database model.Database
}

func (w *DatabaseSyncWorker) Start() {
  go w.syncFollowings()
}

func (w *DatabaseSyncWorker) syncFollowings() error{
  ids, err := w.Twitter.GetSelfFriendIds()
  if err != nil {
    return err
  }

  for _, id := range ids {
    log.Printf("id = %d\n", id)
  }
  return nil
}

func (w *DatabaseSyncWorker) syncFollowers() {

}
