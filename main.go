package main

import (
	"fmt"

	"github.com/ag89201/go_final_project/app/db"
	"github.com/ag89201/go_final_project/app/domain"
	"github.com/ag89201/go_final_project/app/http"

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
	if err != nil {
		panic(err)
	}
	if domain.FileNotExists(dbFile) {
		err := domain.CreateFile(dbFile)
		if err != nil {
			panic(err)
		}
	}
	//open database
	db.Database, err = db.New(dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Database.Close()
	// Create the tables
	db.Database.CreateSchedulerTable()
	db.Database.CreateIndex()

	// Start the web server
	port := domain.GetEnv("TODO_PORT", defPort)
	err = http.StartServer(port, webDir)
	if err != nil {
		fmt.Println(err)
	}
}
