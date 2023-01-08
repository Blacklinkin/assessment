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

func initialServer() (eh *echo.Echo) {
	eh = echo.New()
	go func(e *echo.Echo) {
		h := Handler{}
		h.InitialDB()

		e.POST("/expenses", h.AddExpenses)
		e.GET("/expenses/:id", h.ViewExpensesByID)
		e.PUT("/expenses/:id", h.UpdateExpenses)
		e.GET("/expenses", h.ViewAllExpenses)

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
	return
}

func shutdownServer(eh *echo.Echo, t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errShutDown := eh.Shutdown(ctx)
	assert.NoError(t, errShutDown)
}

func TestITAddExpenses(t *testing.T) {
	// Setup server
	eh := initialServer()

	// Arrange:SetData
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

	byteBodyGot, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	resultStructGot := Expenses{}
	json.Unmarshal(byteBodyGot, &resultStructGot)

	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		assert.NotEqual(t, expWantJSON, byteBodyGot)
		assert.Equal(t, expWant.Title, resultStructGot.Title)
		assert.Equal(t, expWant.Amount, resultStructGot.Amount)
		assert.Equal(t, expWant.Note, resultStructGot.Note)
		assert.Equal(t, expWant.Tags, resultStructGot.Tags)
	}

	//ShutdownServer
	shutdownServer(eh, t)
}

func TestITViewExpensesByID(t *testing.T) {
	// Setup server
	eh := initialServer()

	// Arrange:SetData
	expSendJSON := `{"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food", "beverage"]}`
	// Arrange:SetRequest
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(expSendJSON))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act:Response
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBodyGot, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	resultStructGot := Expenses{}
	json.Unmarshal(byteBodyGot, &resultStructGot)

	resp.Body.Close()

	// Arrange:RequestByID
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, resultStructGot.ID), nil)
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client = http.Client{}

	// Act:ResponseByID
	resp, err = client.Do(req)
	assert.NoError(t, err)

	byteBodyGotByID, _ := ioutil.ReadAll(resp.Body)

	resultStructGotByID := Expenses{}
	json.Unmarshal(byteBodyGotByID, &resultStructGotByID)

	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, byteBodyGot, byteBodyGotByID)
		assert.Equal(t, resultStructGot, resultStructGotByID)
	}

	// ShutdownServer
	shutdownServer(eh, t)
}

func TestITUpdateExpenses(t *testing.T) {
	// Setup server
	eh := initialServer()

	// Arrange:SetData
	expSendJSON := `{"title":"strawberry smoothie","amount":79,"note":"night market promotion discount 10 bath","tags":["food", "beverage"]}`
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%d/expenses", serverPort), strings.NewReader(expSendJSON))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client := http.Client{}

	// Act
	resp, err := client.Do(req)
	assert.NoError(t, err)

	byteBodyGot, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	resultStructGot := Expenses{}
	json.Unmarshal(byteBodyGot, &resultStructGot)

	resp.Body.Close()

	// Arrange:SetData
	ParamUpdateID := resultStructGot.ID
	expSendUpdateJSON := `{"title": "apple smoothie","amount": 89,"note": "no discount","tags": ["beverage"]}`
	expStructUpdateWant := Expenses{ID: 1, Title: "apple smoothie", Amount: 89, Note: "no discount", Tags: []string{"beverage"}}
	expUpdateWantJSON, _ := json.Marshal(expStructUpdateWant)
	req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("http://localhost:%d/expenses/%d", serverPort, ParamUpdateID), strings.NewReader(expSendUpdateJSON))
	assert.NoError(t, err)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	client = http.Client{}

	// Act
	resp, err = client.Do(req)
	assert.NoError(t, err)

	byteBodyUpdateGot, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	resultStructUpdateGot := Expenses{}
	json.Unmarshal(byteBodyUpdateGot, &resultStructUpdateGot)

	resp.Body.Close()

	// Assertions
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotEqual(t, expUpdateWantJSON, byteBodyUpdateGot)
		assert.Equal(t, expStructUpdateWant.Title, resultStructUpdateGot.Title)
		assert.Equal(t, expStructUpdateWant.Amount, resultStructUpdateGot.Amount)
		assert.Equal(t, expStructUpdateWant.Note, resultStructUpdateGot.Note)
		assert.Equal(t, expStructUpdateWant.Tags, resultStructUpdateGot.Tags)
	}

	// ShutdownServer
	shutdownServer(eh, t)
}
