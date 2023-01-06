package expenses

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/lib/pq"
)

type database struct {
	DB     *sql.DB
	err    error
	errMsg string
}

func (db *database) connectDatabase() {
	fmt.Println("address database server:", os.Getenv("DATABASE_URL"))
	db.DB, db.err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if db.err != nil {
		log.Fatal("Connect to database error", db.err)
	}
}

func (db *database) createDatabase() {
	createTB := `CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[] )`
	_, db.err = db.DB.Exec(createTB)
	if db.err != nil {
		db.errMsg = db.err.Error()
		log.Fatal("cant`t create table", db.err)
	}
	log.Println("Okey Database it Have Table")
}

func (db *database) insertExpenses(expenses Expenses) Expenses {
	row := db.DB.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id", expenses.Title, expenses.Amount, expenses.Note, pq.Array(&expenses.Tags))
	var resultExp Expenses
	db.err = row.Scan(&resultExp.ID, &resultExp.Title, &resultExp.Amount, &resultExp.Note, pq.Array(&resultExp.Tags))
	if db.err != nil {
		log.Fatal("cant`t insert data", db.err)
		return expenses
	}
	fmt.Println("insert todo success id : ", resultExp.ID)
	return resultExp
}

func (db *database) InitDatabase() {
	db.connectDatabase()
	db.createDatabase()
}

func (db *database) CloseDatabase() {
	db.DB.Close()
}
