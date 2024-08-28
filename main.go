package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"time"

	vex "github.com/genesisblockid/vex-go"
	"syscall"
)

func main() {
	// Initialize the VEX API
	api := vex.New("https://v2.vexascan.com:2096")
	ctx := context.Background()

	// List of private keys to add to the key bag
	privateKeys := []string{"A", "B", "C"} // Replace with your actual private keys

	// Set up the key bag and signer
	keyBag := &vex.KeyBag{}
	for _, key := range privateKeys {
		keyBag.Add(key)
	}
	api.SetSigner(keyBag)

	// List of accounts to claim rewards for
	accounts := []string{"a", "b", "c"}

	// Create a channel to listen for interrupt signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// Infinite loop to claim rewards periodically
	go func() {
		for {
			select {
			case <-signalChan:
				fmt.Println("Shutdown signal received. Exiting...")
				os.Exit(0)
			default:
				for _, account := range accounts {
					claimReward(api, ctx, account)
				}
				fmt.Println("Waiting for 60 seconds before the next claim...")
				time.Sleep(60 * time.Second)
			}
		}
	}()

	// Block main thread to keep application running
	select {}
}

// Function to claim rewards for block producers
func claimReward(api *vex.API, ctx context.Context, account string) { // Updated to accept `account` argument
	fmt.Printf("Trying to claim reward for account: %s\n", account)

	// Define the action data struct with correct capitalization
	actionData := struct {
		Owner vex.AccountName `json:"owner"`
	}{
		Owner: vex.AccountName(account), // Convert string to vex.AccountName
	}

	// Define the action to claim rewards
	action := vex.Action{
		Account:       vex.AccountName("vexcore"), // Convert string to vex.AccountName
		Name:          vex.ActionName("claimrewards"),
		Authorization: []vex.PermissionLevel{{Actor: vex.AccountName(account), Permission: vex.PermissionName("active")}},
		ActionData:    vex.NewActionData(actionData),
	}

	// Set transaction options
	txOpts := &vex.TxOptions{}
	if err := txOpts.FillFromChain(ctx, api); err != nil {
		fmt.Printf("Error filling transaction options: %v\n", err)
		return
	}

	// Create and sign the transaction
	trx := vex.NewTransaction([]*vex.Action{&action}, txOpts)
	signedTx, packedTx, err := api.SignTransaction(ctx, trx, txOpts.ChainID, vex.CompressionNone)
	if err != nil {
		fmt.Printf("Error signing transaction: %v\n", err)
		return
	}

	// Print signed transaction for debugging
	content, err := json.MarshalIndent(signedTx, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling transaction: %v\n", err)
		return
	}
	fmt.Println("Signed Transaction")
	fmt.Println(string(content))

	// Push the transaction to the blockchain
	response, err := api.PushTransaction(ctx, packedTx)
	if err != nil {
		fmt.Printf("Error pushing transaction: %v\n", err)
		return
	}

	// Print success message
	fmt.Printf("Transaction [%s] submitted successfully for account: %s\n", response.Processed.ID, account)
}
// Save PID to a file
func writePID(pidFile string) {
	pid := os.Getpid()
	err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644)
	if err != nil {
		fmt.Printf("Error writing PID file: %v\n", err)
	}
}

// Cleanup function to remove PID file
func cleanup(pidFile string) {
	err := os.Remove(pidFile)
	if err != nil {
		fmt.Printf("Error removing PID file: %v\n", err)
	}
}
