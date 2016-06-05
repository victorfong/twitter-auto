package workers
import (
  . "github.com/onsi/ginkgo"
  . "github.com/onsi/gomega"

  "github.com/golang/mock/gomock"
  "github.com/victorfong/twitter-auto/mock"
)

var _ = Describe("Unit", func(){
  Context("When there are new people who need to be unfollowed", func(){
    It("unfollow them one at a time", func(){
      ctrl := gomock.NewController(GinkgoT())
	    defer ctrl.Finish()

      unfollow := make([]int64, 2)
      unfollow[0] = 1234123

      database := mock.NewMockDatabase(ctrl)
      database.EXPECT().GetUnfollowList().Return(unfollow, nil)
      database.EXPECT().Unfollow(unfollow).Return(nil)

      twitter := mock.NewMockTwitter(ctrl)
      twitter.EXPECT().Unfollow(unfollow[0])

      worker := AutoUnfollowWorker{
        Twitter: twitter,
        Database: database,
      }

      err := worker.unfollow()
      Expect(err).To(BeNil())

    })
  })

  // Context("When there are new following", func(){
  //   It("sync them in database", func(){
  //     ctrl := gomock.NewController(GinkgoT())
	//     defer ctrl.Finish()
  //
  //     twitter := mock.NewMockTwitter(ctrl)
  //     twitterFollowers := make([]int64, 2)
  //     twitterFollowers[0] = 0
  //     twitterFollowers[1] = 1
  //
  //     twitter.EXPECT().GetSelfFriendIds().Return(twitterFollowers, nil)
  //
  //     database := mock.NewMockDatabase(ctrl)
  //     database.EXPECT().SyncFollowers(twitterFollowers).Return(nil)
  //
  //     worker := DatabaseSyncWorker{
  //       Twitter: twitter,
  //       Database: database,
  //     }
  //
  //     err := worker.syncFollowers()
  //     Expect(err).To(BeNil())
  //
  //   })
  // })
})
