package social
import (
  . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func(){
  Context("When given valid twitter credentials", func(){
    It("establishes an API after Init", func(){
      twitter := TwitterConnection{}
      twitter.Init()
      Expect(twitter.api).ToNot(BeNil())
    })

    It("can return list of friend Ids", func(){
      twitter := TwitterConnection{}
      twitter.Init()
      ids, err := twitter.getSelfFriendIds()
      Expect(err).To(BeNil())
      Expect(len(ids) > 0).To(BeTrue())

    })
  })
})
