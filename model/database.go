package model

import (
  "fmt"
  "os"
  "log"
  "errors"
  "encoding/json"
  "net/url"

  _ "github.com/go-sql-driver/mysql"
  "database/sql"
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
  db *sql.DB
}

func (d DatabaseConnection) Init(uri string) error{
  log.Println("Initializing Database")

  uri, err := getURI()
  if err != nil {
    return err
  }

  db, err := initConnection(uri)
  if err != nil {
    return err
  }
  d.db = db

  log.Println("Finished Initializing Database")
  return nil
}

func formattedUrl(url *url.URL) string {
	return fmt.Sprintf(
		"%v@tcp(%v)%v?parseTime=true",
		url.User,
		url.Host,
		url.Path,
	)
}

func initConnection(uri string) (*sql.DB, error) {

  url, err := url.Parse(uri)
	if err != nil {
		log.Fatalln("Error parsing DATABASE_URL", err)
	}

  databaseUrl := formattedUrl(url)

  db, err := sql.Open("mysql", databaseUrl)
  if err != nil {
    return nil, err
  }

  statement, err := db.Prepare("SELECT 1 FROM testtable LIMIT 1;")
  if err != nil {
    return nil, err
  }

  _, err = statement.Exec()
  if err != nil {
    return nil, err
  }
  //
  // log.Printf("res = %v", res)
  return db, nil
}

// func isDatabaseInitialized() (bool, error) {
//   statement, err := db.Prepare("SELECT 1 FROM testtable LIMIT 1;")
//   if err != nil {
//     return err
//   }
// }

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
