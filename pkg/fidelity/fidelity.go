/*
Provides an interface for fetching information from Fidelity.
*/
package fidelity

type FidelityClient interface {
	// Returns the current balance of the Fidelity account
	GetBalance() (float64, error)
}
