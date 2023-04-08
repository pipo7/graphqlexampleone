package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

type Tutorial struct {
	Title    string
	Author   Author
	ID       int
	Comments []Comment
}
type Author struct {
	Name      string
	Tutorials []int
}
type Comment struct {
	Body string
}

// Manually poplulate the structs
func populate() []Tutorial {

	// 1st item in list
	author := &Author{Name: "PSJohn", Tutorials: []int{1}}
	tutorial1 := Tutorial{
		ID:     1,
		Title:  "Magic Covers",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "First review comment"},
		},
	}
	// 2nd item in list
	author = &Author{Name: "JKR", Tutorials: []int{1}}
	tutorial2 := Tutorial{
		ID:     2,
		Title:  "Harry Potter Covers",
		Author: *author,
		Comments: []Comment{
			{Body: "Second review comment"},
			{Body: "Third review comment"},
		},
	}
	//Add all tutorial into tutorials slice
	var tutorials []Tutorial
	tutorials = append(tutorials, tutorial1)
	tutorials = append(tutorials, tutorial2)
	fmt.Println("Tutuorials:", tutorials)
	return tutorials
}

func main() {
	fmt.Println("Simple graphql example2")
	tutorials := populate()
	// define the fields first
	var commentType = graphql.NewObject( // define object
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{ // object has name and fields
				"body": &graphql.Field{ // Field object and its type
					Type: graphql.String,
				},
			},
		},
	)

	var authorType = graphql.NewObject( // define object
		graphql.ObjectConfig{
			Name: "Author",
			Fields: graphql.Fields{ // object has name and fields
				"name": &graphql.Field{ // Field object and its type
					Type: graphql.String,
				},
				"tutorials": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)

	var tutorialtype = graphql.NewObject( // define object
		graphql.ObjectConfig{
			Name: "Tutorial",
			Fields: graphql.Fields{ // object has name and fields
				"id": &graphql.Field{ // Field object and its type
					Type: graphql.String,
				},
				"title": &graphql.Field{ // Field object and its type
					Type: graphql.String,
				},
				"author": &graphql.Field{
					Type: authorType, // defined above
				},
				"comments": &graphql.Field{
					Type: graphql.NewList(commentType),
				},
			},
		},
	)

	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorialtype,
			Description: "get tutorial by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["id"].(int) // parse the argument id passed from query
				if ok {
					for _, tutorial := range tutorials {
						if int(tutorial.ID) == id { // if id matches to tutorial.ID then return it
							return tutorial, nil
						}
					}
				}
				return nil, nil
			},
		},

		"list": &graphql.Field{
			Type:        graphql.NewList(tutorialtype),
			Description: "Returns the list of tutorials",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}
	// define the rootQuety and then schemaConfig which takes rootQuery as input
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new schema , err %v", err)
	}

	// Lets try to query list of tutorials and retunr id, title and comments
	queryList := `
	{
		list{
			id
			title
		} 
	}
	`

	params := graphql.Params{Schema: schema, RequestString: queryList}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute the query: %v/n", r.Errors)
	}
	rJSON, _ := json.Marshal(r)
	fmt.Printf("%s\n", rJSON)

	// get comments and its sub element body too
	queryList = `
	{
		list{
			id
			title
			comments {
				body
			}
		} 
	}
	`

	params = graphql.Params{Schema: schema, RequestString: queryList}
	r = graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute the query: %v/n", r.Errors)
	}
	rJSON, _ = json.Marshal(r)
	fmt.Printf("%s\n", rJSON)

	// Lets try to query a tutorial by id
	query := `
	{
		tutorial(id:1){
			title
		} 
	}
	`
	params = graphql.Params{Schema: schema, RequestString: query}
	r = graphql.Do(params)
	if len(r.Errors) > 0 {
		log.Fatalf("Failed to execute the query: %v/n", r.Errors)
	}
	rJSON, _ = json.Marshal(r)
	fmt.Printf("%s\n", rJSON)

}
