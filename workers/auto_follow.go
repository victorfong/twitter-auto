package workers

import (
  "time"
  "log"
)

type AutoFollowWorker struct {
  SleepTime time.Duration

}

func (w AutoFollowWorker) RunAutoFollow() <-chan int64 {
  channel := make(chan int64)
  resultChannel := make(chan int64)
  go w.findCandidate(channel)
  go w.consumeCandidate(channel, resultChannel)
  return resultChannel
}

func (w AutoFollowWorker) findCandidate(channel chan int64) {
  var i int64 = 0
  for true {
    log.Printf("Putting %d into channel", i)
    channel <- i
    i = i + 2
  }
}

func (w AutoFollowWorker) consumeCandidate(channel chan int64, resultChannel chan int64) {
  for true {
    candidateId := <- channel
    log.Printf("Consumed %d from channel", candidateId)
    time.Sleep(w.SleepTime)
    resultChannel <- candidateId
  }
}
