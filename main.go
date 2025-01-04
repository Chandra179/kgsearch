package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func main() {
	app := fiber.New()
	driver := connectToNeo4j()

	defer driver.Close(context.Background())

	app.Get("/search", func(c *fiber.Ctx) error {
		category := c.Query("category")
		if category == "" {
			return c.Status(400).SendString("Category is required")
		}

		session := driver.NewSession(context.Background(), neo4j.SessionConfig{})
		defer session.Close(context.Background())

		query := `
			MATCH (p:Product)-[:BELONGS_TO]->(c:Category {name: $category})
			RETURN p.name AS name, p.price AS price, p.brand AS brand
		`
		result, err := session.ExecuteRead(context.Background(), func(tx neo4j.ManagedTransaction) (any, error) {
			records, err := tx.Run(context.Background(), query, map[string]any{"category": category})
			if err != nil {
				return nil, err
			}

			var products []map[string]any
			for records.Next(context.Background()) {
				// Get record values as []any
				values := records.Record().Values
				// Create a map for the record values
				product := map[string]any{
					"name":  values[0],
					"price": values[1],
					"brand": values[2],
				}
				// Append the product map to the products slice
				products = append(products, product)
			}
			if err := records.Err(); err != nil {
				return nil, err
			}
			return products, nil
		})
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(result)
	})

	log.Fatal(app.Listen(":8080"))
}

func connectToNeo4j() neo4j.DriverWithContext {
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7687",
		neo4j.BasicAuth("neo4j", os.Getenv("NEO4J_PASSWORD"), ""),
	)
	if err != nil {
		log.Fatalf("Failed to connect to Neo4j: %v", err)
	}
	return driver
}
