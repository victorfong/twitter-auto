package workers
import (
  . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  // "time"
  "github.com/golang/mock/gomock"
  "github.com/victorfong/twitter-auto/mock"
)

var _ = Describe("Unit", func(){
  Context("Testing channels", func(){
    PIt("runs", func(){
      ctrl := gomock.NewController(GinkgoT())
	    defer ctrl.Finish()

      followerIds := make([]int64, 1)
      followerIds[0] = 1234123

      candidateIds := make([]int64, 2)
      candidateIds[0] = 98237

      database := mock.NewMockDatabase(ctrl)
      // database.EXPECT().GetUnfollowList().Return(unfollow, nil)
      // database.EXPECT().Unfollow(unfollow).Return(nil)
      database.EXPECT().HasAlreadyFollowed(candidateIds[0]).Return(false, nil)

      twitter := mock.NewMockTwitter(ctrl)
      twitter.EXPECT().GetSelfFollowerIds().Return(followerIds, nil)
      twitter.EXPECT().GetFollowerIds(followerIds[0]).Return(candidateIds, nil)

      worker := AutoFollowWorker{
        Twitter: twitter,
        Database: database,
      }

      channel := make(chan int64)
      go worker.findCandidate(channel)

      result := <- channel
      Expect(result).To(Equal(candidateIds[0]))
      // err := worker.follow(channel)
      // Expect(err).To(BeNil())
    })
  })
})
