package expenses

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type HandlerUtil interface {
	AddExpenses(c echo.Context) error
	ViewExpensesByID(c echo.Context) error
}

type MockHandler struct {
	exp           Expenses
	HandlerToCall map[string]bool
}

func (a *MockHandler) AddExpenses(c echo.Context) error {
	a.HandlerToCall["AddExpenses"] = true
	c.Response().Status = http.StatusCreated
	err := c.Bind(&a.exp)
	if err != nil {
		return err
	}
	return nil
}

func (a *MockHandler) ViewExpensesByID(c echo.Context) error {
	a.HandlerToCall["ViewExpensesByID"] = true
	c.Response().Status = http.StatusOK
	if id := c.Param("id"); id == "" {
		a.exp = Expenses{ID: 1, Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}
		return nil
	}
	a.exp = Expenses{}
	return errors.New("Id != 1")
}

func (a *MockHandler) ExpectedTocall(HandlerName string) {
	if a.HandlerToCall == nil {
		a.HandlerToCall = make(map[string]bool)
	}

	a.HandlerToCall[HandlerName] = false
}

func TestAddExpensesHandler(t *testing.T) {
	//Arrenge
	expSendJSON := `{"title": "strawberry smoothie","amount": 79,"note": "night market promotion discount 10 bath", "tags": ["food", "beverage"]}`
	expWant := Expenses{Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}

	//Act
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(expSendJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &MockHandler{}

	// Assertions
	t.Run("Should Call Handler AddExpenses NoError", func(t *testing.T) {
		h.ExpectedTocall("AddExpenses")
		assert.NoError(t, h.AddExpenses(c))
	})

	t.Run("Should Call Be Right Path /expense", func(t *testing.T) {
		assert.Equal(t, "/expenses", req.URL.Path)
	})

	t.Run("Should Call Be Right Method Post", func(t *testing.T) {
		assert.Equal(t, http.MethodPost, req.Method)
	})

	t.Run("Should Call Be Right Handler AddExpenses", func(t *testing.T) {
		assert.Equal(t, true, h.HandlerToCall["AddExpenses"])
	})

	t.Run("Should Be Response Right HTTPCode", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, c.Response().Status)
	})

	t.Run("Should Be Create Expenses Object With JSON Request", func(t *testing.T) {
		assert.Equal(t, expWant.ID, h.exp.ID)
		assert.Equal(t, expWant.Title, h.exp.Title)
		assert.Equal(t, expWant.Amount, h.exp.Amount)
		assert.Equal(t, expWant.Note, h.exp.Note)
		assert.Equal(t, expWant.Tags, h.exp.Tags)
	})
}

func TestViewExpensesByIDHandler(t *testing.T) {
	//Arrenge
	expWant := Expenses{ID: 1, Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}
	expWantJSON, _ := json.Marshal(expWant)

	strPath := fmt.Sprintf("/expenses/%d", expWant.ID)

	//Act
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, strPath, nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &MockHandler{}

	h.ExpectedTocall("ViewExpensesByID")
	err := h.ViewExpensesByID(c)
	expJSON, _ := json.Marshal(h.exp)

	// Assertions
	t.Run("Should Call Handler ViewExpensesByID NoError", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("Should Call Be Right Path /expense/1", func(t *testing.T) {
		assert.Equal(t, strPath, req.URL.Path)
	})

	t.Run("Should Call Be Right Method Get", func(t *testing.T) {
		assert.Equal(t, http.MethodGet, req.Method)
	})

	t.Run("Should Call Be Right Handler ViewExpensesByID", func(t *testing.T) {
		assert.Equal(t, true, h.HandlerToCall["ViewExpensesByID"])
	})

	t.Run("Should Be Response Right HTTPCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, c.Response().Status)
	})

	t.Run("Should Be Create Expenses Object With Get /expenses/1 Request", func(t *testing.T) {
		assert.Equal(t, expWant.ID, h.exp.ID)
		assert.Equal(t, expWant.Title, h.exp.Title)
		assert.Equal(t, expWant.Amount, h.exp.Amount)
		assert.Equal(t, expWant.Note, h.exp.Note)
		assert.Equal(t, expWant.Tags, h.exp.Tags)
	})

	t.Run("Should Be Response JSON String With Get/exoenses/1 Request", func(t *testing.T) {
		assert.Equal(t, string(expWantJSON), string(expJSON))
	})
}
