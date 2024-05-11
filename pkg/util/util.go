/*
Provides utility functions for ease of use.
*/
package util

// Useful function to get address of literal values
func Addr[T any](v T) *T {
	return &v
}
