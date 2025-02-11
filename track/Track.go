package track

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Network struct {
	Name      string
	ChainID   int
	APILink   string
	APIKeyEnv string
}

type Transaction struct {
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	TxReceiptStatus string `json:"txreceipt_status"`
	Input           string `json:"input"`
	ContractAddress string `json:"contractAddress"`
	FunctionName    string `json:"functionName"`
	Confirmations   string `json:"confirmations"`
}

type APIResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Result  []Transaction `json:"result"`
}

func getNetworks(db *sql.DB) ([]Network, error) {
	rows, err := db.Query("SELECT name, chainid, api_link, apikey_env FROM networks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var networks []Network
	for rows.Next() {
		var network Network
		rows.Scan(&network.Name, &network.ChainID, &network.APILink, &network.APIKeyEnv)
		networks = append(networks, network)
	}
	return networks, nil
}

func getContracts(db *sql.DB) (map[string]string, error) {
	rows, err := db.Query("SELECT contract_address, name FROM contracts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	contracts := make(map[string]string)
	for rows.Next() {
		var address, name string
		rows.Scan(&address, &name)
		contracts[address] = name
	}
	return contracts, nil
}

func fetchTransactions(network Network, address string, contracts map[string]string) {
	apiKey := os.Getenv(network.APIKeyEnv)
	if apiKey == "" {
		log.Printf("API key not found for %s", network.Name)
		return
	}

	url := fmt.Sprintf("%s?module=account&action=balancehistory&address=%s&apikey=%s", network.APILink, address, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch transactions: %v", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var apiResp APIResponse
	json.Unmarshal(body, &apiResp)

	if apiResp.Status != "1" {
		log.Println("Error fetching transactions")
		return
	}

	for _, tx := range apiResp.Result {
		if name, exists := contracts[tx.To]; exists {
			fmt.Printf("\nNetwork: %s\nContract: %s\nConfirmations: %s\nBlock: %s\nTimestamp: %s\nHash: %s\nFrom: %s\nTo: %s\nStatus: %s\nInput: %s\nContract Address: %s\nFunction: %s\n", 
				network.Name, name, tx.Confirmations, tx.BlockNumber, tx.TimeStamp, tx.Hash, tx.From, tx.To, tx.TxReceiptStatus, tx.Input, tx.ContractAddress, tx.FunctionName)
		}
	}
}

func main() {
	// Connect to database
	db, err := sql.Open("sqlite3", "network_contracts.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Get networks & contracts
	networks, err := getNetworks(db)
	if err != nil {
		log.Fatal(err)
	}

	contracts, err := getContracts(db)
	if err != nil {
		log.Fatal(err)
	}

	// Get user input
	var address string
	fmt.Print("Enter wallet address: ")
	fmt.Scanln(&address)

	// Fetch transactions for all networks
	for _, network := range networks {
		fetchTransactions(network, address, contracts)
	}
}
