package main_test

import (
	"testing"
	"os"
	"log"

	"github.com/juheyne/ShoppingLedger"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"bytes"
	"strconv"
	"fmt"
	"time"
)

var a main.App

func TestMain(m *testing.M) {
	a = main.App{}
	a.Initialize("./test.db")

	ensureTableExists()

	code := m.Run()

	clearTable()

	// Close DB and remove the sqlite db
	a.DB.Close()
	err := os.Remove("./test.db")
	if err != nil {
		fmt.Println("Could not remove test sqlite db due to: ", err)
	}

	os.Exit(code)
}

func ensureTableExists() {
	var name string
	if err := a.DB.QueryRow("SELECT name FROM sqlite_master WHERE name='expenses'").Scan(&name); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM expenses")
	// TODO Check if something similar to "ALTER SEQUENCE products_id_seq RESTART WITH 1" is necessary
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req,_ := http.NewRequest("GET", "/expenses", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	payload := []byte(`{"payer":"Payer","amount":53.12,"note":"Groceries","date":"2018-01-13T20:47:48+01:00"}`)

	req, _ := http.NewRequest("POST", "/expense", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["payer"] != "Payer" {
		t.Errorf("Expected payer to be 'Payer'. Got '%v'", m["payer"])
	}

	if m["amount"] != 53.12 {
		t.Errorf("Expected amount to be '53.12'. Got '%v'", m["amount"])
	}

	if m["note"] != "Groceries" {
		t.Errorf("Expected note to be 'Groceries'. Got '%v'", m["note"])
	}

	if m["date"] != "2018-01-13T20:47:48+01:00" {
		t.Errorf("Expected time to be '2018-01-13T20:47:48+01:00'. Got '%v'", m["date"])
	}

	if m["id"] != 1.0 {
		t.Errorf("Expected expense ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addExpenses(1)

	req, _ := http.NewRequest("GET", "/expense/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateExpense(t *testing.T) {
	clearTable()
	addExpenses(1)

	req, _ := http.NewRequest("GET", "/expense/1", nil)
	response := executeRequest(req)
	var originalExpense map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalExpense)

	payload := []byte(`{"payer":"New Guy","amount":53.12,"note":"Updated note","date":"2018-01-13T20:47:48+01:00"}`)

	req, _ = http.NewRequest("PUT", "/expense/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalExpense["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalExpense["id"], m["id"])
	}

	if m["payer"] == originalExpense["payer"] {
		t.Errorf("Expected the payer to change from '%v' to 'New Guy'. Got %v", originalExpense["payer"], m["payer"])
	}

	if m["amount"] == originalExpense["amount"] {
		t.Errorf("Expected the amount to change from '%v' to '53.12'. Got %v", originalExpense["amount"], m["amount"])
	}

	if m["note"] == originalExpense["note"] {
		t.Errorf("Expected the note to change from '%v' to 'Updated note'. Got '%v'", originalExpense["note"], m["note"])
	}

	if m["date"] == originalExpense["time"] {
		t.Errorf("Expected date to change from '%v' to '2018-01-13T20:47:48+01:00'. Got '%v'", originalExpense["date"], m["date"])
	}
}

func TestDeleteExpense(t *testing.T) {
	clearTable()
	addExpenses(1)

	req, _ := http.NewRequest("GET", "/expense/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/expense/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/expense/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func addExpenses(n int) {
	if n < 1 {
		n = 1
	}

	for i := 0; i < n; i++ {
		a.DB.Exec("INSERT INTO expenses(payer, amount, note, date) VALUES(?, ?, ?, ?)", fmt.Sprint("Payer ", i), (i+1.0)*10.0, strconv.Itoa(i), time.Now().Format(time.RFC3339))
	}
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}
