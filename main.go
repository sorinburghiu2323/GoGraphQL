package main

import (
	"encoding/json"
	"fmt"
	"log"
	"goGraphQL/utils"

	"github.com/graphql-go/graphql"
)


func populate() []dataTypes.Tutorial {
	author := &dataTypes.Author{Name: "Elliot Forbes", Tutorials: []int{1}}
	tutorial := dataTypes.Tutorial{
		ID: 1,
		Title: "Go tutorial",
		Author: *author,
		Comments: []dataTypes.Comment{
			{Body: "First!"},
		},
	}

	var tutorials []dataTypes.Tutorial
	tutorials = append(tutorials, tutorial)
	return tutorials
}


func main() {
	fmt.Println("Running GraphQL application")

	tutorials := populate()

	// Initialize objects.
	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)
	var authorType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{
				"Name": &graphql.Field{
					Type: graphql.String,
				},
				"Tutorials": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)
	var tutorialType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Tutorial",
			Fields: graphql.Fields{
				"ID": &graphql.Field{
					Type: graphql.Int,
				},
				"Title": &graphql.Field{
					Type: graphql.String,
				},
				"Author": &graphql.Field{
					Type: authorType,
				},
				"Comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	// Create a new GraphQL schema.
	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type: tutorialType,
			Description: "Get Tutorial by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int)
				if ok {
					for _, tutorial := range tutorials {
						if tutorial.ID == id {
							return tutorial, nil
						}
					}
				}
				return "world", nil
			},
		},
		"list": &graphql.Field{
			Type: graphql.NewList(tutorialType),
			Description: "Get all Tutorials",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}

	// Define the object configuration.
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}

	// Define a schema configuration.
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}

	// Create the schema.
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create GraphQL Schema, err: %v", err)
	}

	// Write and execute sample query.
	query := `
		{
			tutorial(id: 1) {
				Title
				Author {
					Name
					Tutorials
				}
			}
		}
	`
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute GraphQL operation, err: %+v", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s \n", rJSON)
}
