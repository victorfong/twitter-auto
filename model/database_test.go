package model
import (
  . "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
  "os"
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
    PIt("creates a database connection", func(){
      uri := "mysql2://twitter-auto2:twitter-auto2@localhost:3306/twitter-auto2?reconnect=true"
      db, err := initConnection(uri)
      Expect(err).To(BeNil())
      Expect(db).ToNot(BeNil())
    })
  })

  Context("when database connection is invalid", func(){
    It("returns an error", func(){
      uri := "mysql://whatever:lalala@localhost:3306/whatever?reconnect=true"
      _, err := initConnection(uri)
      Expect(err).ToNot(BeNil())
    })
  })


})
