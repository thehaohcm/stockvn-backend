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

func TestCorsMiddleware(t *testing.T) {
	// Create a dummy handler to be wrapped by the middleware
	mockHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Create a new recorder to capture the response
	rr := httptest.NewRecorder()

	// Create a request for a GET method
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Apply the middleware
	handler := corsMiddleware(mockHandler)
	handler.ServeHTTP(rr, req)

	// Check CORS headers for a regular GET request
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("CORS middleware failed for GET: Access-Control-Allow-Origin header missing or incorrect")
	}
	if rr.Header().Get("Access-Control-Allow-Methods") != "GET, POST, PUT, DELETE, OPTIONS" {
		t.Errorf("CORS middleware failed for GET: Access-Control-Allow-Methods header missing or incorrect")
	}
	if rr.Header().Get("Access-Control-Allow-Headers") != "Content-Type, Authorization" {
		t.Errorf("CORS middleware failed for GET: Access-Control-Allow-Headers header missing or incorrect")
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for GET: got %v want %v", status, http.StatusOK)
	}

	// Test preflight OPTIONS request
	rr = httptest.NewRecorder() // Reset recorder
	req, err = http.NewRequest("OPTIONS", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	handler.ServeHTTP(rr, req)

	// Check status code for OPTIONS request
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code for OPTIONS: got %v want %v", status, http.StatusOK)
	}
	// For OPTIONS, only the CORS headers should be set, and no body from the mockHandler
	if rr.Body.String() != "" {
		t.Errorf("handler returned unexpected body for OPTIONS: got %v want empty", rr.Body.String())
	}
}
