package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"insider/models"
	"log"
	"os"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func New(indexName, filePath, elasticUrl string) (*elasticsearch.Client, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{elasticUrl},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize elastic client: %w", err)
	}

	resp, err := es.Indices.Exists([]string{indexName})
	if err != nil {
		return nil, fmt.Errorf("failed to checking index: %w", err)
	}

	if resp.StatusCode == 404 {
		IndexFromFile(es, indexName, filePath)
	}

	return es, nil
}

// CreateIndex creates an index in Elasticsearch
func CreateIndex(es *elasticsearch.Client, indexName string) error {
	res, err := es.Indices.Create(indexName)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to create index: %s", res.String())
	}

	log.Printf("index %s created", indexName)
	return nil
}

// BulkIndex bulk indexes documents into Elasticsearch
func BulkIndex(es *elasticsearch.Client, indexName string, items []map[string]interface{}) error {
	var buf bytes.Buffer
	for _, item := range items {
		meta := []byte(fmt.Sprintf(`{ "index" : { "_index" : "%s" }}`, indexName))
		buf.Write(meta)
		buf.WriteByte('\n')

		doc, err := json.Marshal(item)
		if err != nil {
			log.Printf("failed to marshal document: %v", err)
			continue
		}
		buf.Write(doc)
		buf.WriteByte('\n')
	}

	res, err := es.Bulk(bytes.NewReader(buf.Bytes()), es.Bulk.WithIndex(indexName))
	if err != nil {
		return fmt.Errorf("failed to bulk index documents: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("failed to bulk index documents: %s", res.String())
	}

	log.Println("documents indexed successfully")
	return nil
}

// IndexFromFile indexes documents from a JSON file into Elasticsearch
func IndexFromFile(es *elasticsearch.Client, indexName string, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Decode the JSON file
	var items []map[string]interface{}
	err = json.NewDecoder(file).Decode(&items)
	if err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	// Create the index
	err = CreateIndex(es, indexName)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Bulk index the documents
	err = BulkIndex(es, indexName, items)
	if err != nil {
		return fmt.Errorf("failed to bulk index documents: %w", err)
	}

	return nil
}

func ConstructQuery(filterParams models.FilterParams) map[string]interface{} {
	query := make(map[string]interface{})
	boolQuery := make(map[string]interface{})
	boolQuery["must"] = []interface{}{}

	for _, condition := range filterParams.Conditions {
		switch condition.Operation {
		case "equals":
			matchQuery := map[string]interface{}{
				"match": map[string]interface{}{
					condition.FieldName: condition.Value,
				},
			}
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), matchQuery)
		case "less_than":
			numericValue, err := strconv.Atoi(condition.Value)
			if err != nil {
				log.Printf("Error converting value to integer: %s", err)
				continue
			}
			rangeQuery := map[string]interface{}{
				"range": map[string]interface{}{
					condition.FieldName: map[string]interface{}{
						"lt": numericValue,
					},
				},
			}
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), rangeQuery)
		case "greater_than":
			numericValue, err := strconv.Atoi(condition.Value)
			if err != nil {
				log.Printf("Error converting value to integer: %s", err)
				continue
			}
			rangeQuery := map[string]interface{}{
				"range": map[string]interface{}{
					condition.FieldName: map[string]interface{}{
						"gt": numericValue,
					},
				},
			}
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), rangeQuery)
		default:
			log.Printf("Unsupported operation: %s", condition.Operation)
		}
	}

	query["query"] = map[string]interface{}{
		"bool": boolQuery,
	}

	return query
}

func DoSearch(ctxBg context.Context, queryBytes []byte, es *elasticsearch.Client) models.SearchResponse {
	req := esapi.SearchRequest{
		Index: []string{"insider*"},
		Body:  bytes.NewBuffer(queryBytes),
	}

	// Perform the search request
	res, err := req.Do(context.Background(), es)
	if err != nil {
		log.Fatalf("Error executing search request: %s", err)
	}
	defer res.Body.Close()

	// Check if response status code is OK
	if res.IsError() {
		log.Fatalf("Search request failed with status code: %d", res.StatusCode)
	}

	// Print the response body
	fmt.Println("Search results:")
	var searchRes models.SearchResponse
	if err := json.NewDecoder(res.Body).Decode(&searchRes); err != nil {
		log.Fatalf("Error decoding search response: %s", err)
	}

	return searchRes
}
