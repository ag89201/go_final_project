package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	_ "modernc.org/sqlite"
)

const (
	defPort   = 7540
	webDir    = "./web"
	defDbName = "scheduler.db"
)

type Task struct {
	ID      int64  `db:"id"`
	Date    string `db:"date"`
	Title   string `db:"title"`
	Comment string `db:"comment"`
	Repeat  string `db:"repeat"`
}

func getPort() int {

	envPort := os.Getenv("TODO_PORT")
	if len(envPort) > 0 {
		if eport, err := strconv.ParseInt(envPort, 10, 32); err == nil {
			return int(eport)
		}
	}
	return defPort // Default port

}

func ExistDbFile(filePath string) error {

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("Creating database file: %s\n", filePath)
		_, err := os.Create(filePath)
		if err != nil {

			return err
		}
	}

	return nil
}

func getDbFilePath() (string, error) {
	var filePath string
	envFilePath := os.Getenv("TODO_DBFILE")
	if len(envFilePath) > 0 {
		return envFilePath, nil
	}

	appPath, err := os.Executable()
	if err != nil {
		return filePath, err
	}

	filePath = filepath.Join(filepath.Dir(appPath), defDbName)

	return filePath, nil

}

func main() {

	//get db file
	dbFilePath, err := getDbFilePath()
	if err != nil {
		log.Fatal(err)
		return
	}
	err = ExistDbFile(dbFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	//open database
	db, err := sql.Open("sqlite", dbFilePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	//create table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS scheduler (
		id INTEGER PRIMARY KEY,
		date TEXT,
		title TEXT,
		comment TEXT,
		repeat TEXT
		)`)
	if err != nil {
		log.Fatal(err)
		return
	}
	//create index
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date)`)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer db.Close()

	// Start the web server
	portInt := getPort()

	http.Handle("/", http.FileServer(http.Dir(webDir)))

	fmt.Printf("Server started on %d\n port", portInt)
	err = http.ListenAndServe(fmt.Sprintf(":%d", portInt), nil)
	if err != nil {
		log.Fatal(err)
		return
	}

}
