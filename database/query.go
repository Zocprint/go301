package database

import (
	"database/sql"
	"github.com/pmylund/go-cache"
	"fmt"
	"time"
)

type TableConf struct {
	Name, Hash, Url string
}

type CacheConf struct {
	Time int
}

func (database *Database) FindShortenedUrlByHash(hash string, tableConf *TableConf, cacheConf *CacheConf) (string, error) {
	db, err := sql.Open(database.Driver, database.dataSource)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	c := cache.New(time.Duration(10 * int(cacheConf.Time)) , 30*time.Second)
	cachedUrl, found := c.Get(tableConf.Hash)
	if found {
        return cachedUrl.(string), nil
	}

	database_var := "?"
	if database.IsPostgres() {
		database_var = "$1"
	}
	var url string
	err = db.QueryRow(fmt.Sprintf(" SELECT %s FROM %s WHERE %s = %s ", tableConf.Url, tableConf.Name, tableConf.Hash, database_var), hash).Scan(&url)
	if err != nil {
		return url, err
	}
	if len(url) != 0 {
		c.Set(tableConf.Hash, url, 0)
	}
	return url, nil
}

func (database *Database) CreateTables(tableConf *TableConf) {
	db, err := sql.Open(database.Driver, database.dataSource)
	if err != nil {
		panic(err)
	}
	var result sql.Result

	var auto_increment_syntax = "INT AUTO INCREMENT "
	if database.IsPostgres() {
		auto_increment_syntax = "SERIAL "
	}
	fmt.Printf("creating table '%s'... \n", tableConf.Name)
	query  := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id %s PRIMARY KEY, %s varchar(255), %s varchar(255), created_at timestamp DEFAULT CURRENT_TIMESTAMP, updated_at timestamp DEFAULT CURRENT_TIMESTAMP)", tableConf.Name, auto_increment_syntax, tableConf.Hash, tableConf.Url)
	result, err =  db.Exec(query)
	if err != nil {
		panic(err)
	}
	result.RowsAffected()
	
	fmt.Printf("creating table 'goahead_statistics'... \n")
	query  = fmt.Sprintf("CREATE TABLE IF NOT EXISTS goahead_statistics (id %s PRIMARY KEY, created_at timestamp DEFAULT CURRENT_TIMESTAMP)", auto_increment_syntax)
	result, err =  db.Exec(query)
	if err != nil {
		panic(err)
	}
	result.RowsAffected()
}