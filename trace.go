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
	"github.com/manifoldco/promptui"
	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "networks.db"

// API Response structures
type ApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		BlockNumber       string `json:"blockNumber"`
		BlockHash         string `json:"blockHash"`
		TimeStamp         string `json:"timeStamp"`
		Hash              string `json:"hash"`
		Nonce             string `json:"nonce"`
		TransactionIndex  string `json:"transactionIndex"`
		From              string `json:"from"`
		To                string `json:"to"`
		Value             string `json:"value"`
		Gas               string `json:"gas"`
		GasPrice          string `json:"gasPrice"`
		Input             string `json:"input"`
		ContractAddress   string `json:"contractAddress"`
		CumulativeGasUsed string `json:"cumulativeGasUsed"`
		TxreceiptStatus   string `json:"txreceipt_status"`
		GasUsed           string `json:"gasUsed"`
		Confirmations     string `json:"confirmations"`
		IsError           string `json:"isError"`
	} `json:"result"`
}

type TokenApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []struct {
		BlockNumber      string `json:"blockNumber"`
		TimeStamp        string `json:"timeStamp"`
		Hash             string `json:"hash"`
		Nonce            string `json:"nonce"`
		BlockHash        string `json:"blockHash"`
		From             string `json:"from"`
		To               string `json:"to"`
		Value            string `json:"value"`
		TokenName        string `json:"tokenName"`
		TokenSymbol      string `json:"tokenSymbol"`
		TokenDecimal     string `json:"tokenDecimal"`
		TransactionIndex string `json:"transactionIndex"`
		Gas              string `json:"gas"`
		GasUsed          string `json:"gasUsed"`
		GasPrice         string `json:"gasPrice"`
		Input            string `json:"input"`
		Confirmations    string `json:"confirmations"`
	} `json:"result"`
}

type Transaction struct {
	Pathway struct {
		SrcEid int `json:"srcEid"`
		DstEid int `json:"dstEid"`
		Sender struct {
			Address string `json:"address"`
			ID      string `json:"id"`
			Name    string `json:"name"`
			Chain   string `json:"chain"`
		} `json:"sender"`
		Receiver struct {
			Address string `json:"address"`
			ID      string `json:"id"`
			Name    string `json:"name"`
			Chain   string `json:"chain"`
		} `json:"receiver"`
		ID    string `json:"id"`
		Nonce int    `json:"nonce"`
	} `json:"pathway"`
	Source struct {
		Status string `json:"status"`
		Tx     struct {
			TxHash      string `json:"txHash"`
			BlockHash   string `json:"blockHash"`
			BlockNumber string `json:"blockNumber"`
			From        string `json:"from"`
			Payload     string `json:"payload"`
		} `json:"tx"`
	} `json:"source"`
}

type BridgeApiResponse struct {
	Data []Transaction `json:"data"`
}

type Network struct {
	ID     int
	Name   string
	APIURL string
	EnvKey string
}

type Contract struct {
	ID              int
	NetworkID       int
	ContractAddress string
	ContractName    string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("Error opening the database:", err)
	}
	defer db.Close()

	networks, err := fetchNetworks(db)
	if err != nil {
		log.Fatal("Error fetching networks:", err)
	}

	networkName := selectNetwork(networks)
	if networkName == "" {
		log.Fatal("No network selected")
	}

	selectedNetwork, err := getNetworkDetails(db, networkName)
	if err != nil {
		log.Fatal("Error fetching network details:", err)
	}

	contracts, err := fetchContracts(db, selectedNetwork.ID)
	if err != nil {
		log.Fatal("Error fetching contracts:", err)
	}

	addressPrompt := promptui.Prompt{
		Label: "Enter Wallet Address",
	}

	address, err := addressPrompt.Run()
	if err != nil {
		log.Fatalf("Error getting input: %v", err)
	}

	apiKey := os.Getenv(selectedNetwork.EnvKey)
	if apiKey == "" {
		log.Fatalf("API key not found for %s in .env file", selectedNetwork.EnvKey)
	}

	normalTxURL := fmt.Sprintf("%s?module=account&action=txlist&address=%s&startblock=0&endblock=99999999&page=1&offset=10&sort=desc&apikey=%s",
		selectedNetwork.APIURL, address, apiKey)
	internalTxURL := fmt.Sprintf("%s?module=account&action=txlistinternal&address=%s&startblock=0&endblock=99999999&page=1&offset=10&sort=desc&apikey=%s",
		selectedNetwork.APIURL, address, apiKey)
	tokenTxURL := fmt.Sprintf("%s?module=account&action=tokentx&address=%s&startblock=0&endblock=99999999&page=1&offset=10&sort=desc&apikey=%s",
		selectedNetwork.APIURL, address, apiKey)

	fmt.Println("\n--- Matched Normal Transactions ---")
	fetchAndDisplayMatchingTransactions(normalTxURL, contracts)

	fmt.Println("\n--- Matched Internal Transactions ---")
	fetchAndDisplayMatchingTransactions(internalTxURL, contracts)

	fmt.Println("\n--- Matched Token Transactions ---")
	fetchAndDisplayMatchingTokenTransactions(tokenTxURL, contracts)

	lookup(address)
}

func fetchNetworks(db *sql.DB) ([]Network, error) {
	rows, err := db.Query("SELECT id, name FROM networks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var networks []Network
	for rows.Next() {
		var network Network
		if err := rows.Scan(&network.ID, &network.Name); err != nil {
			return nil, err
		}
		networks = append(networks, network)
	}

	return networks, nil
}

func fetchContracts(db *sql.DB, networkID int) ([]Contract, error) {
	rows, err := db.Query("SELECT id, network_id, LOWER(contract_address), contract_name FROM contracts WHERE network_id = ?", networkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contracts []Contract
	for rows.Next() {
		var contract Contract
		if err := rows.Scan(&contract.ID, &contract.NetworkID, &contract.ContractAddress, &contract.ContractName); err != nil {
			return nil, err
		}
		contracts = append(contracts, contract)
	}

	return contracts, nil
}

func selectNetwork(networks []Network) string {
	names := []string{}
	for _, network := range networks {
		names = append(names, network.Name)
	}

	prompt := promptui.Select{
		Label: "Select Network",
		Items: names,
	}

	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Error selecting network: %v", err)
	}

	return result
}

func getNetworkDetails(db *sql.DB, networkName string) (*Network, error) {
	var network Network
	query := "SELECT id, name, api_url, env_key FROM networks WHERE name = ?"
	err := db.QueryRow(query, networkName).Scan(&network.ID, &network.Name, &network.APIURL, &network.EnvKey)
	if err != nil {
		return nil, err
	}
	return &network, nil
}

func fetchAndDisplayMatchingTransactions(url string, contracts []Contract) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error calling the API:", err)
	}
	defer resp.Body.Close()

	var response ApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Error decoding the response:", err)
	}

	if response.Status != "1" {
		log.Fatalf("API request failed: %s", response.Message)
	}

	for _, tx := range response.Result {
		for _, contract := range contracts {
			if tx.To == contract.ContractAddress || tx.From == contract.ContractAddress {
				fmt.Printf("\nMatched Transaction\n")
				fmt.Printf("Block Number: %s\n", tx.BlockNumber)
				fmt.Printf("Hash: %s\n", tx.Hash)
				fmt.Printf("From: %s\n", tx.From)
				fmt.Printf("To: %s\n", tx.To)
				fmt.Printf("Value: %s\n", tx.Value)
				fmt.Println("-----------")
			}
		}
	}
}

func fetchAndDisplayMatchingTokenTransactions(url string, contracts []Contract) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error calling the API:", err)
	}
	defer resp.Body.Close()

	var response TokenApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Fatal("Error decoding the response:", err)
	}

	if response.Status != "1" {
		log.Fatalf("API request failed: %s", response.Message)
	}

	for _, tx := range response.Result {
		for _, contract := range contracts {
			if tx.To == contract.ContractAddress || tx.From == contract.ContractAddress {
				fmt.Printf("\nMatched Token Transaction\n")
				fmt.Printf("Block Number: %s\n", tx.BlockNumber)
				fmt.Printf("Hash: %s\n", tx.Hash)
				fmt.Printf("From: %s\n", tx.From)
				fmt.Printf("To: %s\n", tx.To)
				fmt.Printf("Value: %s\n", tx.Value)
				fmt.Printf("Token Name: %s\n", tx.TokenName)
				fmt.Printf("Token Symbol: %s\n", tx.TokenSymbol)
				fmt.Println("-----------")
			}
		}
	}
}

func lookup(walletAddress string) ([]Transaction, error) {
	// Construct the API URL
	url := fmt.Sprintf("https://scan.layerzero-api.com/v1/messages/wallet/%s?limit=100", walletAddress)

	// Send the GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON response
	var apiResponse BridgeApiResponse
	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		return nil, err
	}

	// Print out the transactions and extract wallet address
	for _, tx := range apiResponse.Data {
		// Extract the wallet address involved in the transaction
		senderAddress := tx.Pathway.Sender.Address
		receiverAddress := tx.Pathway.Receiver.Address

		// Print transaction details
		fmt.Printf("Transaction ID: %s\n", tx.Pathway.ID)
		fmt.Printf("Sender: %s (%s) - Wallet Address: %s\n", tx.Pathway.Sender.Name, tx.Pathway.Sender.Chain, senderAddress)
		fmt.Printf("Receiver: %s (%s) - Wallet Address: %s\n", tx.Pathway.Receiver.Name, tx.Pathway.Receiver.Chain, receiverAddress)
		fmt.Printf("Status: %s\n", tx.Source.Status)
		fmt.Println("----------------------------------------------------")
	}

	// Return the transactions and no error
	return apiResponse.Data, nil
}
