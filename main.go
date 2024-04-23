package main

import (
	"github.com/ag89201/go_final_project/app/domain"
	"github.com/ag89201/go_final_project/app/model"
	"github.com/ag89201/go_final_project/app/server"

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
	model.Database, err = model.NewDataBase(dbFile)
	if err != nil {
		log.Panic(err)
	}
	defer model.Database.Close()
	// Create the tables
	log.Info("open|create table......")
	model.Database.CreateSchedulerTable()

	log.Info("create index......")
	model.Database.CreateIndex()

	// Start the web server
	port := domain.GetEnv("TODO_PORT", defPort)
	log.Fatal(server.Start(port, webDir))
}
