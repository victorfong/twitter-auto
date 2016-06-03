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
  "github.com/coopernurse/gorp"
  "github.com/rubenv/sql-migrate"
)

type MySql struct {
  Cleardb []map[string]map[string]string
}

type Database interface{
  Init()
  Exec(strStatement string) (sql.Result, error)
  InsertFollower(follower Follower) error
}

type DatabaseConnection struct{
  db *sql.DB
  dbMap *gorp.DbMap
}

func (d *DatabaseConnection) Init() error{
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

  // InitDb()
  err = initDbMap(d)
  if err != nil {
    return err
  }

  log.Println("Finished Initializing Database")
  return nil
}

func (d *DatabaseConnection) InsertFollower(follower *Follower) error {
  follower.TwitterId = 1234
  log.Printf("follower.TwitterId = %d", follower.TwitterId)

  return d.dbMap.Insert(follower)
}

func initDbMap(d *DatabaseConnection) error {
  migrationDir := "../db/migrations"
  mysqldbmap := &gorp.DbMap{
    Db: d.db,
    Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"},
  }

	migrations := &migrate.FileMigrationSource{
		Dir: migrationDir,
	}
	n, err := migrate.Exec(d.db, "mysql", migrations, migrate.Up)

	if n > 0 {
		log.Printf("Successfully ran %v migrations\n", n)
	}

  if err != nil {
		return err
	}

  d.dbMap = mysqldbmap
  d.dbMap.AddTableWithName(Follower{}, "follower")
  // .SetKeys(true, "TwitterId")
	d.dbMap.AddTableWithName(Following{}, "following").SetKeys(true, "Id")

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

func (d *DatabaseConnection) Exec(strStatement string) (sql.Result, error){
  statement, err := d.db.Prepare(strStatement)
  if err != nil {
    return nil, err
  }

  result, err := statement.Exec()
  if err != nil {
    return nil, err
  }
  return result, nil
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
