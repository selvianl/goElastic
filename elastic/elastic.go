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

type Elastic struct {
	*elasticsearch.Client
}

func New(indexName, filePath, elasticUrl string) (*Elastic, error) {
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

	e := &Elastic{es}

	if resp.StatusCode == 404 {
		e.IndexFromFile(indexName, filePath)
	}

	return &Elastic{es}, nil
}

// CreateIndex creates an index in Elasticsearch
func (es *Elastic) CreateIndex(indexName string) error {
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
func (es *Elastic) BulkIndex(indexName string, items []map[string]interface{}) error {
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
func (es *Elastic) IndexFromFile(indexName string, filePath string) error {
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
	err = es.CreateIndex(indexName)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	// Bulk index the documents
	err = es.BulkIndex(indexName, items)
	if err != nil {
		return fmt.Errorf("failed to bulk index documents: %w", err)
	}

	return nil
}

func ConstructQuery(filterParams models.FilterParams, sortOption, sortOrder string, page, size int64) (map[string]interface{}, error) {
	query := make(map[string]interface{})
	boolQuery := make(map[string]interface{})
	boolQuery["must"] = []interface{}{}

	for _, condition := range filterParams.Conditions {
		fieldName := models.GetFieldEnumValue(condition.FieldName)

		switch condition.Operation {
		case models.Equals:
			matchQuery := map[string]interface{}{
				"match": map[string]interface{}{
					fieldName: condition.Value,
				},
			}
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), matchQuery)
		case models.Lt:
			numericValue, err := strconv.Atoi(condition.Value)
			if err != nil {
				return nil, fmt.Errorf("error converting value to integer: %v", err)
			}
			rangeQuery := map[string]interface{}{
				"range": map[string]interface{}{
					fieldName: map[string]interface{}{
						"lt": numericValue,
					},
				},
			}
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), rangeQuery)
		case models.Gt:
			numericValue, err := strconv.Atoi(condition.Value)
			if err != nil {
				return nil, fmt.Errorf("error converting value to integer: %v", err)
			}
			rangeQuery := map[string]interface{}{
				"range": map[string]interface{}{
					fieldName: map[string]interface{}{
						"gt": numericValue,
					},
				},
			}
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), rangeQuery)
		case models.Query:
			matchQuery := map[string]interface{}{
				"wildcard": map[string]interface{}{
					fieldName: fmt.Sprintf("*%s*", condition.Value),
				},
			}
			boolQuery["must"] = append(boolQuery["must"].([]interface{}), matchQuery)
		default:
			return nil, fmt.Errorf("unsupported operation: %s", condition.Operation)
		}
	}

	query["query"] = map[string]interface{}{
		"bool": boolQuery,
	}

	// Add sorting
	if sortOption != "" && sortOrder != "" {
		sort := map[string]interface{}{
			sortOption: map[string]interface{}{"order": sortOrder},
		}
		query["sort"] = sort
	}

	query["from"] = page
	query["size"] = size

	return query, nil
}

func (es *Elastic) DoSearch(ctxBg context.Context, queryBytes []byte) models.SearchResponse {
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
