package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Define structures to parse the response
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

	db, err := sql.Open("sqlite3", "./network_contracts.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
    }

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

    for _, table := range networkNames {
        log.Printf("Network Name: %s", table)
    }

	for _, table := range networkApis {
        log.Printf("Network APIs: %s", table)
    }

	for _, table := range networkEnvs {
        log.Printf("Network APIs: %s", table)
    }
/**
	// Prompt for API key, address, and the base URL
	var apiKey, address, baseURL string

	fmt.Print("Enter the API base URL (e.g., https://api.snowscan.xyz/api): ")
	fmt.Scanln(&baseURL)

	fmt.Print("Enter your API key: ")
	fmt.Scanln(&apiKey)

	fmt.Print("Enter the address to query (e.g., 0xaddress): ")
	fmt.Scanln(&address)

	// Construct the API URL
	url := fmt.Sprintf("%s?module=account&action=txlistinternal&startblock=13028500&endblock=13028600&page=1&offset=10&sort=asc&apikey=%s", baseURL, apiKey)

	// Call the API
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error calling API: %v", err)
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Parse JSON response
	var apiResp APIResponse
	err = json.Unmarshal(body, &apiResp)
	if err != nil {
		log.Fatalf("Error unmarshaling response: %v", err)
	}

	// Check if status is success
	if apiResp.Status == "1" {
		fmt.Println("\nTransaction Details:")
		for _, tx := range apiResp.Result {
			fmt.Printf("\nBlock Number: %s\n", tx.BlockNumber)
			fmt.Printf("Timestamp: %s\n", tx.TimeStamp)
			fmt.Printf("Transaction Hash: %s\n", tx.Hash)
			fmt.Printf("From: %s\n", tx.From)
			fmt.Printf("To: %s\n", tx.To)
			fmt.Printf("Value: %s\n", tx.Value)
			fmt.Printf("Contract Address: %s\n", tx.ContractAddr)
			fmt.Printf("Input: %s\n", tx.Input)
			fmt.Printf("Type: %s\n", tx.Type)
			fmt.Printf("Gas: %s\n", tx.Gas)
			fmt.Printf("Gas Used: %s\n", tx.GasUsed)
			fmt.Printf("Trace ID: %s\n", tx.TraceID)
			fmt.Printf("Is Error: %s\n", tx.IsError)
			fmt.Printf("Error Code: %s\n", tx.ErrCode)
			fmt.Println("--------------------------------------------------------")
		}
	} else {
		fmt.Println("Error in response:", apiResp.Message)
	}
		**/
}

func fetchNetworkName(db *sql.DB) ([]string, error) {
    rows, err := db.Query("SELECT name FROM networks")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var network []string
    for rows.Next() {
        var networkName string
        if err := rows.Scan(&networkName); err != nil {
            return nil, err
        }
        network = append(network, networkName)
    }

    return network, rows.Err()
}

func fetchNetworkApi(db *sql.DB) ([]string, error) {
    rows, err := db.Query("SELECT api_link FROM networks")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var api_links []string
    for rows.Next() {
        var api string
        if err := rows.Scan(&api); err != nil {
            return nil, err
        }
        api_links = append(api_links, api)
    }

    return api_links, rows.Err()
}

func fetchNetworkEnv(db *sql.DB) ([]string, error) {
    rows, err := db.Query("SELECT apikey_env FROM networks")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var env_names []string
    for rows.Next() {
        var env string
        if err := rows.Scan(&env); err != nil {
            return nil, err
        }
        env_names = append(env_names, env)
    }

    return env_names, rows.Err()
}