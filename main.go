package main

import (
	"fmt"
	transactions "github.com/antonOO/CardsAndTransactions/transactions"
	"html"
	"log"
	"net/http"
)

func main() {

	// Server supports only a single client as there is no budget :D
	// Nor part of the assignment
	client := transactions.Client{}
	println(client)

	http.HandleFunc("/add-card", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
