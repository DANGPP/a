package main

import (
	"crypto/rand"
	"crypto/rsa"

	// "fmt"
	"log"
	"net/http"
	// "time"
	// "github.com/golang-jwt/jwt/v5"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

var user = map[string]string{
	"name": "1234",
}

func init() {
	var err error
	privateKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("%v", err)
	}
	publicKey = &privateKey.PublicKey

}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	_, ok := user[username]
	log.Fatal(password)
	log.Fatal(ok)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
