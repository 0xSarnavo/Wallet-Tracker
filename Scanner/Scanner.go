package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Result struct {
	BlockNumber  string `json:"blockNumber"`
	TimeStamp    string `json:"timeStamp"`
	Hash         string `json:"hash"`
	From         string `json:"from"`
	To           string `json:"to"`
	Value        string `json:"value"`
	ContractAddr string `json:"contractAddress"`
	Input        string `json:"input"`
	Type         string `json:"type"`
	Gas          string `json:"gas"`
	GasUsed      string `json:"gasUsed"`
	TraceID      string `json:"traceId"`
	IsError      string `json:"isError"`
	ErrCode      string `json:"errCode"`
}

type APIResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  []Result `json:"result"`
}

func main() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Open SQLite database
	db, err := sql.Open("sqlite3", "./network_contracts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test DB connection
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Fetch networks, APIs, and API key environment names
	networkNames, err := fetchNetworkName(db)
	if err != nil {
		log.Fatal(err)
	}

	networkApis, err := fetchNetworkApi(db)
	if err != nil {
		log.Fatal(err)
	}

	networkEnvs, err := fetchNetworkEnv(db)
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through networks and fetch transactions
	for i, network := range networkNames {
		apiKey := os.Getenv(networkEnvs[i])
		if apiKey == "" {
			log.Printf("API key for %s is missing in .env\n", network)
			continue
		}

		fetchTransactions(network, networkApis[i], apiKey)
	}
}

// Fetch transactions from API
func fetchTransactions(network, baseURL, apiKey string) {
	url := fmt.Sprintf("%s?module=account&action=txlistinternal&startblock=13028500&endblock=13028600&page=1&offset=10&sort=asc&apikey=%s", baseURL, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error calling API for %s: %v\n", network, err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response for %s: %v\n", network, err)
		return
	}

	var apiResp APIResponse
	if err = json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("Error unmarshaling response for %s: %v\n", network, err)
		return
	}

	if apiResp.Status == "1" {
		log.Printf("\nTransactions for %s:\n", network)
		for _, tx := range apiResp.Result {
			log.Printf("Block: %s, From: %s, To: %s, Value: %s, Hash: %s\n", tx.BlockNumber, tx.From, tx.To, tx.Value, tx.Hash)
		}
	} else {
		log.Printf("API Error for %s: %s\n", network, apiResp.Message)
	}
}

func fetchNetworkName(db *sql.DB) ([]string, error) {
	return fetchColumn(db, "SELECT name FROM networks")
}

func fetchNetworkApi(db *sql.DB) ([]string, error) {
	return fetchColumn(db, "SELECT api_link FROM networks")
}

func fetchNetworkEnv(db *sql.DB) ([]string, error) {
	return fetchColumn(db, "SELECT apikey_env FROM networks")
}

func fetchColumn(db *sql.DB, query string) ([]string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, rows.Err()
}
