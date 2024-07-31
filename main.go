package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/nt2311-vn/invoice-microservices/api"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	authClient := api.NewAuthClient()
	c, err := api.NewInvoiceClient(authClient)
	if err != nil {
		log.Fatalf("Error creating invoice client: %v\n", err)
	}

	resp, err := c.FetchInvoices()
	if err != nil {
		log.Fatalf("Error fetching invoices: %v\n", err)
	}

	jsonData, err := json.MarshalIndent(resp, " ", "  ")
	fmt.Println(string(jsonData))
}
