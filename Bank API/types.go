package main

import "math/rand"

type Account struct {
	ID            int    `json: "id"`
	FirstName     string `json: "firstName"`
	LastName      string `json: "lastName"`
	AccountNumber int64  `json: "accountNumber"`
	Balance       int64  `json: "balance"`
}

func CreateAccount(firstName, lastName string) *Account {
	return &Account{
		ID:            rand.Intn(10000000),
		FirstName:     firstName,
		LastName:      lastName,
		AccountNumber: rand.Int63n(1000000000000),
	}
}
