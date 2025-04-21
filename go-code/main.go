package main

import (
	"net/http"
	"fmt"
	"log"
)

func main() {
	http.HandleFunc("POST /api/generate-questions", GenerateQuestions())	

	fmt.Println("Server is running at localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal("failed to start the server", err)
	}
}
