//go:build browser

package fidelity

import (
	"runtime"
	"testing"
	"time"

	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/chromedp"
	"kevincao.dev/fidelity2ynab/pkg/secrets"
)

func beforeEach(t testing.TB) (client *fidelityBrowserClient, teardown func()) {
	t.Helper()
	creds, err := secrets.GetSecrets()
	if err != nil {
		t.Fatalf("Failed to get secrets: %s", err)
	}

	opts := []cu.Option{cu.WithTimeout(2 * time.Minute)}
	// Only use headless option if on Linux as the undetected chromedp
	// package does not support headless mode on other OSes.
	if runtime.GOOS == "linux" {
		opts = append(opts, cu.WithHeadless())
	}

	client, err = NewFidelityBrowserClient(
		creds.FidelityUsername,
		creds.FidelityPassword,
		creds.FidelityTotpSecret,
		opts...,
	)
	if err != nil {
		t.Fatalf("Failed to initialize Fidelity Browser client: %s", err)
	}

	return client, func() { client.Close() }
}

// Tests that the browser can log into the Fidelity dashboard.
func TestBrowserLogin(t *testing.T) {
	client, teardown := beforeEach(t)
	defer teardown()
	client.login()

	if err := chromedp.Run(
		client.cuCtx,
		chromedp.WaitReady(".total-balance-value", chromedp.ByQuery),
	); err != nil {
		t.Errorf("Failed to log into Fidelity dashboard: %s", err)
	}
}

// Tests that the browser can get the balance from the Fidelity dashboard.
func TestBrowserGetBalance(t *testing.T) {
	client, teardown := beforeEach(t)
	defer teardown()

	_, err := client.GetBalance()
	if err != nil {
		t.Errorf("Failed to get Fidelity balance: %s", err)
	}
}
