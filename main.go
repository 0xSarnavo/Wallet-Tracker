package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	dbmanager "wallet-tracker/Database"
)

func main() {
	db, err := dbmanager.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nChoose an option:")
		fmt.Println("1. List Networks")
		fmt.Println("2. Add Network")
		fmt.Println("3. Remove Network")
		fmt.Println("4. List Contracts by Network")
		fmt.Println("5. Add Contract")
		fmt.Println("6. Remove Contract")
		fmt.Println("7. Exit")
		fmt.Print("Enter your choice: ")

		choiceStr, _ := reader.ReadString('\n')
		choiceStr = strings.TrimSpace(choiceStr)
		choice, err := strconv.Atoi(choiceStr)
		if err != nil {
			fmt.Println("Invalid choice. Please enter a number.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("\nNetworks in the database:")
			dbmanager.ListNetworks(db)

		case 2:
			fmt.Print("\nEnter network name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Enter Chain ID: ")
			chainidStr, _ := reader.ReadString('\n')
			chainid, err := strconv.Atoi(strings.TrimSpace(chainidStr))
			if err != nil {
				fmt.Println("Invalid Chain ID. Please enter a number.")
				continue
			}

			fmt.Print("Enter API link: ")
			apiLink, _ := reader.ReadString('\n')
			apiLink = strings.TrimSpace(apiLink)

			fmt.Print("Enter API key name in .env: ")
			apiKey, _ := reader.ReadString('\n')
			apiKey = strings.TrimSpace(apiKey)

			dbmanager.AddNetwork(db, name, chainid, apiLink, apiKey)

		case 3:
			fmt.Print("\nEnter network name to remove: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)
			dbmanager.RemoveNetwork(db, name)

		case 4:
			fmt.Print("\nEnter network name to list contracts: ")
			networkName, _ := reader.ReadString('\n')
			networkName = strings.TrimSpace(networkName)
			dbmanager.ListContractsByNetwork(db, networkName)

		case 5:
			fmt.Println("\nAvailable Networks:")
			dbmanager.ListNetworks(db)

			fmt.Print("\nEnter network name: ")
			networkName, _ := reader.ReadString('\n')
			networkName = strings.TrimSpace(networkName)

			fmt.Print("Enter contract address: ")
			contractAddress, _ := reader.ReadString('\n')
			contractAddress = strings.TrimSpace(contractAddress)

			fmt.Print("Enter contract name: ")
			contractName, _ := reader.ReadString('\n')
			contractName = strings.TrimSpace(contractName)

			dbmanager.AddContract(db, networkName, contractAddress, contractName)

		case 6:
			fmt.Print("\nEnter contract address to remove: ")
			contractAddress, _ := reader.ReadString('\n')
			contractAddress = strings.TrimSpace(contractAddress)
			dbmanager.RemoveContract(db, contractAddress)

		case 7:
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
