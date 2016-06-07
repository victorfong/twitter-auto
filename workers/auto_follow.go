package workers

import (
  "github.com/victorfong/twitter-auto/social"
  "github.com/victorfong/twitter-auto/model"

  "time"
  "log"
  "math/rand"
)

type AutoFollowWorker struct {
  SleepTime time.Duration
  Twitter social.Twitter
  Database model.Database
}

func (w AutoFollowWorker) Start() {
  channel := make(chan int64)
  go w.findCandidate(channel)

  for true {
    go w.follow(channel)
    time.Sleep(w.SleepTime)
  }
}

func (w AutoFollowWorker) findCandidate(channel chan int64) {
  r := rand.New(rand.NewSource(time.Now().UnixNano()))
  for true {
    followerIds, err := w.Twitter.GetSelfFollowerIds()
    if err != nil {
      log.Printf("Error while finding candidate: %v", err)
      time.Sleep(w.SleepTime)
    }

    if len(followerIds) > 0 {
      index := r.Intn(len(followerIds))
      candidateIds, err := w.Twitter.GetFollowerIds(followerIds[index])
      if err != nil {
        log.Printf("Error while getting follower id: %v", err)
        time.Sleep(w.SleepTime)
      }

      for i := 0; i < 10 && i < len(candidateIds); i++ {
        index := r.Intn(len(candidateIds))
        candidateId := candidateIds[index]

        alreadyFollowed, err := w.Database.HasAlreadyFollowed(candidateId)

        if err == nil {
          if(alreadyFollowed){
              log.Printf("Already followed candidate %d", candidateId)
          } else {
            log.Printf("Adding candidate id %d to queue", candidateId)
            channel <- candidateId
          }
        } else {
          log.Printf("Error while finding candidate: %v", err)
        }

      }
    }
  }
}

func (w AutoFollowWorker) follow(channel chan int64) error{
  candidateId := <- channel
  err := w.Twitter.Follow(candidateId)
  if err != nil {
    log.Printf("Error while following candidate: %v", err)
    return err
  }

  candidateIds := make([]int64, 1)
  candidateIds[0] = candidateId

  err = w.Database.InsertFollowings(candidateIds)
  if err != nil {
    log.Printf("Error while following candidate: %v", err)
    return err
  }

  return nil
}
