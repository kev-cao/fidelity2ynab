/*
This package is responsible for reading the secrets file and providing access
to the secrets to other packages.
*/
package secrets

import (
	"errors"
	"os"

	"gopkg.in/yaml.v3"
)

type secrets struct {
	FidelityTotpSecret string `yaml:"fidelity_totp_secret"`
	FidelityUsername   string `yaml:"fidelity_username"`
	FidelityPassword   string `yaml:"fidelity_password"`
	YnabApiKey         string `yaml:"ynab_api_key"`
	TwilioAccountSid   string `yaml:"twilio_account_sid"`
	TwilioApiSecret    string `yaml:"twilio_api_secret"`
	TwilioNumber       string `yaml:"twilio_number"`
}

var creds secrets

func init() {
	data, err := os.ReadFile("/run/secrets/secrets")
	if err != nil {
		return
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
