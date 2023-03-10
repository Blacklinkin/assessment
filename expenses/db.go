package expenses

import (
	"database/sql"
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
	log.Println("address database server:", os.Getenv("DATABASE_URL"))
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
	row := db.DB.QueryRow("INSERT INTO expenses (title, amount, note, tags) values ($1, $2, $3, $4) RETURNING id,title, amount, note, tags", expenses.Title, expenses.Amount, expenses.Note, pq.Array(&expenses.Tags))
	resultExp := Expenses{}
	db.err = row.Scan(&resultExp.ID, &resultExp.Title, &resultExp.Amount, &resultExp.Note, pq.Array(&resultExp.Tags))
	if db.err != nil {
		log.Fatal("cant`t insert data", db.err)
		return expenses
	}
	log.Println("insert todo success id : ", resultExp.ID)
	return resultExp
}

func (db *database) viewExpensesDataByID(id int) Expenses {
	stmt, err := db.DB.Prepare("SELECT id, title, amount, note, tags FROM expenses where id=$1")

	db.err = err
	if db.err != nil {
		db.errMsg = "cant`t prepare statement"
		log.Fatal(db.errMsg)
		return Expenses{}
	}

	row := stmt.QueryRow(id)
	expQ := Expenses{}
	db.err = row.Scan(&expQ.ID, &expQ.Title, &expQ.Amount, &expQ.Note, pq.Array(&expQ.Tags))
	if db.err != nil {
		db.errMsg = "cant`t return statement"
		log.Fatal(db.errMsg)
		return Expenses{}
	}
	log.Println("view by id todo success: ", expQ)
	return expQ
}

func (db *database) updateExpensesDataBase(id int, expUpdate Expenses) {
	stmt, err := db.DB.Prepare("UPDATE expenses SET title=$2, amount=$3, note=$4, tags=$5 WHERE id=$1")

	db.err = err
	if db.err != nil {
		db.errMsg = "cant`t prepare statement update"
		log.Fatal(db.errMsg)
	}

	_, err = stmt.Exec(id, expUpdate.Title, expUpdate.Amount, expUpdate.Note, pq.Array(&expUpdate.Tags))

	db.err = err
	if db.err != nil {
		db.errMsg = "cant`t return result update"
		log.Fatal(db.errMsg)
	}

	log.Println("update todo success id : ", id)
}

func (db *database) viewAllExpenses() []Expenses {
	rows, err := db.DB.Query("SELECT id, title, amount, note, tags FROM expenses")
	db.err = err
	if db.err != nil {
		db.errMsg = "Error querying expenses"
		log.Fatal(db.errMsg)
		return []Expenses{}
	}
	defer rows.Close()

	expSet := []Expenses{}
	for rows.Next() {
		var exp Expenses
		err := rows.Scan(&exp.ID, &exp.Title, &exp.Amount, &exp.Note, pq.Array(&exp.Tags))
		db.err = err
		if err != nil {
			db.errMsg = "Error scanning row"
			log.Fatal(db.errMsg)
			return []Expenses{}
		}
		log.Println("add expenses todo success: ", exp)
		expSet = append(expSet, exp)
	}
	log.Println("view all todo success: ", expSet)
	return expSet
}

func (db *database) InitDatabase() {
	db.connectDatabase()
	db.createDatabase()
}

func (db *database) CloseDatabase() {
	db.DB.Close()
}
