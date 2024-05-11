/*
Package secrets is responsible for reading the secrets file and providing access
to the secrets to other packages.
*/
package secrets

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type secrets struct {
	FidelityTotpSecret      string `yaml:"fidelity_totp_secret"`
	FidelityUsername        string `yaml:"fidelity_username"`
	FidelityPassword        string `yaml:"fidelity_password"`
	YnabApiKey              string `yaml:"ynab_api_key"`
	YnabFidelityAccountName string `yaml:"ynab_fidelity_account_name"`
	TwilioAccountSid        string `yaml:"twilio_account_sid"`
	TwilioApiSecret         string `yaml:"twilio_api_secret"`
	TwilioNumber            string `yaml:"twilio_number"`
	TwilioToNumber          string `yaml:"twilio_to_number"`
}

var creds secrets

func init() {
	// Attempt to read secrets assuming in docker container first
	data, err := os.ReadFile("/run/secrets/secrets")
	if err != nil {
		data, err = os.ReadFile("secrets.yaml")
		if err != nil {
			return
		}
	}

	if err = yaml.Unmarshal(data, &creds); err != nil {
		panic(err)
	}
}

func GetSecrets() (secrets, error) {
	if creds == (secrets{}) {
		return creds, errors.New("secrets not initialized")
	}
	return creds, nil
}
