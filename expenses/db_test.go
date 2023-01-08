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

func TestViewDataByID(t *testing.T) {
	////Arrenge
	idParam := 1
	expWant := Expenses{ID: 1, Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}
	db, mock, _ := sqlmock.New()
	row := sqlmock.NewRows([]string{"id", "title", "amount", "note", "tags"}).
		AddRow(expWant.ID, expWant.Title, expWant.Amount, expWant.Note, pq.Array(&expWant.Tags))
	mock.ExpectPrepare("SELECT id, title, amount, note, tags FROM expenses").ExpectQuery().
		WithArgs(idParam).WillReturnRows(row)
	dbt := database{DB: db}

	//Act
	result := dbt.viewExpensesDataByID(idParam)

	//Assert
	assert.Nil(t, dbt.err)
	assert.Equal(t, expWant.ID, result.ID)
	assert.Equal(t, expWant.Title, result.Title)
	assert.Equal(t, expWant.Amount, result.Amount)
	assert.Equal(t, expWant.Note, result.Note)
	assert.Equal(t, expWant.Tags, result.Tags)
}

func TestUpdateDataBase(t *testing.T) {
	////Arrenge
	idParam := 1
	dataUpdate := Expenses{Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}}
	dataUpdated := Expenses{ID: 1, Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}}
	db, mock, _ := sqlmock.New()
	resultExec := sqlmock.NewResult(1, 1)
	mock.ExpectPrepare("UPDATE expenses").ExpectExec().WithArgs(idParam, dataUpdate.Title, dataUpdate.Amount, dataUpdate.Note, pq.Array(&dataUpdate.Tags)).WillReturnResult(resultExec)
	dbt := database{DB: db}

	//Act
	resultID := dbt.updateExpensesDataBase(idParam, dataUpdate)

	//Assert
	assert.Nil(t, dbt.err)
	assert.Equal(t, dataUpdated.ID, resultID)
}
