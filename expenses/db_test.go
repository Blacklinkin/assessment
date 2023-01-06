package expenses

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestCreateDatabase(t *testing.T) {
	//Arrenge
	db, mock, _ := sqlmock.New()
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS expenses").WillReturnResult(sqlmock.NewResult(0, 0))
	dbt := database{DB: db}

	//Act
	dbt.createDatabase()

	//Assert
	assert.Nil(t, dbt.err)
}

func TestInsertDatabese(t *testing.T) {
	////Arrenge
	exp := Expenses{Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}
	db, mock, _ := sqlmock.New()
	row := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).AddRow(1, exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags))
	mock.ExpectQuery("INSERT INTO expenses").WithArgs(exp.Title, exp.Amount, exp.Note, pq.Array(&exp.Tags)).WillReturnRows(row)
	dbt := database{DB: db}

	//Act
	result := dbt.insertExpenses(exp)

	//Assert
	assert.Nil(t, dbt.err)
	assert.NotEqual(t, exp.ID, result.ID)
	assert.Equal(t, exp.Title, result.Title)
	assert.Equal(t, exp.Amount, result.Amount)
	assert.Equal(t, exp.Note, result.Note)
	assert.Equal(t, exp.Tags, result.Tags)
}
