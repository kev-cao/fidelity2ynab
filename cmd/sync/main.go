/*
Entrypoint for sync command. Will fetch the current balance from Fidelity and update it in YNAB.
*/
package main

import (
	"fmt"
	"os"

	"kevincao.dev/fidelity2ynab/pkg/fidelity"
	"kevincao.dev/fidelity2ynab/pkg/log"
	"kevincao.dev/fidelity2ynab/pkg/secrets"
)

func main() {
	secrets, err := secrets.GetSecrets()
	if err != nil {
		log.Error("Failed to get secrets: %s", err)
		os.Exit(1)
	}
	client, err := fidelity.NewFidelityBrowserClient(
		secrets.FidelityUsername,
		secrets.FidelityPassword,
		secrets.FidelityTotpSecret,
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to initialize Fidelity Browser client: %s", err))
		os.Exit(1)
	}
	defer client.Close()
	balance, err := client.GetBalance()
	if err != nil {
		log.Error(fmt.Sprintf("Failed to get Fidelity Balance: %s", err))
		os.Exit(1)
	}
	fmt.Println(balance)
}
