package twilio

import (
	"io"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type twilioWriter struct {
	client     *twilio.RestClient
	FromNumber string
	DestNumber string
}

var _ io.Writer = twilioWriter{}

// Creates a new TwilioWriter with the given account SID, API secret, and destination number
func NewTwilioWriter(
	accountSid string,
	apiSecret string,
	fromNumber string,
	destNumber string,
) *twilioWriter {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: apiSecret,
	})
	return &twilioWriter{
		client:     client,
		FromNumber: fromNumber,
		DestNumber: destNumber,
	}
}

// Sends a message to the destination number
func (s twilioWriter) Write(p []byte) (int, error) {
	content := string(p)
	_, err := s.client.Api.CreateMessage(&twilioApi.CreateMessageParams{
		To:   &s.DestNumber,
		From: &s.FromNumber,
		Body: &content,
	})
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
