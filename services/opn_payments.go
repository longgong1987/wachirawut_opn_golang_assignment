package services

import (
	"time"

	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type OpnClient struct {
	client *omise.Client
}

type CreateTokenPayload struct {
	Name            string
	Number          string
	ExpirationMonth int
	ExpirationYear  int
	SecurityCode    string
}

type CreateChargePayload struct {
	Amount    int64
	Currency  string
	CardToken string
}

func NewOpn(client *omise.Client) *OpnClient {
	return &OpnClient{
		client: client,
	}
}

func (config *OpnClient) CreateToken(payload *CreateTokenPayload) (*omise.Card, error) {

	result := &omise.Card{}

	err := config.client.Do(result, &operations.CreateToken{
		Name:            payload.Name,
		Number:          payload.Number,
		ExpirationMonth: time.Month(payload.ExpirationMonth),
		ExpirationYear:  payload.ExpirationYear,
		SecurityCode:    payload.SecurityCode,
	})

	if err != nil {
		return &omise.Card{}, err
	}

	return result, nil
}

func (config *OpnClient) CreateCharge(payload *CreateChargePayload) (*omise.Charge, error) {

	result := &omise.Charge{}
	err := config.client.Do(result, &operations.CreateCharge{
		Amount:   payload.Amount,
		Currency: payload.Currency,
		Card:     payload.CardToken,
	})

	if err != nil {
		return &omise.Charge{}, err
	}

	return result, nil
}
