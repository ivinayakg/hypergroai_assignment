package models

import (
	"context"
	"fmt"

	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetLastMigration() (*StockMigration, error) {
	var result StockMigration

	err := helpers.CurrentDb.Migration.FindOne(context.TODO(), bson.M{}, options.FindOne().SetSort(bson.M{"datadate": -1})).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &result, nil
}

func GetMigrationForTime(date string) (*StockMigration, error) {
	var result StockMigration

	err := helpers.CurrentDb.Migration.FindOne(context.TODO(), bson.M{"datadate": date}).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &result, nil
}

func GetMigrations() ([]*StockMigration, error) {
	var results []*StockMigration

	curr, err := helpers.CurrentDb.Migration.Find(context.TODO(), bson.M{}, options.Find().SetSort(bson.M{"datadate": 1}))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for curr.Next(context.TODO()) {
		var result StockMigration
		e := curr.Decode(&result)
		if e != nil {
			fmt.Println(err)
			continue
		}
		results = append(results, &result)
	}

	return results, nil
}
