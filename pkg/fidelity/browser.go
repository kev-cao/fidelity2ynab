package fidelity

import (
	"bytes"
	"context"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	cu "github.com/Davincible/chromedp-undetected"
	"github.com/chromedp/chromedp"
	"github.com/pquerna/otp/totp"
	"kevincao.dev/fidelity2ynab/pkg/log"
)

const (
	DELAY = 500 * time.Millisecond
)

type fidelityBrowserClient struct {
	username    string
	password    string
	totp_secret string
	cuCtx       context.Context
	cuCtxCancel context.CancelFunc
}

var _ FidelityClient = fidelityBrowserClient{}

// NewFidelityBrowserClient creates a new FidelityBrowserClient
// If no options are provided, the default options are a 1 minute timeout and headless mode.
func NewFidelityBrowserClient(username, password, totp_secret string, opts ...cu.Option) (*fidelityBrowserClient, error) {
	client := fidelityBrowserClient{
		username:    username,
		password:    password,
		totp_secret: totp_secret,
	}
	if len(opts) == 0 {
		// Delete WithHeadless and run on local machine to view browser in realtime
		opts = []cu.Option{cu.WithTimeout(1 * time.Minute), cu.WithHeadless()}
	}
	if err := client.initializeBrowserContext(opts...); err != nil {
		return nil, err
	}
	return &client, nil
}

func (c *fidelityBrowserClient) initializeBrowserContext(opts ...cu.Option) error {
	ctx, cancel, err := cu.New(cu.NewConfig(opts...))
	if err != nil {
		return errors.New("Failed to create undetected chromedp context: " + err.Error())
	}
	c.cuCtx = ctx
	c.cuCtxCancel = cancel
	return nil
}

// Close closes the browser client context
func (c fidelityBrowserClient) Close() {
	c.cuCtxCancel()
}

func (c fidelityBrowserClient) login() error {
	if err := chromedp.Run(
		c.cuCtx, chromedp.Navigate("https://digital.fidelity.com/prgw/digital/login/full-page"),
	); err != nil {
		return errors.New("Failed to navigate to Fidelity login page: " + err.Error())
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

func (c fidelityBrowserClient) submitCredentials() error {
	if err := chromedp.Run(c.cuCtx, chromedp.SendKeys("#dom-username-input", c.username, chromedp.ByQuery)); err != nil {
		return errors.New("Could not find username input element: " + err.Error())
	}
	log.Debug("Found username input element")
	time.Sleep(DELAY)

	if err := chromedp.Run(c.cuCtx, chromedp.SendKeys("#dom-pswd-input", c.password, chromedp.ByQuery)); err != nil {
		return errors.New("Could not find password input element: " + err.Error())
	}
	log.Debug("Found password input element")
	time.Sleep(DELAY)

	if err := chromedp.Run(c.cuCtx, chromedp.Click("#dom-login-button", chromedp.ByQuery)); err != nil {
		return errors.New("Could not find login button: " + err.Error())
	}
	log.Debug("Found login button")
	return nil
}

func (c fidelityBrowserClient) submitTotp() error {
	code, err := totp.GenerateCode(c.totp_secret, time.Now())
	if err != nil {
		return errors.New("Failed to generate TOTP code: " + err.Error())
	}
	if err := chromedp.Run(c.cuCtx, chromedp.SendKeys("#dom-svip-security-code-input", code, chromedp.ByQuery)); err != nil {
		return errors.New("Could not find TOTP element: " + err.Error())
	}
	log.Debug("Found TOTP input element")
	time.Sleep(DELAY)

	if err := chromedp.Run(c.cuCtx, chromedp.Click("#dom-svip-code-submit-button", chromedp.ByQuery)); err != nil {
		return errors.New("Could not find TOTP submit button: " + err.Error())
	}
	log.Debug("Found TOTP submit button")

	return nil
}

func (c fidelityBrowserClient) GetBalance() (float64, error) {
	if err := c.login(); err != nil {
		return 0, err
	}

	content := bytes.Buffer{}
	if err := chromedp.Run(
		c.cuCtx,
		chromedp.Dump(".balance-total-value", &content, chromedp.ByQuery),
	); err != nil {
		return 0, errors.New("Failed to read balance element: " + err.Error())
	}

	balancePattern, _ := regexp.Compile("\\$[0-9,]+\\.[0-9]+")
	balanceString := balancePattern.FindString(content.String())
	if balanceString == "" {
		return 0, errors.New("Failed to find balance in element")
	}
	balanceString = balanceString[1:] // Remove dollar sign
	return strconv.ParseFloat(strings.ReplaceAll(balanceString, ",", ""), 64)
}
