package main

import (
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID                int       `json:"id"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	AccountNumber     int64     `json:"accountNumber"`
	EncryptedPassword string    `json:"-"`
	Balance           int64     `json:"balance"`
	CreatedAt         time.Time `json:"createdAt"`
}

type CreateAccountRequest struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  string `json:"password"`
}

type TransferRequest struct {
	ToAccount int `json:"toAccount"`
	Amount    int `json:"amount"`
}

type SignInRequest struct {
	AccountNumber int    `json:"accountNumber"`
	Password      string `json:"password"`
}

type SignInResponse struct {
	JWTToken      string `json:"jwtToken"`
	AccountNumber int    `json:"accountNumber"`
}

func CreateAccount(firstName, lastName, password string) (*Account, error) {
	encryptedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		AccountNumber:     rand.Int63n(100000000),
		EncryptedPassword: string(encryptedPw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
