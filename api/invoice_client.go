package api

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	"github.com/nt2311-vn/invoice-microservices/models"
)

const baseURL = "https://vinvoice.viettel.vn/api/"

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
	authObj map[string]string
}

func NewAuthClient() *AuthClient {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return nil
	}

	client := resty.New()
	return &AuthClient{
		client: client,
		authObj: map[string]string{
			"username": os.Getenv("INVOICE_CLIENT_USERNAME"),
			"password": os.Getenv("INVOICE_CLIENT_PASSWORD"),
		},
	}
}

func (c *AuthClient) GetToken() (*TokenResponse, error) {
	resp, err := c.client.R().
		SetHeader("Accept", "application/json").
		SetHeader("Accept-Language", "en-US,en;q=0.9,vi;q=0.8").
		SetHeader("Content-Type", "application/json;charset=UTF-8").
		SetHeader("Host", "vinvoice.viettel.vn").
		SetHeader("Referer", "https://vinvoice.viettel.vn/account/login").
		SetBody(c.authObj).Post(baseURL + "auth/login")
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

type InvoiceClient struct {
	client     *resty.Client
	authClient *AuthClient
}

func NewInvoiceClient(authClient *AuthClient) (*InvoiceClient, error) {
	if authClient.IsTokenExpired() {
		_, err := authClient.GetToken()
		if err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
	}

	client := resty.New().SetAuthToken(authClient.token.AccessToken)
	return &InvoiceClient{client: client, authClient: authClient}, nil
}

func (c *InvoiceClient) FetchInvoices() ([]models.InvoiceResponse, error) {
	if c.authClient.IsTokenExpired() {
		token, err := c.authClient.GetToken()
		if err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}

		c.client.SetAuthToken(token.AccessToken)
	}

	// Set queries with bearer token
	page := 1

	resp, err := c.client.R().
		SetQueryParams(map[string]string{
			"adjustmentType.equals":          "1",
			"supplierId.equals":              "16087",
			"dataType.equals":                "0",
			"taxCode.equals":                 "0301482205",
			"sort":                           "issueDate,desc",
			"invoiceSeri.equals":             "C24MAA",
			"invoiceStatus.equals":           "1",
			"size":                           "10",
			"page":                           fmt.Sprintf("%d", page),
			"createdDate.greaterThanOrEqual": "2021-01-01T00:00:00Z",
			"createdDate.lessThanOrEqual":    "2021-12-31T23:59:59Z",
		}).SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.authClient.token.AccessToken)).Get(baseURL + "services/einvoiceapplication/api/invoice/search")
	if err != nil {
		return nil, fmt.Errorf("Cannot request invoice: %v\n", err)
	}

	invResps := []models.InvoiceResponse{}
	invResp := models.InvoiceResponse{}

	err = json.Unmarshal(resp.Body(), &invResp)
	if err != nil {
		return nil, fmt.Errorf("Cannot unmarshal json: %v\n", err)
	}
	fmt.Println(invResps)

	return invResps, nil
}
