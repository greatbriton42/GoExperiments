package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const url = "http://localhost:8080/write"
const contentType = "application/json"

type Transaction struct {
	Account string
	Action  string
	Amount  float64
}

type Account struct {
	Account string
	Amount  string
}

func main() {
	client := http.DefaultClient

	transaction := Transaction{
		Account: "12345",
		Action:  "deposit",
		Amount:  2.50,
	}

	content, _ := json.Marshal(transaction)

	bodyReader := bytes.NewReader(content)

	client.Post(url, contentType, bodyReader)
}
