package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHealthCheck(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(healthCheck)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestGetPotentialSymbols(t *testing.T) {
	// Create a new mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Mock the expected query and its result
	rows := sqlmock.NewRows([]string{"symbol", "highest_price", "lowest_price"}).
		AddRow("AAPL", 150.0, 140.0).
		AddRow("GOOG", 2500.0, 2400.0)

	mock.ExpectQuery("SELECT symbol, highest_price, lowest_price FROM symbols_watchlist").WillReturnRows(rows)

	// Mock the query for the latest updated time
	now := time.Now()
	timeRow := sqlmock.NewRows([]string{"max"}).AddRow(now)
	mock.ExpectQuery("SELECT MAX\\(updated_at\\) FROM symbols_watchlist LIMIT 1").WillReturnRows(timeRow)

	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/getPotentialSymbols", nil)
	if err != nil {
		t.Fatal(err)
	}

	// We create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the refactored handler function with the mock DB
	handler := getPotentialSymbolsHandler(db)
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response SymbolDataResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("could not unmarshal response body: %v", err)
	}

	if len(response.Data) != 2 {
		t.Errorf("handler returned unexpected number of symbols: got %v want %v",
			len(response.Data), 2)
	}

	// You can add more specific checks for the data content if needed
	// For example:
	// if response.Data[0].Symbol != "AAPL" { ... }

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
