package main

// This is a fork from https://github.com/viniciuswebdev/goahead

import (
	"code.google.com/p/gcfg"
	"flag"
	"fmt"
	"github.com/Zocprint/go301/database"
	"github.com/Zocprint/go301/server"
	"log"
)

type Config struct {
	General  server.Server
	Database database.DatabaseConf
	Table    database.TableConf
	Cache    database.CacheConf
}

var _cfg Config

func main() {
	configFilePath := flag.String("config", "./etc/goahead.ini", "Configuration file path")
	help := flag.Bool("help", false, "Show me the help!")
	run := flag.Bool("run", false, "Run server")
	build := flag.Bool("build", false, "Create the scaffold")

	flag.Parse()

	if *help || (!(*run) && !(*build)) {
		showHelp()
		return
	}

	err := gcfg.ReadFileInto(&_cfg, *configFilePath)
	if err != nil {
		log.Fatal(err)
	}
	db := database.Create(&(_cfg.Database))
	err = db.IsValid()
	if err != nil {
		log.Fatal(err)
	}

	if *build {
		buildScaffold(db)
		return
	}

	if *run {
		runServer(db)
		return
	}

}

func showHelp() {
	fmt.Println("Go301 is a simple service that redirects routes\n")
	fmt.Println("Usage: ")
	fmt.Println("\t goahead [argument [-config filepath with ini configuration]]")
	fmt.Println("")
	fmt.Println("The commands are: ")
	fmt.Println("\t -help      show the help")
	fmt.Println("\t -build     create the scaffold")
	fmt.Println("\t -run       run server!")
	fmt.Println("")
	fmt.Println("Visit http://github.com/Zocprint/go301 for many informations")
}

func buildScaffold(db *database.Database) {
	fmt.Println("Creating the scaffold...")
	db.CreateTables(&(_cfg.Table))
}

func runServer(db *database.Database) {
	_cfg.General.TurnOn(db, &(_cfg.Table), &(_cfg.Cache))
}