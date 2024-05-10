package fidelity

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/chromedp"
	"github.com/pquerna/otp/totp"
	"kevincao.dev/fidelity2ynab/pkg/log"
)

type fidelitySeleniumClient struct {
	username    string
	password    string
	totp_secret string
	cuCtx       context.Context
	cuCtxCancel context.CancelFunc
}

func NewFidelityBrowserClient(username, password, totp_secret string) (*fidelitySeleniumClient, error) {
	client := fidelitySeleniumClient{
		username:    username,
		password:    password,
		totp_secret: totp_secret,
	}
	if err := client.initializeBrowser(); err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *fidelitySeleniumClient) initializeBrowser() error {
	ctx, cancel, err := cu.New(cu.NewConfig(
		cu.WithHeadless(),
		cu.WithTimeout(1*time.Minute),
	))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to create undetected chromedp context: %s", err))
	}
	c.cuCtx = ctx
	c.cuCtxCancel = cancel
	return nil
}

func (c fidelitySeleniumClient) Close() {
	c.cuCtxCancel()
}

func (c fidelitySeleniumClient) login() error {
	if err := chromedp.Run(
		c.cuCtx, chromedp.Navigate("https://digital.fidelity.com/prgw/digital/login/full-page"),
	); err != nil {
		return errors.New(fmt.Sprintf("Failed to navigate to Fidelity login page: %s", err))
	}
	log.Debug("Navigated to login page")
	if err := c.submitCredentials(); err != nil {
		return errors.New(err.Error())
	}
	if err := c.submitTotp(); err != nil {
		return errors.New(err.Error())
	}
	return nil
}

func (c fidelitySeleniumClient) submitCredentials() error {
	if err := chromedp.Run(c.cuCtx, chromedp.SendKeys("#dom-username-input", c.username, chromedp.ByQuery)); err != nil {
		return errors.New(fmt.Sprintf("Could not find username input element: %s", err))
	}
	log.Debug("Found username input element")
	time.Sleep(1 * time.Second)

	if err := chromedp.Run(c.cuCtx, chromedp.SendKeys("#dom-pswd-input", c.password, chromedp.ByQuery)); err != nil {
		return errors.New(fmt.Sprintf("Could not find password input element: %s", err))
	}
	log.Debug("Found password input element")
	time.Sleep(1 * time.Second)

	if err := chromedp.Run(c.cuCtx, chromedp.Click("#dom-login-button", chromedp.ByQuery)); err != nil {
		return errors.New(fmt.Sprintf("Could not find login button: %s", err))
	}
	log.Debug("Found login button")
	return nil
}

func (c fidelitySeleniumClient) submitTotp() error {
	code, err := totp.GenerateCode(c.totp_secret, time.Now())
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to generate TOTP code: %s", err))
	}
	if err := chromedp.Run(c.cuCtx, chromedp.SendKeys("#dom-svip-security-code-input", code, chromedp.ByQuery)); err != nil {
		return errors.New(fmt.Sprintf("Could not find TOTP element: %s", err))
	}
	log.Debug("Found TOTP input element")
	time.Sleep(1 * time.Second)

	if err := chromedp.Run(c.cuCtx, chromedp.Click("#dom-svip-code-submit-button", chromedp.ByQuery)); err != nil {
		return errors.New(fmt.Sprintf("Could not find TOTP submit button: %s", err))
	}
	log.Debug("Found TOTP submit button")

	return nil
}

func (c fidelitySeleniumClient) GetBalance() (float64, error) {
	if err := c.login(); err != nil {
		return 0, err
	}

	content := bytes.Buffer{}
	if err := chromedp.Run(
		c.cuCtx,
		chromedp.Dump(".total-balance-value", &content, chromedp.ByQuery),
	); err != nil {
		return 0, errors.New(fmt.Sprintf("Failed to read balance element: %s", err))
	}

	balancePattern, _ := regexp.Compile("\\$[0-9,]+\\.[0-9]+")
	balanceString := balancePattern.FindString(content.String())
	if balanceString == "" {
		return 0, errors.New("Failed to find balance in element")
	}
	balanceString = balanceString[1:] // Remove dollar sign
	return strconv.ParseFloat(strings.ReplaceAll(balanceString, ",", ""), 64)
}
