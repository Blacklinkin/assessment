package expenses

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type database struct {
	DB  *sql.DB
	err error
}

func (db *database) connectDatabase() {
	url := os.Getenv("DATABASE_URL")
	fmt.Println("address database server:", url)
	db.DB, db.err = sql.Open("postgres", url)
	if db.err != nil {
		log.Fatal("Connect to database error", db.err)
	}
}

func (db *database) createDatabase() {
	createTB := `CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[] );`
	_, db.err = db.DB.Exec(createTB)
	if db.err != nil {
		log.Fatal("cant`t create table", db.err)
	}

	log.Println("Okey Database it Have Table")
}

func (db *database) InitDatabase() {
	db.connectDatabase()
	db.createDatabase()
}

func (db *database) CloseDatabase() {
	db.DB.Close()
}
