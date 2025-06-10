package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = ""
	port     = 0
	user     = ""
	password = ""
	dbname   = ""
)

type SymbolData struct {
	Symbol       string  `json:"symbol"`
	HighestPrice float64 `json:"highest_price"`
	LowestPrice  float64 `json:"lowest_price"`
}

type SymbolDataResponse struct {
	Data          []SymbolData `json:"data"`
	LatestUpdated time.Time    `json:"latest_updated"`
}

type UserInfo struct {
	ID  int `json:"ID"`
	OTP int `json:"OTP"`
}

type StockWithEntryPrice struct {
	Symbol     string `json:"symbol"`
	EntryPrice int    `json:"entry_price"`
}

type UserTradeRequest struct {
	UserID   string                `json:"user_id"`
	Stocks   []StockWithEntryPrice `json:"stocks"`
	Operator string                `json:"operator"` // "Add", "Update", or "Delete"
}

type UserTradeResponse struct {
	Symbol        string  `json:"symbol"`
	EntryPrice    int     `json:"entry_price"`
	Signal        string  `json:"signal"`
	AvgPrice      int     `json:"avg_price"`
	CurrentPrice  int     `json:"current_price"`
	PercentChange float64 `json:"percent_change"`
}

// Add this struct definition
type UpdateSignalRequest struct {
	UserID         string `json:"user_id"`
	Symbol         string `json:"symbol"`
	BreakEvenPrice int    `json:"break_even_price"`
}

func getPotentialSymbolsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := db.Ping()
		if err != nil {
			http.Error(w, "Failed to ping database", http.StatusInternalServerError)
			log.Println("Failed to ping database:", err)
			return
		}

		rows, err := db.Query("SELECT symbol, highest_price, lowest_price FROM symbols_watchlist")
		if err != nil {
			http.Error(w, "Failed to query database", http.StatusInternalServerError)
			log.Println("Failed to query database:", err)
			return
		}
		defer rows.Close()

		var symbols []SymbolData
		for rows.Next() {
			var s SymbolData
			if err := rows.Scan(&s.Symbol, &s.HighestPrice, &s.LowestPrice); err != nil {
				http.Error(w, "Failed to scan row", http.StatusInternalServerError)
				log.Println("Failed to scan row:", err)
				return
			}
			symbols = append(symbols, s)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, "Error during row iteration", http.StatusInternalServerError)
			log.Println("Error during row iteration:", err)
			return
		}

		// Query to get the latest updated
		row := db.QueryRow("SELECT MAX(updated_at) FROM symbols_watchlist LIMIT 1")
		var latestUpdated time.Time
		if err = row.Scan(&latestUpdated); err != nil {
			http.Error(w, "Failed to scan row", http.StatusInternalServerError)
			log.Println("Failed to scan row:", err)
			return
		}

		response := SymbolDataResponse{
			Data:          symbols,
			LatestUpdated: latestUpdated,
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(response)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	http.HandleFunc("/getPotentialSymbols", getPotentialSymbolsHandler(db))
	http.HandleFunc("/health", healthCheck)
	fmt.Println("Server listening on :3000")
	addr := net.JoinHostPort("::", "3000")
	server := &http.Server{Addr: addr}
	log.Fatalln(server.ListenAndServe())
}
