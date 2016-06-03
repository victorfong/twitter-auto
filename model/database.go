package model

import (
  "os"
  "log"
  "errors"
  "encoding/json"

  // _ "github.com/go-sql-driver/mysql"
  // "database/sql"
)

type MySql struct {
  Cleardb []map[string]map[string]string
}

type Credentials struct {
  uri string
}

type Database interface{
  Init(uri string)
}

type DatabaseConnection struct{

}

func (d DatabaseConnection) Init(uri string) error{
  log.Println("Initializing Database")

  uri, err := getURI()
  if err != nil {
    return err
  }

  log.Println("Finished Initializing Database")
  return nil
}

func getURI() (string, error) {
	envVar := os.Getenv("VCAP_SERVICES")

  if envVar == "" {
    return "", errors.New("MySql Service not set in env path")
  }

  // var vcapServices map[string]interface{}
  var vcapServices MySql
  err := json.Unmarshal([]byte(envVar), &vcapServices)
  if err != nil {
    return "", errors.New("MySql Service cannot be found")
  }

  if len(vcapServices.Cleardb) == 0{
    return "", errors.New("MySql Service cannot be found")
  }

  return vcapServices.Cleardb[0]["credentials"]["uri"], nil
}
