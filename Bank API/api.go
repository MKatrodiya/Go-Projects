package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) error {
	w.WriteHeader(statusCode)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(payload)
}

type APIError struct {
	Error string
}

type APIServer struct {
	listenPort string
}

func createServer(listenPort string) *APIServer {
	apiServer := APIServer{
		listenPort: listenPort,
	}
	return &apiServer
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", s.handleAccount)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
