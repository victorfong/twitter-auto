package social
import (
  . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func(){
  Context("When given valid twitter credentials", func(){
    var twitter TwitterConnection

    BeforeEach(func() {
      twitter = NewTwitter()
    })

    It("establishes an API after Init", func(){
      Expect(twitter.api).ToNot(BeNil())
    })

    It("can return list of friend Ids", func(){
      ids, err := twitter.GetSelfFriendIds()
      Expect(err).To(BeNil())
      Expect(len(ids) > 0).To(BeTrue())
    })

    It("can return list of follower ids", func(){
      friendIds, err := twitter.GetSelfFriendIds()
      ids, err := twitter.GetFollowerIds(friendIds[0])
      Expect(err).To(BeNil())
      Expect(len(ids) > 0).To(BeTrue())
    })
  })
})
