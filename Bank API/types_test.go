package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewAccount(t *testing.T) {
	account, err := CreateAccount("A", "B", "pass1234")
	assert.Nil(t, err)

	fmt.Printf("%+v\n", account)
}
