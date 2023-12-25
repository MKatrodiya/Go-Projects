package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type APIServer struct {
	listenPort string
	store      Storage
}

func CreateAPIServer(listenPort string, store Storage) *APIServer {
	apiServer := APIServer{
		listenPort: listenPort,
		store:      store,
	}
	return &apiServer
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/signin", makeHandleFunc(s.HandleSignIn))
	router.HandleFunc("/account", makeHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHandleFunc(s.handleGetAccountByID), s.store))
	router.HandleFunc("/transfer", makeHandleFunc(s.handleTransfer))

	log.Println("Server running on port ", s.listenPort)

	http.ListenAndServe(s.listenPort, router)
}

func (s *APIServer) HandleSignIn(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not supported %s", r.Method)
	}

	req := new(SignInRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := s.store.GetAccountByNumber(req.AccountNumber)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(account.EncryptedPassword), []byte(req.Password)); err != nil {
		permissionDeniedResponse(w)
		return nil
	}

	token, err := createJWT(account)

	if err != nil {
		return err
	}

	resp := SignInResponse{
		JWTToken:      token,
		AccountNumber: req.AccountNumber,
	}
	WriteJSON(w, http.StatusOK, resp)
	return nil
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not supported %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {

		id, err := getIDParam(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not supported %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := &CreateAccountRequest{}

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		return err
	}

	account, err := CreateAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDParam(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return WriteJSON(w, http.StatusOK, transferReq)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func makeHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo1MjY0MDMwNSwiZXhwaXJlc0F0IjoxNzAzNDYzNDI4fQ.so1ijQ9xJXBCJ4V0bC_NeUJPsIunlmtW2aoN2dtY9w4
func createJWT(account *Account) (string, error) {
	secret := os.Getenv("Bank_JWT_SECRET")

	claims := &jwt.MapClaims{
		"expiresAt":     jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		"accountNumber": account.AccountNumber,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func permissionDeniedResponse(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, APIError{
		Error: "permission denied",
	})
}

func withJWTAuth(handler http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Passing through JWT middleware")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)

		if err != nil {
			permissionDeniedResponse(w)
			return
		}
		if !token.Valid {
			permissionDeniedResponse(w)
			return
		}

		id, err := getIDParam(r)
		if err != nil {
			permissionDeniedResponse(w)
			return
		}
		account, err := s.GetAccountByID(id)
		if err != nil {
			permissionDeniedResponse(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		accountNumber, ok := claims["accountNumber"].(float64)
		if !ok {
			permissionDeniedResponse(w)
			return
		}

		if account.AccountNumber != int64(accountNumber) {
			permissionDeniedResponse(w)
			return
		}
		handler(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	jwtSecret := os.Getenv("Bank_JWT_SECRET") // Declare jwtSecret variable
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jwtSecret), nil
	})
}

func WriteJSON(w http.ResponseWriter, statusCode int, payload any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(payload)
}

func getIDParam(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id: %s", idStr)
	}
	return id, nil
}
