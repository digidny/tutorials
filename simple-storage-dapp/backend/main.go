package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	contract "github.com/digidny/simple-storage-dapp/backend/internal/contract/storage" // Import the generated contract binding package.
	"github.com/joho/godotenv"                                                          // For loading environment variables from a .env file.
	"github.com/jumbochain/jumbochain-go"                                               // Main Jumbochain Go library.
	"github.com/jumbochain/jumbochain-go/accounts/abi"                                  // For interacting with contract ABIs.
	"github.com/jumbochain/jumbochain-go/accounts/abi/bind"                             // For binding Go code to smart contracts.
	"github.com/jumbochain/jumbochain-go/common"                                        // Common data types and utilities.
	"github.com/jumbochain/jumbochain-go/core/types"                                    // Ethereum transaction and receipt types.
	"github.com/jumbochain/jumbochain-go/crypto"                                        // Cryptographic utilities for key management.
	"github.com/jumbochain/jumbochain-go/jumboclient"                                   // Ethereum client interface.
)

// Ensure this matches the contract ABI.  Use `abigen` to generate.
//go:generate abigen --abi=../build/SimpleStorage.abi --pkg=storage --out=./internal/contract/storage/storage.go

// Storage struct holds information about the deployed contract.
type Storage struct {
	address common.Address // The address of the deployed smart contract.
	abi     string         // The ABI (Application Binary Interface) of the smart contract.
}

func main() {
	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Connect to the Jumbochain client. Use the URL from the environment.
	rpcURL := os.Getenv("RPC_URL") // e.g., "http://localhost:8545"
	if rpcURL == "" {
		log.Fatal("RPC_URL environment variable not set")
	}
	client, err := jumboclient.DialContext(context.Background(), rpcURL)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close() // Ensure the client connection is closed when the function exits.

	// 1. Get the contract address from the environment.
	contractAddressStr := os.Getenv("CONTRACT_ADDRESS")
	if contractAddressStr == "" {
		log.Fatal("CONTRACT_ADDRESS environment variable not set")
	}
	contractAddress := common.HexToAddress(contractAddressStr) // Convert the hex string address to a common.Address type.
	fmt.Println("Contract Address:", contractAddress)

	// 2. Create an instance of the contract binding.
	instance, err := contract.NewStorage(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Get the initial value stored in the contract.
	initialValue, err := instance.Get(nil) // `nil` indicates we are making a read-only call and don't need transaction options.
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Initial value:", initialValue)

	// 4. Set a new value in the contract.
	// Parse the contract ABI. This is needed to pack the function arguments.
	contractABI, err := abi.JSON(strings.NewReader(contract.StorageABI))
	if err != nil {
		log.Fatal("Error parsing contract ABI:", err)
	}

	// Pack the arguments for the `set` function. The `set` function in the
	// SimpleStorage contract likely takes one uint256 argument.
	newValue := big.NewInt(150)
	fmt.Printf("Setting a new value: %d\n", newValue)

	txData, err := contractABI.Pack("set", newValue) // "set" is the name of the function in the contract.
	if err != nil {
		log.Fatal("Error packing arguments for 'set':", err)
	}

	// Get transaction authorization options for sending a transaction.
	auth, err := getTransactionAuthorizer(client, contractAddress, txData)
	if err != nil {
		log.Fatal(err)
	}

	// Send the transaction to set the new value.
	tx, err := instance.Set(auth, newValue)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Set transaction hash: %s\n", tx.Hash().Hex())

	// Wait for the transaction to be mined on the blockchain.
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatalf("Transaction %s mining failed: %v", tx.Hash().Hex(), err)
	}

	// Check if the transaction was successful.
	if receipt.Status == types.ReceiptStatusFailed {
		log.Fatalf("Transaction %s failed", tx.Hash().Hex())
	}
	fmt.Printf("Transaction mined in block %d\n", receipt.BlockNumber.Uint64())

	// 6. Get the updated value from the contract.
	updatedValue, err := instance.Get(nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Updated value:", updatedValue)

	// 7. Call the `add` function of the contract.
	addValue := big.NewInt(10)
	fmt.Printf("Adding %d to value\n", addValue)

	// Pack the arguments for the `add` function. The `add` function likely
	// takes one uint256 argument.
	txData, err = contractABI.Pack("add", addValue) // "add" is the name of the function in the contract.
	if err != nil {
		log.Fatal("Error packing arguments for 'add':", err)
	}

	// Get transaction authorization options for the `add` transaction.
	authAdd, err := getTransactionAuthorizer(client, contractAddress, txData)
	if err != nil {
		log.Fatal(err)
	}

	// Send the transaction to call the `add` function.
	txAdd, err := instance.Add(authAdd, addValue)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Add transaction hash: %s\n", txAdd.Hash().Hex())

	// Wait for the `add` transaction to be mined.
	receiptAdd, err := bind.WaitMined(context.Background(), client, txAdd)
	if err != nil {
		log.Fatalf("Transaction %s mining failed: %v", txAdd.Hash().Hex(), err)
	}
	if receiptAdd.Status == types.ReceiptStatusFailed {
		log.Fatalf("Transaction %s failed", txAdd.Hash().Hex())
	}

	// Get the value after calling the `add` function.
	newValueAfterAdd, err := instance.Get(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New Value After Add:", newValueAfterAdd)
}

// getTransactionAuthorizer creates a `bind.TransactOptions` struct
// for signing and submitting transactions. It reads the private key
// from the environment.
func getTransactionAuthorizer(client *jumboclient.Client, contractAddress common.Address, txData []byte) (*bind.TransactOpts, error) {
	privateKeyHex := os.Getenv("PRIVATE_KEY") // The sender's private key in hexadecimal format.
	if privateKeyHex == "" {
		return nil, fmt.Errorf("PRIVATE_KEY environment variable not set")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex) // Convert the hex string private key to an ECDSA private key.
	if err != nil {
		return nil, err
	}

	// Get the sender's Ethereum address from the private key's public key.
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	// Get the nonce (transaction count) for the sender's address. This prevents replay attacks.
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
		return nil, err
	}

	// Get the chain ID of the Jumbochain network. This is used for EIP-155 signing.
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, err
	}

	// Create a new `bind.TransactOptions` struct. This struct holds
	// all the necessary information for signing and sending a transaction.
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, err
	}
	auth.Nonce = big.NewInt(int64(nonce)) // Set the nonce for the transaction.
	auth.Value = big.NewInt(0)            // Amount of Ether to send with the transaction (in wei). Set to 0 for contract calls.

	// Estimate the gas required for the transaction. This is a more accurate way
	// to set the gas limit and avoid "out of gas" errors.
	gas, err := client.EstimateGas(context.Background(), jumbochain.CallMsg{
		From: auth.From,
		To:   &contractAddress,
		Data: txData, // The encoded function call data.
	})
	if err != nil {
		log.Println("Warning: Error estimating gas:", err)
		// Fallback to a higher default gas limit if estimation fails.
		auth.GasLimit = uint64(500000)
	} else {
		fmt.Println("Estimated gas:", gas)
		auth.GasLimit = gas + 20000 // Add a buffer to the estimated gas.
	}

	return auth, nil
}
