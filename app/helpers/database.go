package helpers

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	User           *mongo.Collection
	Stock          *mongo.Collection
	StockInfo      *mongo.Collection
	Migration      *mongo.Collection
	UserFavourites *mongo.Collection
}

var CurrentDb *DB

const stockCodeIndexName = "stock_code_index_1"
const stockInfoCodeDateIndexName = "stock_info_code_date_index_1"
const stockInfoGainIndexName = "stock_info_gain_index_1"
const stockInfoNameIndexName = "stock_info_name_index_1"

func doesIndexExist(ctx context.Context, collection *mongo.Collection, indexName string) (bool, error) {
	cursor, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}

	defer cursor.Close(ctx)

	var indexDoc bson.M
	for cursor.Next(ctx) {
		if err := cursor.Decode(&indexDoc); err != nil {
			return false, err
		}

		// Check if the index name matches
		if name, ok := indexDoc["name"].(string); ok && name == indexName {
			return true, nil
		}
	}

	return false, nil
}

func CreateDBInstance() {
	connectionString := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	userCollName := os.Getenv("DB_USER_COLLECTION_NAME")
	stockCollName := os.Getenv("DB_STOCK_COLLECTION_NAME")
	stockInfoCollName := os.Getenv("DB_STOCK_INFO_COLLECTION_NAME")
	migrationCollName := os.Getenv("DB_MIGRATION_COLLECTION_NAME")
	userFavouritesCollName := os.Getenv("DB_USER_FAVOURITES_COLLECTION_NAME")
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
		return
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
		return
	}

	userCollection := client.Database(dbName).Collection(userCollName)
	stockCollection := client.Database(dbName).Collection(stockCollName)
	stockInfoCollection := client.Database(dbName).Collection(stockInfoCollName)
	migrationCollection := client.Database(dbName).Collection(migrationCollName)
	userFavouritesCollection := client.Database(dbName).Collection(userFavouritesCollName)

	stockCodeIndexExists, err := doesIndexExist(context.Background(), stockCollection, stockCodeIndexName)
	if err != nil {
		log.Fatal(err)
	}

	if !stockCodeIndexExists {
		// Create the index
		indexModel := mongo.IndexModel{
			Keys:    bson.M{"stockCode": 1},
			Options: options.Index().SetUnique(true).SetName(stockCodeIndexName),
		}

		_, err := stockCollection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Stock Code Index created successfully.")
	} else {
		fmt.Println("Stock Code Index already exists.")
	}

	stockInfoCodeDateIndexExists, err := doesIndexExist(context.Background(), stockInfoCollection, stockInfoCodeDateIndexName)
	if err != nil {
		log.Fatal(err)
	}

	if !stockInfoCodeDateIndexExists {
		// Create the index
		indexModel := mongo.IndexModel{
			Keys: bson.D{
				{Key: "stockCode", Value: 1},
				{Key: "date", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName(stockInfoCodeDateIndexName),
		}

		_, err := stockInfoCollection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Stock Info Code Date Index created successfully.")
	} else {
		fmt.Println("Stock Info Code Date Index already exists.")
	}

	stockInfoGainIndexExists, err := doesIndexExist(context.Background(), stockInfoCollection, stockInfoGainIndexName)
	if err != nil {
		log.Fatal(err)
	}

	if !stockInfoGainIndexExists {
		// Create the index
		indexModel := mongo.IndexModel{
			Keys: bson.D{
				{Key: "gain", Value: 1},
				{Key: "date", Value: 1},
			},
			Options: options.Index().SetName(stockInfoGainIndexName),
		}

		_, err := stockInfoCollection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Stock Info Gain Index created successfully.")
	} else {
		fmt.Println("Stock Info Gain Index already exists.")
	}

	stockInfoNameIndexExists, err := doesIndexExist(context.TODO(), stockInfoCollection, stockInfoNameIndexName)
	if err != nil {
		log.Fatal(err)
	}

	if !stockInfoNameIndexExists {
		// Create the index
		indexModel := mongo.IndexModel{
			Keys: bson.D{
				{Key: "stockName", Value: 1},
			},
			Options: options.Index().SetName(stockInfoNameIndexName),
		}

		_, err := stockInfoCollection.Indexes().CreateOne(context.Background(), indexModel)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Stock Info Name Index created successfully.")
	} else {
		fmt.Println("Stock Info Name Index already exists.")
	}

	CurrentDb = &DB{User: userCollection, Stock: stockCollection, StockInfo: stockInfoCollection, Migration: migrationCollection, UserFavourites: userFavouritesCollection}
	fmt.Println("Connected to MongoDB")
}
