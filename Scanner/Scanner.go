package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type NormalTxResult struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	TransactionIndex  string `json:"transactionIndex"`
	From              string `json:"from"`
	To                string `json:"to"`
	Value             string `json:"value"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	IsError           string `json:"isError"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	Input             string `json:"input"`
	ContractAddress   string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed           string `json:"gasUsed"`
	Confirmations     string `json:"confirmations"`
}

type NormalTxAPIResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []NormalTxResult `json:"result"`
}
type InternalTxResult struct {
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

type InternalTxAPIResponse struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Result  []InternalTxResult `json:"result"`
}

type Erc20TxResult struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Input             string `json:"input"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Confirmations     string `json:"confirmations"`
}

type Erc20TxAPIResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Result  []Erc20TxResult `json:"result"`
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

	fmt.Print("Enter Ethereum address: ")
	reader := bufio.NewReader(os.Stdin)
	userAddress, _ := reader.ReadString('\n')
	userAddress = strings.TrimSpace(userAddress)

	for i, network := range networkNames {
		apiKey := os.Getenv(networkEnvs[i])
		if apiKey == "" {
			log.Printf("API key for %s is missing in .env\n", network)
			continue
		}

		fetchTransactions(network, networkApis[i], apiKey, userAddress)
	}
}


// Fetch transactions from API
func fetchTransactions(network, baseURL, apiKey, address string) {
	
	if address == "" {
		log.Println("No address provided. Skipping API call.")
		return
	}

	// normal transactions
	normalTxUrl := fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=10&sort=asc&apikey=%s", baseURL, address, apiKey)

	normalResp, err := http.Get(normalTxUrl)
	if err != nil {
		log.Printf("Error calling API for %s: %v\n", network, err)
		return
	}
	defer normalResp.Body.Close()

	normalBody, err := ioutil.ReadAll(normalResp.Body)
	if err != nil {
		log.Printf("Error reading response for %s: %v\n", network, err)
		return
	}

	var normalApiResp NormalTxAPIResponse
	if err = json.Unmarshal(normalBody, &normalApiResp); err != nil {
		log.Printf("Error unmarshaling response for %s: %v\n", network, err)
		return
	}

	if normalApiResp.Status == "1" {
		log.Printf("\nTransactions for %s:\n", network)
		for _, tx := range normalApiResp.Result {
			log.Printf("Block: %s, From: %s, To: %s, Value: %s, Hash: %s\n", tx.BlockNumber, tx.From, tx.To, tx.Value, tx.Hash)
		}
	} else {
		log.Printf("API Error for %s: %s\n", network, normalApiResp.Message)
	}

	//internal transactions
	internalTxUrl := fmt.Sprintf("%s?module=account&action=txlistinternal&startblock=0&endblock=99999999&page=1&offset=10&sort=asc&apikey=%s", baseURL, apiKey)

	internalResp, err := http.Get(internalTxUrl)
	if err != nil {
		log.Printf("Error calling API for %s: %v\n", network, err)
		return
	}
	defer internalResp.Body.Close()

	internalBody, err := ioutil.ReadAll(internalResp.Body)
	if err != nil {
		log.Printf("Error reading response for %s: %v\n", network, err)
		return
	}

	var internalApiResp InternalTxAPIResponse
	if err = json.Unmarshal(internalBody, &internalApiResp); err != nil {
		log.Printf("Error unmarshaling response for %s: %v\n", network, err)
		return
	}

	if internalApiResp.Status == "1" {
		log.Printf("\nTransactions for %s:\n", network)
		for _, tx := range internalApiResp.Result {
			log.Printf("Block: %s, From: %s, To: %s, Value: %s, Hash: %s\n", tx.BlockNumber, tx.From, tx.To, tx.Value, tx.Hash)
		}
	} else {
		log.Printf("API Error for %s: %s\n", network, internalApiResp.Message)
	}

	//erc20 transactions
	erc20TxUrl := fmt.Sprintf("%s?module=account&action=tokentx&address=%s&startblock=0&endblock=99999999&page=1&offset=10&sort=asc&apikey=%s", baseURL, address, apiKey)

	erc20Resp, err := http.Get(erc20TxUrl)
	if err != nil {
		log.Printf("Error calling API for %s: %v\n", network, err)
		return
	}
	defer erc20Resp.Body.Close()

	erc20Body, err := ioutil.ReadAll(erc20Resp.Body)
	if err != nil {
		log.Printf("Error reading response for %s: %v\n", network, err)
		return
	}

	var erc20ApiResp Erc20TxAPIResponse
	if err = json.Unmarshal(erc20Body, &erc20ApiResp); err != nil {
		log.Printf("Error unmarshaling response for %s: %v\n", network, err)
		return
	}

	if erc20ApiResp.Status == "1" {
		log.Printf("\nTransactions for %s:\n", network)
		for _, tx := range erc20ApiResp.Result {
			log.Printf("Block: %s, From: %s, To: %s, Value: %s, Hash: %s\n", tx.BlockNumber, tx.From, tx.To, tx.Value, tx.Hash)
		}
	} else {
		log.Printf("API Error for %s: %s\n", network, erc20ApiResp.Message)
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
