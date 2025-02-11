package dbmanager

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "network_contracts.db")
	if err != nil {
		return nil, err
	}
	createTables(db)
	return db, nil
}

func createTables(db *sql.DB) {
	networkTable := `CREATE TABLE IF NOT EXISTS networks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE NOT NULL,
		chainid INTEGER UNIQUE NOT NULL,
		api_link TEXT NOT NULL,
		apikey_env TEXT NOT NULL
	);`

	contractTable := `CREATE TABLE IF NOT EXISTS contracts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		network_name TEXT NOT NULL,
		chainid INTEGER NOT NULL,
		contract_address TEXT UNIQUE NOT NULL,
		name TEXT NOT NULL,
		FOREIGN KEY(network_name) REFERENCES networks(name),
		FOREIGN KEY(chainid) REFERENCES networks(chainid)
	);`

	_, err := db.Exec(networkTable)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(contractTable)
	if err != nil {
		log.Fatal(err)
	}
}

func AddNetwork(db *sql.DB, name string, chainid int, apiLink, apikey string) {
	_, err := db.Exec("INSERT INTO networks (name, chainid, api_link, apikey_env) VALUES (?, ?, ?, ?)", name, chainid, apiLink, apikey)
	if err != nil {
		log.Println("Error adding network:", err)
	} else {
		fmt.Println("Network added successfully!")
	}
}

func RemoveNetwork(db *sql.DB, name string) {
	_, err := db.Exec("DELETE FROM networks WHERE name = ?", name)
	if err != nil {
		log.Println("Error removing network:", err)
	} else {
		fmt.Println("Network removed successfully!")
	}
}

func AddContract(db *sql.DB, networkName, contractAddress, contractName string) {
	var chainid int
	err := db.QueryRow("SELECT chainid FROM networks WHERE name = ?", networkName).Scan(&chainid)
	if err != nil {
		fmt.Println("Network not found!")
		return
	}

	_, err = db.Exec("INSERT INTO contracts (network_name, chainid, contract_address, name) VALUES (?, ?, ?, ?)", networkName, chainid, contractAddress, contractName)
	if err != nil {
		log.Println("Error adding contract:", err)
	} else {
		fmt.Println("Contract added successfully!")
	}
}

func RemoveContract(db *sql.DB, contractAddress string) {
	_, err := db.Exec("DELETE FROM contracts WHERE contract_address = ?", contractAddress)
	if err != nil {
		log.Println("Error removing contract:", err)
	} else {
		fmt.Println("Contract removed successfully!")
	}
}

func ListNetworks(db *sql.DB) {
	rows, err := db.Query("SELECT name, chainid, api_link, apikey_env FROM networks")
	if err != nil {
		log.Println("Error fetching networks:", err)
		return
	}
	defer rows.Close()

	fmt.Println("Available Networks:")
	for rows.Next() {
		var name, apiLink, apikey string
		var chainid int
		rows.Scan(&name, &chainid, &apiLink, &apikey)
		fmt.Printf("Name: %s, ChainID: %d, API: %s, APIKey: %s\n", name, chainid, apiLink, apikey)
	}
}

func ListContractsByNetwork(db *sql.DB, networkName string) {
	rows, err := db.Query("SELECT contract_address, name FROM contracts WHERE network_name = ?", networkName)
	if err != nil {
		log.Println("Error fetching contracts:", err)
		return
	}
	defer rows.Close()

	fmt.Printf("Contracts under network %s:\n", networkName)
	for rows.Next() {
		var contractAddress, name string
		rows.Scan(&contractAddress, &name)
		fmt.Printf("Address: %s, Name: %s\n", contractAddress, name)
	}
}