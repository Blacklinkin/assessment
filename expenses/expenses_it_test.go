package expenses

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const serverPort = 80

func TestITAddExpenses(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		h := Handler{}
		h.InitialDB()

		e.POST("/expenses", h.AddExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
		h.CloseDB()
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	expSendJSON := `{"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food", "beverage"]}`
	expWant := Expenses{Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}
	expWantJSON, _ := json.Marshal(expWant)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(expSendJSON))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	resultStruc := Expenses{}
	json.Unmarshal(byteBody, &resultStruc)

	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEqual(t, expWantJSON, byteBody)
		assert.Equal(t, expWant.Title, resultStruc.Title)
		assert.Equal(t, expWant.Amount, resultStruc.Amount)
		assert.Equal(t, expWant.Note, resultStruc.Note)
		assert.Equal(t, expWant.Tags, resultStruc.Tags)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITViewExpensesByID(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		h := Handler{}
		h.InitialDB()

		e.POST("/expenses", h.AddExpenses)
		e.GET("/expenses/:id", h.ViewExpensesByID)
		e.Start(fmt.Sprintf(":%d", serverPort))
		h.CloseDB()
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	// Arrange
	expSendJSON := `{"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food", "beverage"]}`
	expWant := Expenses{Title: "strawberry smoothie", Amount: 79, Note: "night market promotion discount 10 bath", Tags: []string{"food", "beverage"}}
	expWantJSON, _ := json.Marshal(expWant)
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(expSendJSON))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	resultStruc := Expenses{}
	json.Unmarshal(byteBody, &resultStruc)

	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEqual(t, expWantJSON, byteBody)
		assert.Equal(t, expWant.Title, resultStruc.Title)
		assert.Equal(t, expWant.Amount, resultStruc.Amount)
		assert.Equal(t, expWant.Note, resultStruc.Note)
		assert.Equal(t, expWant.Tags, resultStruc.Tags)
	}

	// Arrange
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, resultStruc.ID), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client = http.Client{}

	// Act
	resp, err = client.Do(req)
	assert.NoError(t, err)

	byteBodyGetID, _ := ioutil.ReadAll(resp.Body)

	resultStrucGetID := Expenses{}
	json.Unmarshal(byteBodyGetID, &resultStrucGetID)

	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, byteBody, byteBodyGetID)
		assert.Equal(t, resultStruc, resultStrucGetID)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}

func TestITUpdateExpenses(t *testing.T) {
	// Setup server
	eh := echo.New()
	go func(e *echo.Echo) {
		h := Handler{}
		h.InitialDB()

		e.PUT("/expenses/:id", h.UpdateExpenses)
		e.Start(fmt.Sprintf(":%d", serverPort))
		h.CloseDB()
	}(eh)
	for {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", serverPort), 30*time.Second)
		if err != nil {
			log.Println(err)
		}
		if conn != nil {
			conn.Close()
			break
		}
	}

	//Arrange
	reqBody := `{"title":"apple smoothie","amount":89,"note":"no discount","tags": ["beverage"]}`
	//strucReqBody := Expenses{Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}}
	strucResBody := Expenses{ID: 1, Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}}
	//expSendJSON, _ := json.Marshal(strucReqBody)
	expWantJSON, _ := json.Marshal(strucResBody)

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/:id", serverPort), strings.NewReader(reqBody))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, expWantJSON, strings.TrimSpace(string(byteBody)))
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = eh.Shutdown(ctx)
	assert.NoError(t, err)
}
