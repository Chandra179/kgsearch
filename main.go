package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j/config"
)

type Neo4jClient struct {
	driver neo4j.DriverWithContext
}

func NewNeo4jClient() (*Neo4jClient, error) {
	config := func(config *config.Config) {
		config.MaxConnectionPoolSize = 100           
		config.MaxConnectionLifetime = 1 * time.Hour
		config.ConnectionAcquisitionTimeout = 2 * time.Minute 
		config.ConnectionLivenessCheckTimeout = 2 * time.Second
	}
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://neo4j:7687",
		neo4j.NoAuth(),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}
	ctx := context.Background()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		driver.Close(ctx)
		return nil, fmt.Errorf("failed to verify connectivity: %w", err)
	}

	return &Neo4jClient{driver: driver}, nil
}

func (c *Neo4jClient) Close(ctx context.Context) error {
	return c.driver.Close(ctx)
}

// Example of a method using the connection pool
func (c *Neo4jClient) CreatePerson(ctx context.Context, name string, age int) error {
	session := c.driver.NewSession(ctx, neo4j.SessionConfig{
		DatabaseName: "neo4j",
		AccessMode:   neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	_, err := session.ExecuteWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx,
			"MERGE (p:Person {name: $name}) SET p.age = $age RETURN p",
			map[string]any{
				"name": name,
				"age":  age,
			})
		if err != nil {
			return nil, err
		}
		return result.Single(ctx)
	})
	return err
}

func main() {
	client, err := NewNeo4jClient()
	if err != nil {
		log.Fatalf("Failed to create Neo4j client: %v", err)
	}
	ctx := context.Background()
	defer client.Close(ctx)

	// Example of concurrent operations using the connection pool
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := fmt.Sprintf("Person-%d", i)
			err := client.CreatePerson(ctx, name, 25+i)
			if err != nil {
				log.Printf("Error creating person %s: %v", name, err)
				return
			}
			log.Printf("Successfully created %s", name)
		}(i)
	}
	wg.Wait()

	fmt.Println("All operations completed")
}