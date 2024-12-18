package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
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

	var numberOfRecordsStr string
	var numOfRecToCreate int
	if len(os.Args) > 1 {
		numberOfRecordsStr = os.Args[1]
		numOfRecs, err := strconv.ParseInt(numberOfRecordsStr, 10, 0)
		if err != nil {
			numOfRecToCreate = 1
		} else {
			numOfRecToCreate = int(numOfRecs)
		}
	} else {
		numOfRecToCreate = 1
	}

	var wg sync.WaitGroup
	for i := 0; i < numOfRecToCreate; i++ {
		wg.Add(1)
		go sendWrite(&wg, i)
	}

	wg.Wait()

}

func sendWrite(wg *sync.WaitGroup, routineNumber int) {
	defer wg.Done()
	client := http.Client{Timeout: time.Duration(20) * time.Second}

	randomAccount := rand.Intn(100000)
	randomDollar := rand.Intn(100)
	randomCents := rand.Intn(100)
	randomPositive := rand.Intn(100)

	var action string

	if randomPositive <= 30 {
		action = "deposit"
	} else {
		action = "withdraw"
	}

	amount := fmt.Sprintf("%d.%02d", randomDollar, randomCents)
	var amountFloat, _ = strconv.ParseFloat(amount, 64)

	transaction := Transaction{
		Account: strconv.Itoa(randomAccount),
		Action:  action,
		Amount:  amountFloat,
	}

	content, _ := json.Marshal(transaction)

	bodyReader := bytes.NewReader(content)

	client.Post(url, contentType, bodyReader)

	fmt.Printf("Finished routine %d\n", routineNumber)
}
