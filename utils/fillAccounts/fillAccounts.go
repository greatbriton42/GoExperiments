package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

const accountFile = "../../server/db/accounts.db"

func main() {
	fmt.Println("Running")
	f, err := os.OpenFile(accountFile, os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Error Opening %s", err)
	}
	defer f.Close()

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

	fmt.Printf("Creating %d Records\n", numOfRecToCreate)
	for i := 0; i < int(numOfRecToCreate); i++ {
		randomAccount := rand.Intn(100000)
		randomDollar := rand.Intn(100)
		randomCents := rand.Intn(100)
		randomPositive := rand.Intn(100)

		var symbol string

		if randomPositive <= 80 {
			symbol = ""
		} else {
			symbol = "-"
		}

		record := fmt.Sprintf("%05d:%s%d.%02d", randomAccount, symbol, randomDollar, randomCents)
		//We will pad to give each line some space to be edited. 5 for account + ':' + 9 for account amount
		paddedRecord := fmt.Sprintf("%-15s\n", record)
		fmt.Printf("%s", paddedRecord)
		_, err = f.WriteString(paddedRecord)
		if err != nil {
			fmt.Printf("Error Writing %s", err)
		}
	}
	f.Sync()
	fmt.Println("\nFinished Writing to File")
}
