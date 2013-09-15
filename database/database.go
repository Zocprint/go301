package database

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DatabaseConf struct {
    User     string
    Password string
    Database string
}

type Database struct {
	User, Password, Database string
}

func (database *Database) FindShortenerUrlByHash(hash string) string {
	db, err := sql.Open("mysql", database.User+":"+database.Password+"@/"+database.Database)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var url string
	err = db.QueryRow("SELECT url FROM shortened_url WHERE hash = '" + hash + "'").Scan(&url)
	if err != nil {
		panic(err.Error())
	}
	return url
}

func Create(conf *DatabaseConf) *Database {
	db := new(Database)
	db.User = conf.User
    db.Password = conf.Password
    db.Database = conf.Database

    return db
}
