//go:build unit
// +build unit

package expenses

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type HandlerUtil interface {
	AddExpenses(c echo.Context) error
	ViewExpensesByID(c echo.Context) error
	UpdateExpensesHandler(c echo.Context) error
	ViewAllExpensesHandler(c echo.Context) error
}

type MockHandler struct {
	exp           Expenses
	expReqGot     Expenses
	expSetGot     []Expenses
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
	if id := c.Param("id"); id != "" {
		Id, _ := strconv.Atoi(id)
		a.exp = Expenses{ID: Id, Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}
		return nil
	}
	return echo.ErrBadRequest
}

func (a *MockHandler) UpdateExpensesHandler(c echo.Context) error {
	a.HandlerToCall["UpdateExpenses"] = true
	err := c.Bind(&a.exp)
	if err != nil {
		return err
	}
	c.Response().Status = http.StatusCreated
	if id := c.Param("id"); id != "" {
		a.expReqGot = a.exp
		Id, _ := strconv.Atoi(id)
		a.exp = Expenses{ID: Id, Title: a.exp.Title, Amount: a.exp.Amount, Note: a.exp.Note, Tags: a.exp.Tags}
		return nil
	}
	return echo.ErrBadRequest
}

func (a *MockHandler) ViewAllExpensesHandler(c echo.Context) error {
	a.HandlerToCall["ViewAllExpenses"] = true
	c.Response().Status = http.StatusOK
	a.expSetGot = append(a.expSetGot,
		Expenses{ID: 1, Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}},
		Expenses{ID: 2, Title: "iPhone 14 Pro Max 1TB", Amount: 66900, Note: "birthday gift from my love", Tags: []string{"gadget"}},
	)
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
	expSendJSON := `{"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food", "beverage"]}`
	expWant := Expenses{Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}

	//Act
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(expSendJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &MockHandler{}

	h.ExpectedTocall("AddExpenses")
	err := h.AddExpenses(c)

	// Assertions
	t.Run("Should Call Handler AddExpenses NoError", func(t *testing.T) {
		assert.NoError(t, err)
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

	t.Run("Should Be Response Right HTTPCode On AddExpenses", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, c.Response().Status)
	})

	t.Run("Should Be Create Expenses Object With JSON Request On AddExpenses", func(t *testing.T) {
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

	//Act
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &MockHandler{}

	h.ExpectedTocall("ViewExpensesByID")
	err := h.ViewExpensesByID(c)
	expJSON, _ := json.Marshal(h.exp)

	// Assertions
	t.Run("Should Call Handler ViewExpensesByID NoError", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("Should Call Be Right Path /expense/1", func(t *testing.T) {
		assert.Equal(t, "/expenses/:id", req.URL.Path)
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

	t.Run("Should Be Response JSON String With Get /expenses/1 Request", func(t *testing.T) {
		assert.Equal(t, string(expWantJSON), string(expJSON))
	})
}

func TestUpdateExpensesHandler(t *testing.T) {
	//Arrenge
	reqBody := `{"title":"apple smoothie","amount":89,"note":"no discount","tags": ["beverage"]}`
	strucReqBody := Expenses{Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}}
	resBody := `{"id":1,"title":"apple smoothie","amount":89,"note":"no discount","tags":["beverage"]}`
	strucResBody := Expenses{ID: 1, Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}}

	//Act
	e := echo.New()
	req := httptest.NewRequest(http.MethodPut, "/expenses/:id", strings.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("1")
	h := &MockHandler{}

	h.ExpectedTocall("UpdateExpenses")
	err := h.UpdateExpensesHandler(c)
	expJSON, _ := json.Marshal(h.exp)

	// Assertions
	t.Run("Should Call Handler UpdateExpenses NoError", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("Should Call Be Right Path /expense/1", func(t *testing.T) {
		assert.Equal(t, "/expenses/:id", req.URL.Path)
	})

	t.Run("Should Call Be Right Method Put", func(t *testing.T) {
		assert.Equal(t, http.MethodPut, req.Method)
	})

	t.Run("Should Call Be Right Handler UpdateExpenses", func(t *testing.T) {
		assert.Equal(t, true, h.HandlerToCall["UpdateExpenses"])
	})

	t.Run("Should Be Response Right HTTPCode", func(t *testing.T) {
		assert.Equal(t, http.StatusCreated, c.Response().Status)
	})

	t.Run("Should be Receive Expenses Struct From Request", func(t *testing.T) {
		assert.Equal(t, strucReqBody.ID, h.expReqGot.ID)
		assert.Equal(t, strucReqBody.Title, h.expReqGot.Title)
		assert.Equal(t, strucReqBody.Amount, h.expReqGot.Amount)
		assert.Equal(t, strucReqBody.Note, h.expReqGot.Note)
		assert.Equal(t, strucReqBody.Tags, h.expReqGot.Tags)
	})

	t.Run("Should Be Update Expenses Object With Put /expenses/1 Request", func(t *testing.T) {
		assert.Equal(t, strucResBody.ID, h.exp.ID)
		assert.Equal(t, strucResBody.Title, h.exp.Title)
		assert.Equal(t, strucResBody.Amount, h.exp.Amount)
		assert.Equal(t, strucResBody.Note, h.exp.Note)
		assert.Equal(t, strucResBody.Tags, h.exp.Tags)
	})

	t.Run("Should Be Response JSON String With Put /exoenses/1 ", func(t *testing.T) {
		assert.Equal(t, string(resBody), string(expJSON))
	})
}

func TestViewAllExpensesHandler(t *testing.T) {
	//Arrenge
	expStrucWant := []Expenses{
		{ID: 1, Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}},
		{ID: 2, Title: "iPhone 14 Pro Max 1TB", Amount: 66900, Note: "birthday gift from my love", Tags: []string{"gadget"}},
	}
	expJSONSetWant, _ := json.Marshal(expStrucWant)

	//Act
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := &MockHandler{}

	h.ExpectedTocall("ViewAllExpenses")
	err := h.ViewAllExpensesHandler(c)
	expJSONSetGot, _ := json.Marshal(h.expSetGot)

	// Assertions
	t.Run("Should Call Handler ViewAllExpenses NoError", func(t *testing.T) {
		assert.NoError(t, err)
	})

	t.Run("Should Call Be Right Path /expense", func(t *testing.T) {
		assert.Equal(t, "/expenses", req.URL.Path)
	})

	t.Run("Should Call Be Right Method Get", func(t *testing.T) {
		assert.Equal(t, http.MethodGet, req.Method)
	})

	t.Run("Should Call Be Right Handler ViewAllExpenses", func(t *testing.T) {
		assert.Equal(t, true, h.HandlerToCall["ViewAllExpenses"])
	})

	t.Run("Should Be Response Right HTTPCode", func(t *testing.T) {
		assert.Equal(t, http.StatusOK, c.Response().Status)
	})

	t.Run("Should Be Create Expenses Object With Get /expenses Request", func(t *testing.T) {
		assert.Equal(t, expStrucWant, h.expSetGot)
	})

	t.Run("Should Be Response JSON String With Get /expenses Request", func(t *testing.T) {
		assert.Equal(t, string(expJSONSetWant), string(expJSONSetGot))
	})
}
