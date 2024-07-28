package api

import (
	"encoding/json"
	"time"

	"github.com/go-resty/resty/v2"
)

const baseURL = "https://vinvoice.viettel.vn/api"

type AuthClient struct {
	client  *resty.Client
	token   string
	tokenAt time.Time
}

type InvoiceClient struct {
	client     *resty.Client
	authClient *AuthClient
}

func NewAuthClient() *AuthClient {
	client := resty.New()

	return &AuthClient{client: client}
}

func (c *AuthClient) GetToken(authObj map[string]string) (string, error) {
	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "en-US,en;q=0.9,vi;q=0.8").
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		SetHeader("Host", "vinvoice.viettel.vn").
		SetHeader("Referer", "https://vinvoice.viettel.vn/account/login").
		SetBody(authObj).Post("https://vinvoice.viettel.vn/api/auth/login")
	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", resp.Error().(error)
	}

	var result map[string]interface{}
	err = json.Unmarshal(resp.Body(), &result)
	if err != nil {
		return "", err
	}

	token, ok := result["token"].(string)
	if !ok {
		return "", err
	}

	c.token = token
	c.tokenAt = time.Now()

	return token, nil
}

func (c *AuthClient) IsTokenExpired() bool {
	return time.Since(c.tokenAt).Seconds() >= 250
}

func NewApiClient(baseURL string, authObj map[string]string) (*InvoiceClient, error) {
	authClient := NewAuthClient()
	token, err := authClient.GetToken(authObj)
	if err != nil {
		return nil, err
	}

	client := resty.New().SetBaseURL(baseURL).SetAuthToken(token)

	return &InvoiceClient{client: client, authClient: authClient}, nil
}
