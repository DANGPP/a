package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func getVaultSecretPathFromDB(clusterName string) (string, error) {
	dbURL := "host=127.0.0.1 port=5433 user=dangdn2 password=1 dbname=db-part-k8s sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return "", fmt.Errorf("connect db: %w", err)
	}
	defer db.Close()

	var path string
	query := `SELECT path FROM k8s_configs WHERE name = $1 LIMIT 1`
	err = db.QueryRow(query, clusterName).Scan(&path)
	if err != nil {
		return "", fmt.Errorf("query path: %w", err)
	}

	return path, nil
}
func main() {
	ss, err := getVaultSecretPathFromDB("k8s-destop")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(ss)
}
