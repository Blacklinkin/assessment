package expenses

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type HandlerUtil interface {
	AddExpenses(c echo.Context) error
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
		assert.Equal(t, expWant.Note, h.exp.Note)
		assert.Equal(t, expWant.Tags, h.exp.Tags)
	})
}
