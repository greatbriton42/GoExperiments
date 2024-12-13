package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

const ledgerFile = "ledger.db"
const accountFile = "accounts.db"

type Transaction struct {
	Account string
	Name    string
	Action  string
	Amount  float64
}

type Account struct {
	Account string
	Amount  string
}

func main() {
	mux := http.NewServeMux()

	//Set up routing
	mux.HandleFunc("POST /write", handleWrite)
	mux.HandleFunc("GET /read/{id}", handleRead)

	fmt.Println("Listening...")
	http.ListenAndServe(":8080", mux)

}

func handleWrite(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var transaction Transaction

	dec.Decode(&transaction)

	fmt.Println("Writing ->")
	fmt.Printf("\tname: %s\n", transaction.Name)
	fmt.Printf("\tamount: %f\n", transaction.Amount)
	fmt.Printf("\taction: %s\n", transaction.Action)

	f, err := os.OpenFile(ledgerFile, os.O_APPEND, 0644)
	checkFileOpenErr(err, w)
	defer f.Close()

	record, _ := json.Marshal(transaction)

	_, err = f.WriteString(string(record) + "|")
	if err != nil {
		w.WriteHeader(500)
		errorMes := fmt.Sprintf("Error Writing to file: %s", err)
		w.Write([]byte(errorMes))
	}
	f.Sync()
}
func handleRead(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	fmt.Printf("Reading: %s \n", idString)

	f, err := os.Open(accountFile)
	checkFileOpenErr(err, w)

	var account Account

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, idString+":") {
			accountInfo := strings.Split(line, ":")
			account.Account = accountInfo[0]
			account.Amount = accountInfo[1]
			break
		}
	}

	if account.Account == "" {
		w.WriteHeader(404)
	} else {
		response, _ := json.Marshal(account)
		w.Header().Set("content-type", "application/json")
		w.Write(response)
	}
}

func checkFileOpenErr(e error, w http.ResponseWriter) {
	if e != nil {
		w.WriteHeader(500)
		var errorMes string
		if os.IsNotExist(e) {
			errorMes = fmt.Sprintf("File does not exist: %s", e)
		} else if os.IsPermission(e) {
			errorMes = fmt.Sprintf("Permission denied: %s", e)
		} else {
			errorMes = fmt.Sprintf("Error opening file: %s", e)
		}
		w.Write([]byte(errorMes))
	}
}
