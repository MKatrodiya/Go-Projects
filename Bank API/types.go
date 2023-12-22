package main

import (
	"math/rand"
	"time"
)

type Account struct {
	ID            int       `json: "id"`
	FirstName     string    `json: "firstName"`
	LastName      string    `json: "lastName"`
	AccountNumber int64     `json: "accountNumber"`
	Balance       int64     `json: "balance"`
	CreatedAt     time.Time `json: "createdAt"`
}

type CreateAccountRequest struct {
	FirstName string `json: "firstname"`
	LastName  string `json: "lastname"`
}

func CreateAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName:     firstName,
		LastName:      lastName,
		AccountNumber: rand.Int63n(100000000),
		CreatedAt:     time.Now().UTC(),
	}
}
