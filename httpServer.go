package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Transaction struct {
	Account int
	Name    string
	Action  string
	Amount  float32
}

func main() {
	mux := http.NewServeMux()

	//Set up routing
	mux.HandleFunc("POST /write", handleWrite)
	mux.HandleFunc("GET /read/{id}", handleRead)

	http.ListenAndServe(":8080", mux)

	fmt.Println("Listening...")
}

func handleWrite(w http.ResponseWriter, r *http.Request) {

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	var transaction Transaction

	dec.Decode(&transaction)

	fmt.Println("Writing ->")
	fmt.Printf("\tname: %s\n", transaction.Name)
	fmt.Printf("\tamount: %v\n", transaction.Amount)
	fmt.Printf("\taction: %s\n", transaction.Action)
}
func handleRead(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")

	fmt.Printf("Reading: %s \n", idString)
	transaction := Transaction{Name: "MyName", Account: 12345, Action: "Deposit", Amount: 10.50}
	response, _ := json.Marshal(transaction)
	w.Write(response)
}
