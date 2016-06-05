package workers
import (
  . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  "time"
)

var _ = Describe("Unit", func(){
  Context("Testing channels", func(){
    It("runs", func(){
      w := AutoFollowWorker{
        SleepTime: 1 * time.Nanosecond,
      }
      channel := w.RunAutoFollow()
      var result int64 = 0
      for i := 0; i < 10 ;i++ {
        result = <- channel
      }
      Expect(result).ToNot(Equal(0))
    })
  })
})
