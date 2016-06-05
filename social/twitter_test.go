package social
import (
  . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Integration", func(){
  Context("When given valid twitter credentials", func(){
    It("establishes an API after Init", func(){
      twitter := TwitterConnection{}
      err := twitter.Init()
      Expect(err).To(BeNil())
      Expect(twitter.api).ToNot(BeNil())
    })
  })
})
