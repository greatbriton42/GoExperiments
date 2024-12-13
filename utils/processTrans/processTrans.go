package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

const accountFile = "../../server/db/accounts.db"
const ledgerFile = "../../server/db/ledger.db"

type Transaction struct {
	Account string
	Name    string
	Action  string
	Amount  float64
}

type Account struct {
	Account string
	Amount  float64
}

func main() {
	fAccounts, err := os.OpenFile(accountFile, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Error Opening %s", err)
	}
	defer fAccounts.Close()

	fLedger, err := os.OpenFile(ledgerFile, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("Error Opening %s", err)
	}
	defer fLedger.Close()

	var transaction Transaction

	scanner := bufio.NewScanner(fLedger)
	scanner.Split(ScanTransactions)
	for scanner.Scan() {
		transactionString := scanner.Text()

		json.Unmarshal([]byte(transactionString), &transaction)
		fmt.Printf("Processing Transaction: %s\n", transactionString)
		updateAccount(transaction, fAccounts)
	}

	//Clear all transactions
	//TODO: leave errored transaction for further processing
	fLedger.Truncate(0)
}

func ScanTransactions(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if r != '|' {
			break
		}
	}
	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if r == '|' {
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}

func updateAccount(transaction Transaction, f *os.File) {
	var account Account

	f.Seek(0, 0)
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	byteNumber := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, transaction.Account+":") {
			accountInfo := strings.Split(line, ":")
			account.Account = strings.TrimSpace(accountInfo[0])
			account.Amount, _ = strconv.ParseFloat(strings.TrimSpace(accountInfo[1]), 64)
			break
		} else {
			byteNumber += len(line) + 1
		}
	}

	if account.Amount > 0.0 && len(account.Account) > 0 {
		if strings.ToLower(transaction.Action) == "withdraw" {
			account.Amount -= transaction.Amount
		} else if strings.ToLower(transaction.Action) == "deposit" {
			account.Amount += transaction.Amount
		} else {
			fmt.Printf("Not a valid action (%s)", transaction.Action)
			return
		}

		f.Seek(int64(byteNumber), 0)
		record := fmt.Sprintf("%s:%.2f", account.Account, account.Amount)
		paddedRecord := fmt.Sprintf("%-15s\n", record)
		fmt.Printf("Writing to byte offset: %d with %s", byteNumber, paddedRecord)
		f.WriteString(paddedRecord)
		f.Sync()
	}
}
