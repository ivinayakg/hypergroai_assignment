package helpers

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

var bucketName string
var projectID string

type FileType struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func CreateGoogleCloudStorageClient() {
	bucketName = os.Getenv("GOOGLE_CLOUD_STORAGE_BUCKET_NAME")
	projectID = os.Getenv("GOOGLE_CLOUD_STORAGE_PROJECT_ID")
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Println(fmt.Errorf("storage.NewClient: %w", err))
	}
	defer client.Close()

	_, err = client.Bucket(bucketName).Attrs(ctx)
	if err != nil && err.Error() != "storage: bucket doesn't exist" {
		fmt.Println(err.Error() == "storage: bucket doesn't exist")
	} else if err != nil && err.Error() == "storage: bucket doesn't exist" {
		// Creates a Bucket instance.
		bucket := client.Bucket(bucketName)
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		if err := bucket.Create(ctx, projectID, nil); err != nil {
			log.Fatalf("Failed to create bucket: %v", err)
		}

		fmt.Printf("Bucket %v created.\n", bucketName)
	}

	fmt.Println("Google Cloud Storage Setup Success")
}

// listFiles lists objects within specified bucket.
func ListFiles() ([]*FileType, error) {
	results := []*FileType{}
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return results, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	it := client.Bucket(bucketName).Objects(ctx, nil)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err == nil {
			file := &FileType{
				Name: attrs.Name, Type: attrs.ContentType,
			}
			results = append(results, file)
		}
	}
	return results, nil
}

func DownloadCSVFileIntoMemory(object string) ([][]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	rc, err := client.Bucket(bucketName).Object(object).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	data, err := io.ReadAll(rc)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	reader := csv.NewReader(bytes.NewReader(data))

	// Read all records from the CSV data
	csvData, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV data: %w", err)
	}

	return csvData, nil
}

// uploadFile uploads an object.
func UploadCSVFile(w io.Writer, fileName string, file io.Reader) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	o := client.Bucket(bucketName).Object(fileName)

	// Optional: set a generation-match precondition to avoid potential race
	// conditions and data corruptions. The request to upload is aborted if the
	// object's generation number does not match your precondition.
	// ...

	// Upload an object with storage.Writer.
	wc := o.NewWriter(ctx)
	wc.ContentType = "text/csv"
	if _, err = io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %w", err)
	}
	return nil
}
