package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/opensearch-project/opensearch-go/v2"
)

type Event struct{}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, event Event) (string, error) {
	esURL := os.Getenv("ES_URL")
	esUsername := os.Getenv("ES_USERNAME")
	esPassword := os.Getenv("ES_PASSWORD")
	dryRunStr := os.Getenv("DRY_RUN")
	indexAgeLimitDaysStr := os.Getenv("INDEX_AGE_LIMIT_DAYS")

	dryRun, err := strconv.ParseBool(dryRunStr)
	if err != nil {
		return "", fmt.Errorf("invalid DRY_RUN value: %v", err)
	}

	indexAgeLimitDays, err := strconv.Atoi(indexAgeLimitDaysStr)
	if err != nil {
		return "", fmt.Errorf("invalid INDEX_AGE_LIMIT_DAYS value: %v", err)
	}

	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{esURL},
		Username:  esUsername,
		Password:  esPassword,
	})
	if err != nil {
		return "", fmt.Errorf("error creating OpenSearch client: %v", err)
	}

	indices, err := getIndices(client)
	if err != nil {
		return "", fmt.Errorf("error getting indices: %v", err)
	}

	cutoff := time.Now().AddDate(0, 0, -indexAgeLimitDays)
	fmt.Printf("Cutoff date: %s\n", cutoff.Format("2006-01-02"))

	var indicesToDelete []string

	// Count total indices matching the pattern
	var totalLogstashIndices int
	for _, index := range indices {
		if strings.HasPrefix(index, "logstash-logs-") {
			totalLogstashIndices++
			if shouldDelete(index, cutoff) {
				indicesToDelete = append(indicesToDelete, index)
			}
		}
	}

	// Print the total number of logstash-logs-* indices and the number of indices to delete
	fmt.Printf("Total logstash-logs-* indices: %d\n", totalLogstashIndices)
	fmt.Printf("Cutoff date: %s\n", cutoff.Format("2006-01-02"))

	if len(indicesToDelete) > 0 {
		fmt.Println("Indices to delete:")
		for _, index := range indicesToDelete {
			fmt.Println(index)
		}
		fmt.Printf("Total Indices to delete: %v\n", len(indicesToDelete))
	} else {
		fmt.Println("No indices to delete.")
	}

	if !dryRun {
		for _, index := range indicesToDelete {
			err := deleteIndex(client, index)
			if err != nil {
				fmt.Printf("Error deleting index %s: %v\n", index, err)
			}
		}
	} else {
		fmt.Printf("Dry run: No indices have been deleted.\n")
	}

	return "Process completed", nil
}

func getIndices(client *opensearch.Client) ([]string, error) {
	var indices []string
	res, err := client.Cat.Indices()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Read the raw response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Split the response into lines
	lines := strings.Split(string(body), "\n")

	for _, line := range lines {
		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split the line into fields
		fields := strings.Fields(line)
		if len(fields) > 0 {
			indexName := fields[2] // The index name is typically the third field
			if strings.HasPrefix(indexName, "logstash-logs-") {
				indices = append(indices, indexName)
			}
		}
	}

	return indices, nil
}

func shouldDelete(index string, cutoff time.Time) bool {
	// Extract the date from the index name
	parts := strings.Split(index, "-")
	if len(parts) == 0 {
		return false
	}

	indexDateStr := parts[len(parts)-1] // Assuming date is at the end of the index name

	// Parse the date from the index name
	indexDate, err := time.Parse("2006.01.02", indexDateStr)
	if err != nil {
		fmt.Printf("Error parsing date from index %s: %v\n", index, err)
		return false
	}

	// Print dates for debugging
	// fmt.Printf("Index date: %s, Cutoff date: %s | will be deleted = %v\n", indexDate.Format("2006-01-02"), cutoff.Format("2006-01-02"), indexDate.Before(cutoff))

	// Compare the index date with the cutoff date
	return indexDate.Before(cutoff)
}

func deleteIndex(client *opensearch.Client, index string) error {
	res, err := client.Indices.Delete([]string{index})
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
