package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	seed := flag.Bool("seed", false, "seed the database")
	flag.Parse()

	store, err := CreatePostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	if *seed {
		seedAccounts(store)
	}

	server := CreateAPIServer(":8080", store)
	server.Run()
}

func seedAccount(store Storage, firstName, lastName, pw string) *Account {
	acc, err := CreateAccount(firstName, lastName, pw)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.CreateAccount(acc); err != nil {
		log.Fatal(err)
	}

	fmt.Println("new account => ", acc.AccountNumber)

	return acc
}

func seedAccounts(s Storage) {
	seedAccount(s, "Seed", "User", "seedPw")
}
