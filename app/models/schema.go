package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Stock struct {
	ID         primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	StockCode  string             `json:"stockCode"`
	StockName  string             `json:"stockName" validate:"required"`
	StockGroup string             `json:"stockGroup" validate:"required"`
	StockType  string             `json:"stockType" validate:"required"`
}

type StockInfo struct {
	ID        primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	Date      string             `json:"date" validate:"required"`
	Open      float64            `json:"open" validate:"required"`
	High      float64            `json:"high" validate:"required"`
	Low       float64            `json:"low" validate:"required"`
	Close     float64            `json:"close" validate:"required"`
	Last      float64            `json:"last" validate:"required"`
	Gain      float64            `json:"gain" validate:"required"`
	StockCode string             `json:"stockCode" validate:"required"`
	StockName string             `json:"stockName" validate:"required"`
}

type StockJson struct {
	ID         primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	StockCode  string             `json:"stockCode"`
	StockName  string             `json:"stockName"`
	StockGroup string             `json:"stockGroup"`
	StockType  string             `json:"stockType"`
	Info       *StockInfo         `json:"info"`
}

type StockMigration struct {
	ID       primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	FileName string             `json:"fileName"`
	// MigratedAt string             `json:"migratedAt"`
	DataDate string `json:"dataDate"`
}

type User struct {
	ID        primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	Email     string             `json:"email"`
	CreatedAt time.Time          `json:"createdAt"`
	Token     string             `json:"token" bson:"-"`
}

type UserFavourite struct {
	ID   primitive.ObjectID `json:"_id,omitempty"  bson:"_id,omitempty"`
	User primitive.ObjectID `json:"user"`
	// save stockCode
	Stocks []string `json:"stocks"`
}
