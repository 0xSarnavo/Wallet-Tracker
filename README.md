# Blockchain Transaction Lookup

This Go program facilitates the lookup of blockchain transactions on different networks using API endpoints. It fetches and displays the transactions for a specified wallet address, including normal transactions, internal transactions, and token transactions. Additionally, it supports transactions from LayerZero-based cross-chain bridges.

## Features

1. **Network Selection**: The program allows the user to select a network from a list of available networks stored in an SQLite database.
2. **Wallet Address Input**: The user is prompted to enter a wallet address for querying transactions.
3. **API Interaction**: The program interacts with blockchain APIs to fetch transaction data based on the selected network.
4. **Transaction Matching**: The transactions fetched from the blockchain are matched with contract addresses stored in the database, and relevant transactions are displayed.
5. **Cross-Chain Bridge Lookup**: Fetches transaction details related to a wallet from the LayerZero API, which tracks cross-chain interactions.
6. **SQLite Database**: Uses SQLite to store and fetch network and contract details.

## Dependencies

- `github.com/joho/godotenv` – Load environment variables from a `.env` file.
- `github.com/manifoldco/promptui` – A Go package for building command-line prompts.
- `github.com/mattn/go-sqlite3` – SQLite driver for Go.

## Setup

### Step 1: Install Dependencies

Run the following command to install the required dependencies:
```bash
go mod tidy
```

### Step 2: Set up the `.env` File

Create a `.env` file at the root of the project to store the API keys for each network. The `.env` file should look like this:
```
NETWORK_API_KEY=<API_KEY_FOR_SELECTED_NETWORK>
```

### Step 3: Set up the SQLite Database

The program uses an SQLite database (`networks.db`) to store network and contract information. You need to set up the database with the following tables:

- **Networks**: Contains network details like name, API URL, and environment key.
- **Contracts**: Contains contract addresses and names for each network.

### Step 4: Run the Program

To start the program, run:
```bash
go run main.go
```

### Step 5: Select a Network

You will be prompted to select a network from the available list stored in the database. After selecting the network, you will need to enter a wallet address to fetch transactions.

### Step 6: Display Transactions

The program will display the following types of transactions for the given wallet address:

- **Matched Normal Transactions**: Transactions related to the given wallet address.
- **Matched Internal Transactions**: Internal transactions associated with the wallet.
- **Matched Token Transactions**: Token transactions for the wallet.

Additionally, the program will display transactions from LayerZero-based cross-chain bridges.

## Code Breakdown

### 1. **Environment Variable Loading**
The `godotenv` package loads the environment variables from the `.env` file, ensuring the program has access to the required API keys.

### 2. **Database Operations**
The program interacts with an SQLite database (`networks.db`) to fetch network and contract details. It performs the following database operations:
- Fetch all networks stored in the `networks` table.
- Fetch contracts related to the selected network from the `contracts` table.

### 3. **Transaction Fetching**
The program constructs API URLs based on the selected network and wallet address. It uses the `http.Get` method to make requests to the API and fetch transaction data in the following formats:
- **Normal Transactions**: General wallet transactions.
- **Internal Transactions**: Transactions involving contract interactions.
- **Token Transactions**: Transactions involving tokens.

### 4. **Transaction Matching**
For each transaction fetched, the program checks if the transaction's `from` or `to` address matches any contract addresses from the database. If a match is found, the transaction details are displayed.

### 5. **LayerZero Bridge Lookup**
The program also fetches transaction details from the LayerZero API for the given wallet address. This section fetches cross-chain transactions and prints detailed information about the sender, receiver, and status of each transaction.

## Example Output

```bash
--- Matched Normal Transactions ---
Block Number: 123456
Hash: 0xabc123...
From: 0xdef456...
To: 0x123456...
Value: 1.5 ETH
-----------

--- Matched Internal Transactions ---
Block Number: 123457
Hash: 0xdef123...
From: 0x123456...
To: 0x789012...
Value: 0.3 ETH
-----------

--- Matched Token Transactions ---
Block Number: 123458
Hash: 0xabc123...
From: 0x987654...
To: 0x123456...
Value: 1000 Tokens
Token Name: MyToken
Token Symbol: MTK
-----------
```

## Notes

- Ensure that the `networks.db` file contains valid data before running the program.
- The program is designed to be extendable for additional networks or transaction types.
```

This README should provide clarity on the functionality, setup, and operation of the code.