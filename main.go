package main

import (
	"github.com/ag89201/go_final_project/app/db"
	"github.com/ag89201/go_final_project/app/domain"
	"github.com/ag89201/go_final_project/app/http"

	log "github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

const (
	DefDbName = "scheduler.db"
	defPort   = "7540"
	webDir    = "./web"
)

func main() {
	// get db filename
	dbFile, err := domain.GetFileName("TODO_DBFILE", DefDbName)
	log.Info("open database: " + dbFile)
	if err != nil {
		log.Panic(err)
	}
	if domain.FileNotExists(dbFile) {
		log.Info("file not exists......")
		err := domain.CreateFile(dbFile)
		if err != nil {
			log.Panic(err)

		}
		log.Info("creating new file: " + dbFile)
	}
	//open database
	db.Database, err = db.New(dbFile)
	if err != nil {
		log.Panic(err)
	}
	defer db.Database.Close()
	// Create the tables
	log.Info("open|create table......")
	db.Database.CreateSchedulerTable()

	log.Info("create index......")
	db.Database.CreateIndex()

	// Start the web server
	port := domain.GetEnv("TODO_PORT", defPort)
	log.Fatal(http.StartServer(port, webDir))
}
