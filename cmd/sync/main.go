/*
Entrypoint for sync command. Will fetch the current balance from Fidelity and update it in YNAB.
*/
package main

import (
	"fmt"

	"kevincao.dev/fidelity2ynab/pkg/fidelity"
	"kevincao.dev/fidelity2ynab/pkg/log"
	"kevincao.dev/fidelity2ynab/pkg/secrets"
	"kevincao.dev/fidelity2ynab/pkg/ynab"
)

func main() {
	secrets, err := secrets.GetSecrets()
	if err != nil {
		log.Fatal("Failed to get secrets: %s", err)
	}
	fidelityClient, err := fidelity.NewFidelityBrowserClient(
		secrets.FidelityUsername,
		secrets.FidelityPassword,
		secrets.FidelityTotpSecret,
	)
	if err != nil {
		log.Fatal("Failed to initialize Fidelity Browser client: " + err.Error())
	}
	defer fidelityClient.Close()
	fidelityBalance, err := fidelityClient.GetBalance()
	if err != nil {
		log.Fatal("Failed to get Fidelity Balance: " + err.Error())
	}
	log.Debug(fmt.Sprintf("Fidelity balance is %f", fidelityBalance))

	ynabClient := ynab.NewClient(secrets.YnabApiKey)
	ynabFidelity, budget, err := ynabClient.GetAccountWithName(secrets.YnabFidelityAccountName)
	if err != nil {
		log.Fatal("Failed to find YNAB account with name " + secrets.YnabFidelityAccountName)
	}
	ynabClient.UpdateAccountBalance(budget, ynabFidelity, fidelityBalance)
	if err != nil {
		log.Fatal("Failed to update YNAB account: " + err.Error())
	}
}
