package model
import (
  . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  "os"
  "log"
  "math/rand"
)

var _ = Describe("Unit", func(){
  Context("when VCAP_SERVICES is set", func(){
    Context("when MySQL Service URI is set", func(){
      BeforeEach(func() {
        envVar := "{\"cleardb\": [{\"credentials\": {\"uri\": \"test_string\"}}]}"
        os.Setenv("VCAP_SERVICES", envVar)
      })

      It("returns MySQL Service URI", func(){
        uri, err := getURI()
        Expect(err).To(BeNil())
        Expect(uri).To(Equal("test_string"))
      })

    })

    Context("when MySQL Service URI is not set", func(){
      BeforeEach(func() {
        envVar := "{\"something\":\"\"}"
        os.Setenv("VCAP_SERVICES", envVar)
      })

      It("returns error", func(){
        _, err := getURI()
        Expect(err).ToNot(BeNil())
      })
    })
  })

  Context("when VCAP_SERVICES is not set", func(){
    BeforeEach(func() {
      os.Unsetenv("VCAP_SERVICES")
    })

    It("returns error", func(){
      _, err := getURI()
      Expect(err).ToNot(BeNil())
    })
  })
})

var _ = Describe("Integration", func(){
  Context("when database connection is valid", func(){
    BeforeEach(func() {
      envVar := `{
        "cleardb": [
          {
            "credentials": {
              "uri": "mysql://twitter-auto2:twitter-auto2@localhost:3306/twitter-auto2?reconnect=true"
            }
          }
        ]
      }`
      os.Setenv("VCAP_SERVICES", string(envVar))
    })

    var db DatabaseConnection
    It("creates a database connection", func(){
      db = DatabaseConnection{}
      err := db.Init()
      Expect(err).To(BeNil())
    })

    It("can insert a new follower", func() {
      db = DatabaseConnection{}
      err := db.Init()
      Expect(err).To(BeNil())

      var follower Follower = Follower{ TwitterId: 153 }
      log.Printf("follower.TwitterId = %d", follower.TwitterId)
      err = db.InsertFollower(&follower)
      Expect(err).To(BeNil())
    })

    It("can insert new temp followings", func() {
      db = DatabaseConnection{}
      err := db.Init()
      Expect(err).To(BeNil())

      r := rand.New(rand.NewSource(99))
      ids := make([]int64, 2)
      ids[0] = r.Int63()
      ids[1] = r.Int63()

      err = db.insertTempFollowings(ids)
      Expect(err).To(BeNil())

      err = db.clearTempFollowings()
      Expect(err).To(BeNil())
    })

    It("can sync followers", func(){
      db = DatabaseConnection{}
      err := db.Init()
      Expect(err).To(BeNil())

      r := rand.New(rand.NewSource(99))
      ids := make([]int64, 2)
      ids[0] = r.Int63()
      ids[1] = r.Int63()

      err = db.SyncFollowers(ids)
      Expect(err).To(BeNil())

      err = db.clearFollowers()
      Expect(err).To(BeNil())
    })

  })

  Context("when database connection is invalid", func(){
    BeforeEach(func() {
      envVar := `{
        "cleardb": [
          {
            "credentials": {
              "uri": "mysql://whatever:lalala@localhost:3306/whatever?reconnect=true"
            }
          }
        ]
      }`
      os.Setenv("VCAP_SERVICES", string(envVar))
    })

    It("returns an error", func(){
      db := DatabaseConnection{}
      err := db.Init()

      Expect(err).ToNot(BeNil())
    })
  })




})
