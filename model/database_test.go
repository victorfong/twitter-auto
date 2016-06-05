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
              "uri": "mysql://twitter-auto-t1:twitter-auto-t1@localhost:3306/twitter-auto-t1?reconnect=true"
            }
          }
        ]
      }`
      os.Setenv("VCAP_SERVICES", string(envVar))
    })

    var db DatabaseConnection
    It("creates a database connection", func(){
      db = DatabaseConnection{}
      err := InitDatabase(&db)
      Expect(err).To(BeNil())
    })

    It("can insert a new follower", func() {
      db = DatabaseConnection{}
      err := InitDatabase(&db)
      Expect(err).To(BeNil())

      var follower Follower = Follower{ TwitterId: 153 }
      log.Printf("follower.TwitterId = %d", follower.TwitterId)
      err = db.InsertFollower(&follower)
      Expect(err).To(BeNil())
    })

    Context("When a user is already followed", func(){
      It("can get determine if a user is already followed", func(){
        db = DatabaseConnection{}
        err := InitDatabase(&db)
        Expect(err).To(BeNil())

        r := rand.New(rand.NewSource(99))
        ids := make([]int64, 2)
        ids[0] = r.Int63()

        result, err := db.HasAlreadyFollowed(ids[0])
        Expect(err).To(BeNil())
        Expect(result).To(BeFalse())

        err = db.InsertFollowings(ids)
        Expect(err).To(BeNil())

        result, err = db.HasAlreadyFollowed(ids[0])
        Expect(err).To(BeNil())
        Expect(result).To(BeTrue())

        err = db.clearFollowings()
        Expect(err).To(BeNil())
      })
    })

    It("can get unfollow list", func(){
      db = DatabaseConnection{}
      err := InitDatabase(&db)
      Expect(err).To(BeNil())

      r := rand.New(rand.NewSource(99))
      ids := make([]int64, 2)
      ids[0] = r.Int63()
      ids[1] = r.Int63()

      err = db.InsertFollowings(ids)
      Expect(err).To(BeNil())

      err = db.insertFollowers(ids[1:])
      Expect(err).To(BeNil())

      result, err := db.GetUnfollowList()
      Expect(err).To(BeNil())

      Expect(len(result)).To(Equal(1))
      Expect(result[0]).To(Equal(ids[0]))

      err = db.clearFollowings()
      Expect(err).To(BeNil())

      err = db.clearFollowers()
      Expect(err).To(BeNil())
    })

    It("can sync followings", func(){
      db = DatabaseConnection{}
      err := InitDatabase(&db)
      Expect(err).To(BeNil())

      r := rand.New(rand.NewSource(99))
      ids := make([]int64, 2)
      ids[0] = r.Int63()
      ids[1] = r.Int63()

      err = db.InsertFollowings(ids[1:])
      Expect(err).To(BeNil())

      err = db.SyncFollowings(ids)
      Expect(err).To(BeNil())

      err = db.clearFollowings()
      Expect(err).To(BeNil())
    })

    It("can find new no longer followings", func() {
      db = DatabaseConnection{}
      err := InitDatabase(&db)
      Expect(err).To(BeNil())

      r := rand.New(rand.NewSource(99))
      ids := make([]int64, 2)
      ids[0] = r.Int63()
      ids[1] = r.Int63()

      err = db.InsertFollowings(ids)
      Expect(err).To(BeNil())

      err = db.insertTempFollowings(ids[1:])
      Expect(err).To(BeNil())

      result, err := db.getNoLongerFollowings()
      Expect(err).To(BeNil())
      Expect(len(result)).To(Equal(1))

      err = db.Unfollow(result)
      Expect(err).To(BeNil())

      result, err = db.getNoLongerFollowings()
      Expect(err).To(BeNil())
      Expect(len(result)).To(Equal(0))

      err = db.clearTempFollowings()
      Expect(err).To(BeNil())

      err = db.clearFollowings()
      Expect(err).To(BeNil())
    })

    It("can find new following", func() {
      db = DatabaseConnection{}
      err := InitDatabase(&db)
      Expect(err).To(BeNil())

      r := rand.New(rand.NewSource(99))
      ids := make([]int64, 2)
      ids[0] = r.Int63()
      ids[1] = r.Int63()

      err = db.insertTempFollowings(ids)
      Expect(err).To(BeNil())

      err = db.InsertFollowings(ids[1:])
      Expect(err).To(BeNil())

      result, err := db.getNewFollowings()
      Expect(err).To(BeNil())

      for _, id := range result {
        log.Printf("id = %d\n", id)
      }

      Expect(len(result)).To(Equal(1))

      err = db.clearTempFollowings()
      Expect(err).To(BeNil())

      err = db.clearFollowings()
      Expect(err).To(BeNil())
    })

    It("can insert new temp followings", func() {
      db = DatabaseConnection{}
      err := InitDatabase(&db)
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
      err := InitDatabase(&db)
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
      err := InitDatabase(&db)

      Expect(err).ToNot(BeNil())
    })
  })




})
