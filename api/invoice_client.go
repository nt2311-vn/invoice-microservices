package api

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
)

const baseURL = "https://vinvoice.viettel.vn/api"

type TokenResponse struct {
	AccessToken    string `json:"access_token"`
	ExpiresIn      int    `json:"expires_in"`
	Iat            int    `json:"iat"`
	InvoiceCluster string `json:"invoice_cluster"`
	Jti            string `json:"jti"`
	RefreshToken   string `json:"refresh_token"`
	Scope          string `json:"scope"`
	TokenType      string `json:"token_type"`
	Type           int    `json:"type"`
}

type AuthClient struct {
	client  *resty.Client
	token   *TokenResponse
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

func (c *AuthClient) GetToken() (*TokenResponse, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	authObj := map[string]string{
		"username": os.Getenv("INVOICE_CLIENT_USERNAME"),
		"password": os.Getenv("INVOICE_CLIENT_PASSWORD"),
	}

	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "en-US,en;q=0.9,vi;q=0.8").
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		SetHeader("Host", "vinvoice.viettel.vn").
		SetHeader("Referer", "https://vinvoice.viettel.vn/account/login").
		SetBody(authObj).Post("https://vinvoice.viettel.vn/api/auth/login")
	if err != nil {
		return nil, err
	}

	if resp.IsError() {
		return nil, fmt.Errorf("error: %v", resp.Error())
	}

	token := &TokenResponse{}
	err = json.Unmarshal(resp.Body(), token)
	if err != nil {
		return nil, err
	}

	c.token = token
	c.tokenAt = time.Now()

	return token, nil
}

func (c *AuthClient) IsTokenExpired() bool {
	return time.Since(c.tokenAt).Seconds() >= 250
}

func NewApiClient() (*InvoiceClient, error) {
	authClient := NewAuthClient()
	token, err := authClient.GetToken()
	if err != nil {
		return nil, err
	}

	client := resty.New().SetBaseURL(baseURL).SetAuthToken(token.AccessToken)

	return &InvoiceClient{client: client, authClient: authClient}, nil
}
