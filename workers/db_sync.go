package workers

import (
  "github.com/victorfong/twitter-auto/social"
  "github.com/victorfong/twitter-auto/model"
  "time"
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
  // ids, err := w.Twitter.GetSelfFriendIds()
  // if err != nil {
  //   return err
  // }
  //
  //
  return nil
}

func (w *DatabaseSyncWorker) syncFollowers() error {
  ids, err := w.Twitter.GetSelfFollowerIds()
  if err != nil {
    return err
  }

  err = w.Database.SyncFollowers(ids)
  return err
}
