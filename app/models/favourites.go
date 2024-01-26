package models

import (
	"context"
	"fmt"

	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func AddFavourite(userId string, stockCode string) (*UserFavourite, error) {
	ctx := context.TODO()

	userOID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var user User
	err = helpers.CurrentDb.User.FindOne(ctx, bson.M{"_id": userOID}).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var stock Stock
	err = helpers.CurrentDb.Stock.FindOne(ctx, bson.M{"stockCode": stockCode}).Decode(&stock)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var favourite UserFavourite
	err = helpers.CurrentDb.UserFavourites.FindOne(ctx, bson.M{"user": userOID}).Decode(&favourite)
	if err != nil && err != mongo.ErrNoDocuments {
		fmt.Println(err)
		return nil, err
	}
	if err == mongo.ErrNoDocuments {
		stocksFavouritesList := []string{stockCode}
		result, err := helpers.CurrentDb.UserFavourites.InsertOne(ctx, UserFavourite{User: userOID, Stocks: stocksFavouritesList})
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		favourite = UserFavourite{}
		favourite.ID = result.InsertedID.(primitive.ObjectID)
		favourite.User = userOID
		favourite.Stocks = stocksFavouritesList
	} else {
		stocksFavouritesList := favourite.Stocks
		stocksFavouritesList = append(stocksFavouritesList, stockCode)
		_, err := helpers.CurrentDb.UserFavourites.UpdateByID(ctx, favourite.ID, bson.M{"$set": bson.M{"stocks": stocksFavouritesList}})
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		favourite.Stocks = stocksFavouritesList
	}

	return &favourite, nil
}

func RemoveFavourite(userId string, stockCode string) (*UserFavourite, error) {
	ctx := context.TODO()

	userOID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var user User
	err = helpers.CurrentDb.User.FindOne(ctx, bson.M{"_id": userOID}).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var favourite UserFavourite
	err = helpers.CurrentDb.UserFavourites.FindOne(ctx, bson.M{"user": userOID}).Decode(&favourite)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	stocksFavouritesList := []string{}
	for _, entry := range favourite.Stocks {
		if entry == stockCode {
			continue
		}
		stocksFavouritesList = append(stocksFavouritesList, entry)
	}

	_, err = helpers.CurrentDb.UserFavourites.UpdateByID(ctx, favourite.ID, bson.M{"$set": bson.M{"stocks": stocksFavouritesList}})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	favourite.Stocks = stocksFavouritesList

	return &favourite, nil
}

func GetFavouriteStocks(userId string) ([]*StockInfo, error) {
	ctx := context.TODO()

	userOID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var user User
	err = helpers.CurrentDb.User.FindOne(ctx, bson.M{"_id": userOID}).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var userFavourite UserFavourite
	err = helpers.CurrentDb.UserFavourites.FindOne(ctx, bson.M{"user": userOID}).Decode(&userFavourite)
	if err != nil {
		fmt.Println(err)
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("no user favourites exists")
		}
		return nil, err
	}

	// latest data only
	migration, _ := GetLastMigration()

	var results []*StockInfo
	stockInfoFilters := bson.M{"stockCode": bson.M{"$in": userFavourite.Stocks}, "date": migration.DataDate}

	curr, err := helpers.CurrentDb.StockInfo.Find(context.TODO(), stockInfoFilters)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var result StockInfo
		e := curr.Decode(&result)
		if e != nil {
			fmt.Println(err)
			continue
		}
		results = append(results, &result)
	}

	return results, nil
}
