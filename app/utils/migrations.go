package utils

import (
	"context"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"time"

	"github.com/ivinayakg/hypergroai_assignment/helpers"
	"github.com/ivinayakg/hypergroai_assignment/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var CSVFileDataIndex = map[string]int32{
	"SC_CODE":  0,
	"SC_NAME":  1,
	"SC_GROUP": 2,
	"SC_TYPE":  3,
	"OPEN":     4,
	"HIGH":     5,
	"LOW":      6,
	"CLOSE":    7,
	"LAST":     8,
}

func CreateMigration(wg *sync.WaitGroup, fileCh chan struct{}, fileData [][]string, fileName string) {
	defer wg.Done()
	{
		// for creating stock instances
		// Bulk write options
		bulkOptions := options.BulkWrite().SetOrdered(false)

		// // Prepare the bulk write operations
		var bulkOps []mongo.WriteModel

		for _, data := range fileData {
			filter := bson.M{"stockCode": data[CSVFileDataIndex["SC_CODE"]]}
			update := bson.M{
				"$set": bson.M{
					"stockCode":  data[CSVFileDataIndex["SC_CODE"]],
					"stockType":  data[CSVFileDataIndex["SC_TYPE"]],
					"stockGroup": data[CSVFileDataIndex["SC_GROUP"]],
					"stockName":  strings.TrimSpace(data[CSVFileDataIndex["SC_NAME"]]),
				},
			}

			// Create an upsert operation
			upsert := mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(update).
				SetUpsert(true)

			bulkOps = append(bulkOps, upsert)
		}

		// Execute the bulk write operations
		result, err := helpers.CurrentDb.Stock.BulkWrite(context.Background(), bulkOps, bulkOptions)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Inserted %d Stock documents\n", result.UpsertedCount)
		// fmt.Printf("Inserted %d documents\n", 0)
	}
	{
		// stock info instances
		// Bulk write options
		bulkOptions := options.BulkWrite().SetOrdered(false)

		// // Prepare the bulk write operations
		var bulkOps []mongo.WriteModel

		fileDateTimeString, _ := helpers.ConvertFilenameToDate(fileName)
		fileDateString, _ := helpers.ConvertTimeIntoDate(fileDateTimeString)

		for _, data := range fileData {
			filter := bson.M{"stockCode": data[CSVFileDataIndex["SC_CODE"]], "date": fileDateString}
			close := helpers.ConvertStringToFloat(data[CSVFileDataIndex["CLOSE"]])
			last := helpers.ConvertStringToFloat(data[CSVFileDataIndex["LAST"]])
			netGain := ((close - last) / last) * 100
			update := bson.M{
				"$set": bson.M{
					"stockCode": data[CSVFileDataIndex["SC_CODE"]],
					"date":      fileDateString,
					"open":      helpers.ConvertStringToFloat(data[CSVFileDataIndex["OPEN"]]),
					"high":      helpers.ConvertStringToFloat(data[CSVFileDataIndex["HIGH"]]),
					"low":       helpers.ConvertStringToFloat(data[CSVFileDataIndex["LOW"]]),
					"close":     close,
					"last":      last,
					"gain":      math.Round(netGain*100) / 100,
					"stockName": strings.TrimSpace(data[CSVFileDataIndex["SC_NAME"]]),
				},
			}

			// Create an upsert operation
			upsert := mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(update).
				SetUpsert(true)

			bulkOps = append(bulkOps, upsert)
		}

		// Execute the bulk write operations
		result, err := helpers.CurrentDb.StockInfo.BulkWrite(context.Background(), bulkOps, bulkOptions)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Inserted %d Stock Info documents\n", result.UpsertedCount)
		// fmt.Printf("Inserted %d documents\n", 0)
	}
	fileCh <- struct{}{}
}

func RunMigrations() {
	filesName, err := helpers.ListFiles()
	if err != nil {
		fmt.Println("Error while migrating, ", err)
	}

	ctx := context.TODO()
	migratedFilesFilter := bson.M{}

	curr, err := helpers.CurrentDb.Migration.Find(ctx, migratedFilesFilter)
	if err != nil {
		fmt.Println(err)
	}
	defer curr.Close(context.TODO())

	var migratedFilesNames []string
	for curr.Next(context.TODO()) {
		var result models.StockMigration
		e := curr.Decode(&result)
		if e != nil {
			fmt.Println(err)
		}
		migratedFilesNames = append(migratedFilesNames, result.FileName)
	}

	var toMigrateFileNames []*helpers.FileType
	for _, file := range filesName {
		if !helpers.ContainsString(&migratedFilesNames, &file.Name) {
			toMigrateFileNames = append(toMigrateFileNames, file)
		}
	}

	// Create a wait group to wait for the completion of goroutines
	var wg sync.WaitGroup

	for _, file := range toMigrateFileNames {
		if file.Type != "text/csv" {
			continue
		}

		// Create a channel for each file to coordinate sub-goroutines
		fileCh := make(chan struct{}, 1)
		csvData, err := helpers.DownloadCSVFileIntoMemory(file.Name)
		if err != nil {
			fmt.Println(err)
		} else {
			noOfOperations := 4
			var operationsSizePerRoutine int = (len(csvData) + noOfOperations - 1) / noOfOperations

			for i := 1; i <= noOfOperations; i++ {
				wg.Add(1)
				start := (i - 1) * operationsSizePerRoutine
				if start == 0 {
					start = 1
				}
				end := i * operationsSizePerRoutine
				if end > len(csvData) {
					end = len(csvData)
				}
				go CreateMigration(&wg, fileCh, csvData[start:end], file.Name)
				<-fileCh // Wait for the previous sub-goroutine to complete
				time.Sleep(time.Second * 5)
			}
			// create a migration data on the database
			dataTimeDate, _ := helpers.ConvertFilenameToDate(file.Name)
			dateDate, _ := helpers.ConvertTimeIntoDate(dataTimeDate)
			stockMigration := &models.StockMigration{FileName: file.Name, DataDate: dateDate}
			res, err := helpers.CurrentDb.Migration.InsertOne(context.Background(), stockMigration)
			if err != nil {
				fmt.Println(err)
			}

			stockMigration.ID = res.InsertedID.(primitive.ObjectID)

			fmt.Printf("Migration Successfully Completed of %v, of type %v. The log is saved in the DB with migration Id - %+v\n", file.Name, file.Type, stockMigration.ID.Hex())
		}

		time.Sleep(time.Second * 20)
	}

	fmt.Println("Migrations Completed Successfully for all the current files in the DB")
}
