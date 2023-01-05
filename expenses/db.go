package expenses

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type DataBaseUtil interface {
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	Close() error
}

type database struct {
	DB  DataBaseUtil
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
