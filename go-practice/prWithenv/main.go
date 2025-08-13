package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("gay vl")
	}

	stringgL := os.Getenv("GAY")
	fmt.Println(stringgL)
}
