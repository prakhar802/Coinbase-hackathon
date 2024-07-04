package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Transaction struct {
	ID      primitive.ObjectID `json: "id,omitempty" bson:"_id,omitempty"`
	Amount  int                `json: "amount"`
	Company string             `json: "company"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello World")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal("Connection error : ", err)
	}

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected MONGDB Atlas")

	collection = client.Database("golang_db").Collection("transactions")
	app := fiber.New()

	app.Get("/api/transactions", getTransaction)
	app.Post("/api/transactions", createTransaction)

	port := "5000"

	log.Fatal(app.Listen("0.0.0.0:" + port))

}

func getTransaction(c *fiber.Ctx) error {
	var transactions []Transaction
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return err
	}

	for cursor.Next(context.Background()) {
		var transaction Transaction
		if err := cursor.Decode(&transaction); err != nil {
			return err
		}
		transactions = append(transactions, transaction)
	}
	return c.JSON(transactions)
}

func createTransaction(c *fiber.Ctx) error {
	transaction := new(Transaction)

	err := c.BodyParser(transaction)

	if err != nil {
		return err
	}

	if transaction.Company == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Transaction field cannot be empty"})
	}
	insertResult, err := collection.InsertOne(context.Background(), transaction)
	if err != nil {
		return err
	}
	transaction.ID = insertResult.InsertedID.(primitive.ObjectID)
	return c.Status(201).JSON(transaction)

}
