package expenses

import (
<<<<<<< HEAD
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

=======
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockDB struct {
	QueryStr     string
	lastInsertID int64
	rowsAffected int64
}

func (m *MockDB) Exec(query string, args ...any) (sql.Result, error) {
	m.QueryStr = query
	return m, nil
}

func (m *MockDB) LastInsertId() (int64, error) {
	return m.lastInsertID, nil
}

func (m *MockDB) RowsAffected() (int64, error) {
	return m.rowsAffected, nil
}

func (m *MockDB) Query(query string, args ...any) (*sql.Rows, error) {
	m.QueryStr = query
	return nil, nil
}

func (m *MockDB) QueryRow(query string, args ...any) *sql.Row {
	m.QueryStr = query
	return nil
}

func (m *MockDB) Prepare(query string) (*sql.Stmt, error) {
	m.QueryStr = query
	return nil, nil
}

func (m *MockDB) Close() error {
	return nil
}

func TestCreateDatabase(t *testing.T) {
	mockDB := new(MockDB)
	db := database{DB: mockDB}

	createTB := `CREATE TABLE IF NOT EXISTS expenses ( id SERIAL PRIMARY KEY, title TEXT, amount FLOAT, note TEXT, tags TEXT[] );`

	db.createDatabase()

	assert.Nil(t, db.err)
	assert.Equal(t, createTB, mockDB.QueryStr)
>>>>>>> 01172d83f7a83fc23cde769d28c2d2d9b8859582
}
