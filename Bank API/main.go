package main

import (
	"log"
)

func main() {
	store, err := CreatePostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := CreateAPIServer(":8080", store)
	server.Run()
}
