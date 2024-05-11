/*
Provides helper functions for syncing balances with brunomvsouza/ynab
*/
package ynab

import (
	"errors"
	"fmt"
	"time"

	brynab "github.com/brunomvsouza/ynab.go"
	"github.com/brunomvsouza/ynab.go/api"
	"github.com/brunomvsouza/ynab.go/api/account"
	"github.com/brunomvsouza/ynab.go/api/transaction"
	"kevincao.dev/fidelity2ynab/pkg/log"
	"kevincao.dev/fidelity2ynab/pkg/util"
)

type ynabClient struct {
	brynab.ClientServicer
}

// NewClient creates a new YNAB client with the given access token.
func NewClient(accessToken string) *ynabClient {
	return &ynabClient{
		brynab.NewClient(accessToken),
	}
}

// Searches all YNAB budgets for an account with the given name. Returns
// the first account found with the given name.
func (c *ynabClient) GetAccountWithName(name string) (
	account *account.Account, budgetID string, err error,
) {
	budgets, err := c.ClientServicer.Budget().GetBudgets()
	if err != nil {
		return
	}
	for _, budget := range budgets {
		accounts, err := c.ClientServicer.Account().GetAccounts(budget.ID, nil)
		if err != nil {
			return account, budgetID, err
		}
		for _, a := range accounts.Accounts {
			if a.Name == name {
				budgetID = budget.ID
				account = a
				break
			}
		}
	}
	if account == nil {
		return account, budgetID, errors.New("Could not find YNAB account with name " + name)
	}
	return
}

// Updates the balance of the given account to the given balance by adding
// a transaction to the account.
func (c *ynabClient) UpdateAccountBalance(
	budgetID string,
	account *account.Account,
	balance float64,
) error {
	currentBalance := float64(account.Balance) / 1000
	difference := balance - currentBalance
	log.Info(fmt.Sprintf(
		"Updating balance from %.2f to %.2f by adding transaction %.2f",
		currentBalance, balance, difference,
	))
	_, err := c.ClientServicer.Transaction().CreateTransaction(
		budgetID,
		transaction.PayloadTransaction{
			AccountID: account.ID,
			Date:      api.Date{Time: time.Now()},
			Amount:    int64(difference * 1000),
			Cleared:   transaction.ClearingStatusCleared,
			Approved:  true,
			PayeeName: util.Addr("[F2Y] Automatic Balance Adjustment"),
			Memo:      util.Addr("Automatically updated by Fidelity2Ynab"),
		},
	)
	return err
}
