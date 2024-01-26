package models

import (
	"context"
	"fmt"
	"time"

	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetStocksData(page int, dataDate string, stockName string, pageSize int64, top bool) ([]*StockInfo, error) {
	var cacheKey = fmt.Sprintf("stock:page-%d:date-%v:name-%v:pageSize-%d:top-%v", page, dataDate, stockName, pageSize, top)
	var results = []*StockInfo{}
	err := helpers.Redis.GetJSON(cacheKey, results)
	if err == nil {
		return results, nil
	}

	stockInfoFilters := bson.M{"date": dataDate}
	queryOptions := options.Find().SetLimit(pageSize).SetSkip(pageSize * int64(page-1))

	if stockName != "" {
		stockInfoFilters["stockName"] = bson.M{"$regex": primitive.Regex{Pattern: fmt.Sprintf("^%v", stockName), Options: "i"}}
	}

	if top {
		queryOptions.SetSort(bson.M{"gain": -1})
	}

	curr, err := helpers.CurrentDb.StockInfo.Find(context.TODO(), stockInfoFilters, queryOptions)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer curr.Close(context.TODO())

	for curr.Next(context.TODO()) {
		var result StockInfo
		e := curr.Decode(&result)
		if e != nil {
			fmt.Println(err.Error())
			continue
		}
		results = append(results, &result)
	}

	go helpers.Redis.SetJSON(cacheKey, results, time.Until(time.Now().Add(time.Second*600)))

	return results, nil
}

func GetStockDetail(stockCode string) (*StockJson, error) {
	ctx := context.TODO()
	var result StockJson

	stockFilter := bson.M{"stockCode": stockCode}

	err := helpers.CurrentDb.Stock.FindOne(ctx, stockFilter).Decode(&result)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var stockInfo StockInfo
	err = helpers.CurrentDb.StockInfo.FindOne(ctx, stockFilter).Decode(&stockInfo)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result.Info = &stockInfo

	return &result, nil
}
