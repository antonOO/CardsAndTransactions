package main

import (
	"encoding/json"
	"fmt"
	transactions "github.com/antonOO/CardsAndTransactions/transactions"
	"log"
	"net/http"
)

type CardDto struct {
	id int
}

func main() {

	// Server supports only a single client as there is no budget :D
	// Nor part of the assignment. Serves for testing purposes by the reviewers
	client := transactions.Client{}

	http.HandleFunc("/add-card", func(w http.ResponseWriter, r *http.Request) {
		var cardDto CardDto
		err := json.NewDecoder(r.Body).Decode(&cardDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	})

	http.HandleFunc("/activate", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	http.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		client.GetActiveCards()
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
