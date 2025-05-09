package config

import (
	"log"

	"github.com/olivere/elastic/v7"
)

func NewClient() *elastic.Client {
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
	)
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %v", err)
	}
	return client
}
