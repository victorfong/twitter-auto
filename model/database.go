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
  // InsertFollower(follower Follower) error

  SyncFollowers(ids []int64) error
  SyncFollowings(id []int64) error
}

type DatabaseConnection struct{
  db *sql.DB
  dbMap *gorp.DbMap
}

func InitDatabase(d *DatabaseConnection) error{
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

func (d DatabaseConnection) SyncFollowers(ids []int64) error {
  tx, err := d.db.Begin()
  if err != nil {
    return err
  }

  err = d.clearFollowers()
  if err != nil {
    return err
  }

  err = d.insertFollowers(ids)
  if err != nil {
    return err
  }

  err = tx.Commit()
  if err != nil {
    return err
  }

  return nil
}

func (d DatabaseConnection) SyncFollowings(ids []int64) error {
  err := d.clearTempFollowings()
  if err != nil {
    return err
  }

  err = d.insertTempFollowings(ids)
  if err != nil {
    return err
  }

  newFollowings, err := d.getNewFollowings()
  if err != nil {
    return err
  }

  noLongerFollowings, err := d.getNoLongerFollowings()
  // _, err = d.getNoLongerFollowings()
  if err != nil {
    return err
  }

  err = d.unfollow(noLongerFollowings)
  if err != nil {
    return err
  }

  err = d.insertFollowings(newFollowings)
  if err != nil {
    return err
  }

  err = d.clearTempFollowings()
  return err
}

func (d DatabaseConnection) unfollow(ids []int64) error {

  if len(ids) > 0 {
    sqlStr := "UPDATE following SET unfollowed = true WHERE twitter_id IN ("
    vals := []interface{}{}

    for _, id := range ids {
      sqlStr += "?, "
      vals = append(vals, id)
    }

    sqlStr = sqlStr[0:len(sqlStr)-2]
    sqlStr += ")"

    log.Printf("sqlStr = %s", sqlStr)

    statement, err := d.db.Prepare(sqlStr)
    if err != nil {
      return err
    }

    _, err = statement.Exec(vals...)
    if err != nil {
      return err
    }
  }

  return nil
}

func (d DatabaseConnection) getNoLongerFollowings() ([]int64, error) {
  sqlStr := "SELECT twitter_id FROM following WHERE twitter_id NOT IN (SELECT twitter_id from temp_following) AND unfollowed = false ORDER BY since"

  log.Printf("sqlStr = %s", sqlStr)

  statement, err := d.db.Prepare(sqlStr)
  if err != nil {
    return nil, err
  }

  rows, err := statement.Query()
  if err != nil {
    return nil, err
  }

  result := []int64{}
  for rows.Next() {
    var id int64
    err = rows.Scan(&id)

    if err != nil {
      return nil, err
    }

    result = append(result, id)
  }

  return result, nil
}

func (d DatabaseConnection) getNewFollowings() ([]int64, error) {

  sqlStr := "SELECT twitter_id FROM temp_following WHERE twitter_id NOT IN (SELECT twitter_id from following)"

  log.Printf("sqlStr = %s", sqlStr)

  statement, err := d.db.Prepare(sqlStr)
  if err != nil {
    return nil, err
  }

  rows, err := statement.Query()
  if err != nil {
    return nil, err
  }

  result := []int64{}
  for rows.Next() {
    var id int64
    err = rows.Scan(&id)

    if err != nil {
      return nil, err
    }

    result = append(result, id)
  }

  return result, nil
}

func (d DatabaseConnection) clearFollowers() error {
  sqlStr := "DELETE FROM follower"

  log.Printf("sqlStr = %s", sqlStr)

  statement, err := d.db.Prepare(sqlStr)
  if err != nil {
    return err
  }

  _, err = statement.Exec()
  if err != nil {
    return err
  }

  return nil
}

func (d DatabaseConnection) clearTempFollowings() error {
  sqlStr := "DELETE FROM temp_following"

  log.Printf("sqlStr = %s", sqlStr)

  statement, err := d.db.Prepare(sqlStr)
  if err != nil {
    return err
  }

  _, err = statement.Exec()
  if err != nil {
    return err
  }

  return nil
}


func (d DatabaseConnection) clearFollowings() error {
  sqlStr := "DELETE FROM following"

  log.Printf("sqlStr = %s", sqlStr)

  statement, err := d.db.Prepare(sqlStr)
  if err != nil {
    return err
  }

  _, err = statement.Exec()
  if err != nil {
    return err
  }

  return nil
}

func (d DatabaseConnection) insertFollowings(ids []int64) error {
  if len(ids) > 0 {
    sqlStr := "INSERT INTO following(twitter_id) VALUES "
    vals := []interface{}{}

    for _, id := range ids {
      sqlStr += "(?),"
      vals = append(vals, id)
    }

    sqlStr = sqlStr[0:len(sqlStr)-1]
    log.Printf("sqlStr = %s", sqlStr)

    statement, err := d.db.Prepare(sqlStr)
    if err != nil {
      return err
    }

    _, err = statement.Exec(vals...)
    if err != nil {
      return err
    }
  }

  return nil
}

func (d DatabaseConnection) insertTempFollowings(ids []int64) error {
  if len(ids) > 0 {
    sqlStr := "INSERT INTO temp_following(twitter_id) VALUES "
    vals := []interface{}{}

    for _, id := range ids {
      sqlStr += "(?),"
      vals = append(vals, id)
    }

    sqlStr = sqlStr[0:len(sqlStr)-1]
    log.Printf("sqlStr = %s", sqlStr)

    statement, err := d.db.Prepare(sqlStr)
    if err != nil {
      return err
    }

    _, err = statement.Exec(vals...)
    if err != nil {
      return err
    }
  }
  return nil
}

func (d DatabaseConnection) insertFollowers(ids []int64) error {
  sqlStr := "INSERT INTO follower(twitter_id) VALUES "
  vals := []interface{}{}

  for _, id := range ids {
    sqlStr += "(?),"
    vals = append(vals, id)
  }

  sqlStr = sqlStr[0:len(sqlStr)-1]
  log.Printf("sqlStr = %s", sqlStr)

  statement, err := d.db.Prepare(sqlStr)
  if err != nil {
    return err
  }

  _, err = statement.Exec(vals...)
  if err != nil {
    return err
  }

  return nil
}



func (d DatabaseConnection) InsertFollower(follower *Follower) error {
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
	d.dbMap.AddTableWithName(Following{}, "following")

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

func (d DatabaseConnection) Exec(strStatement string) (sql.Result, error){
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
