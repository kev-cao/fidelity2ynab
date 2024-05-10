package fidelity

import (
	"github.com/tebeka/selenium"
)

type FidelitySeleniumClient struct {
	email       string
	password    string
	totp_secret string
	service     selenium.Service
	webdriver   selenium.WebDriver
}

func (c *FidelitySeleniumClient) initializeSelenium() {
	selenium.NewChromeDriverService("/app/bin/chromedriver", 4444)
}

func (c FidelitySeleniumClient) GetBalance() float64 {
	return 0
}
