package main

import (
	"encoding/json"
	"fmt"
	transactions "github.com/antonOO/CardsAndTransactions/transactions"
	"log"
	"net/http"
)

type CardDto struct {
	Id int `json:"id"`
}

func main() {

	// Server supports only a single client as there is no budget :D
	// Nor part of the assignment. Serves for testing purposes by the reviewers
	client := transactions.NewClient()

	http.HandleFunc("/add-card", func(w http.ResponseWriter, r *http.Request) {
		var cardDto CardDto

		err := json.NewDecoder(r.Body).Decode(&cardDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		card, err := transactions.NewCard(cardDto.Id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		client.AddCard(card)
		fmt.Fprintf(w, "Added")
	})

	http.HandleFunc("/activate", func(w http.ResponseWriter, r *http.Request) {
		var cardDto CardDto
		err := json.NewDecoder(r.Body).Decode(&cardDto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		client.ReceiveTransaction(cardDto.Id)
		fmt.Fprintf(w, "Activated")
	})

	http.HandleFunc("/cards", func(w http.ResponseWriter, r *http.Request) {
		cards := client.GetActiveCards()
		for _, card := range cards {
			fmt.Fprintf(w, "ID - %v, Last Used - %v \n", card.Id, card.LastUsed)
		}
		fmt.Fprintf(w, "^ ACTIVE CARDS")
	})

	log.Printf("Listening on 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
