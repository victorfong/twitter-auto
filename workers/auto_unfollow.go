package workers

import (
  "github.com/victorfong/twitter-auto/social"
  "github.com/victorfong/twitter-auto/model"
  "time"
  "log"
  "math/rand"
  "regexp"
  "strconv"
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
    time.Sleep(time.Duration(r.Intn(80) + 1) * time.Second + w.SleepTime)
    // time.Sleep(w.SleepTime)
  }

}

func IsEnglishName(name string) bool{
  match, _ := regexp.MatchString("[A-Za-z\\s\\'\\,\\-]{"+strconv.Itoa(len(name))+"}", name)
  return match
}

func (w AutoUnfollowWorker) unfollow() error{

  unfollowIds, err := w.Database.GetUnfollowList()
  if err != nil {
    log.Printf("Error while unfollowing: %v", err)
    return err
  }

  if len(unfollowIds) > 50 {
    var unfollowerId int64
    for i := 0; i < 20; i++ {
      r := rand.New(rand.NewSource(time.Now().UnixNano()))
      unfollowerId = unfollowIds[r.Intn(len(unfollowIds) - 5) + 5]
      userShow, err := w.Twitter.GetUsersShowById(unfollowerId)
      if err != nil{
        log.Printf("Received error while fetching user show for %d. Skipping...", unfollowerId)
      } else {
        if(!IsEnglishName(userShow.Name)){
          log.Printf("Candidate name %s appears not to be English", userShow.Name)
          log.Printf("Attempting to unfollow %d %s.", unfollowerId, userShow.Name)
          break;
        }
      }
    }

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
