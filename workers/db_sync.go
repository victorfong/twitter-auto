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

  for true {
    time.Sleep(w.SleepTime)
    go w.syncFollowings()
    go w.syncFollowers()
  }

}

func (w *DatabaseSyncWorker) syncFollowings() error{
  log.Print("Running Sync Following")
  ids, err := w.Twitter.GetSelfFriendIds()
  if err != nil {
    log.Printf("Found error while sync'ing following: %v", err)
    return err
  }

  log.Printf("Found %d following", len(ids))

  err = w.Database.SyncFollowings(ids)

  if err != nil {
    log.Printf("Found error while sync'ing following: %v", err)
    return err
  }

  log.Print("Finished Running Sync Following")
  return nil
}

func (w *DatabaseSyncWorker) syncFollowers() error {
  log.Print("Running Sync Followers")

  ids, err := w.Twitter.GetSelfFollowerIds()
  if err != nil {
    log.Printf("Found error while sync'ing follower: %v", err)
    return err
  }

  log.Printf("Found %d followers", len(ids))
  err = w.Database.SyncFollowers(ids)
  if err != nil {
    log.Printf("Found error while sync'ing follower: %v", err)
    return err
  }

  log.Print("Finished Running Sync Followers")
  return nil
}
