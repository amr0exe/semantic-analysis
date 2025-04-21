package main

import (
	"fmt"
	"log"
	"context"
	"regexp"
	"strings"
	"encoding/json"
	"time"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func InitDB() (context.Context, *mongo.Collection, context.CancelFunc) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file...")
	}
	dbURL := os.Getenv("DB_URL")
	
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbURL))
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}
	collection := client.Database("papergeneration").Collection("paper")
	return ctx, collection, cancel
}

func extractEnglishText(input string) string {
	re := regexp.MustCompile(`[A-Za-z0-9\s\p{P}]+`)
	matches := re.FindAllString(input, -1)
	englishText := strings.Join(matches, " ")
	englishText = strings.Join(strings.Fields(englishText), " ")

	// remove 8.
	numberPattern := regexp.MustCompile(`^\d+\.\s*`)
	englishText = numberPattern.ReplaceAllString(englishText, "")

	return englishText
}

func saveToFile(data Prequest) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal("Error marshalling data", err)
	}

	err = os.WriteFile("json.txt", jsonData, 0644)
	if err != nil {
		log.Fatal("Error writing to json.txt", err)
	}

	fmt.Print("writing success!")
}
