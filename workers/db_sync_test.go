package workers
import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/golang/mock/gomock"
  "github.com/victorfong/twitter-auto/mock"
)

var _ = Describe("Unit", func(){
  Context("When there are new follower", func(){
    It("persists them in database", func(){
      ctrl := gomock.NewController(GinkgoT())
	    defer ctrl.Finish()

      twitter := mock.NewMockTwitter(ctrl)
      twitterFollowers := make([]int64, 2)
      twitterFollowers[0] = 0
      twitterFollowers[1] = 1

      twitter.EXPECT().GetSelfFollowerIds().Return(twitterFollowers, nil)

      database := mock.NewMockDatabase(ctrl)
      database.EXPECT().SyncFollowers(twitterFollowers).Return(nil)

      worker := DatabaseSyncWorker{
        Twitter: twitter,
        Database: database,
      }

      err := worker.syncFollowers()
      Expect(err).To(BeNil())

    })
  })
})
