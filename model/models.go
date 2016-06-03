package model
import (
  "time"
)

type Following struct {
  Id          int64       `db:"twitter_id"`
  Since time.Time `db:"since"`
}

type Follower struct {
  Id int64 `db:"id"`
  TwitterId   int64 `db:"twitter_id"`
}
