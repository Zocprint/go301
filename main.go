package main

import (
	"code.google.com/p/gcfg"
	"flag"
	"github.com/viniciuswebdev/goahead/database"
	"log"
	"net/http"
)

type Config struct {
	General struct {
		Port string
	}
	Database database.DatabaseConf
	Table database.TableConf
}

var _cfg Config
var _db *database.Database

func handler(w http.ResponseWriter, r *http.Request) {
	hash := r.URL.Path[1:]
	log.Printf("Searching url with hash '%s' \n", hash)

	url, error := _db.FindShortenerUrlByHash(hash, &(_cfg.Table))
	if error != nil {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

func main() {
	configFilePath := flag.String("config", "./etc/goahead.ini", "Configuration file path")
	flag.Parse()

	err := gcfg.ReadFileInto(&_cfg, *configFilePath)
	if err != nil {
		panic(err.Error())
	}
	_db = database.Create(&(_cfg.Database))
	err = _db.IsValid() 
	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("/", handler)
	log.Printf("Starting Goahead on localhost:%s ...\n", _cfg.General.Port)
	err = http.ListenAndServe(":"+_cfg.General.Port, nil)
	if err != nil {
		panic(err.Error())
	}
}
