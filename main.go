package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/graphql-go/graphql"
)

//Example to query Tutorials list and a tutorial by id using GraphQL

// Tutorial
type Tutorial struct {
	Title    string
	Author   Author
	ID       int
	Comments []Comment
}

// Author
type Author struct {
	Name      string
	Tutorials []int
}

// Comment
type Comment struct {
	Body string
}

// Manually poplulate the structs with 2 tutorials ... You can later use API for this
func populate() []Tutorial {

	// 1st item in Tutorial list
	author := &Author{Name: "PSJohn", Tutorials: []int{1}}
	tutorial1 := Tutorial{
		ID:     1,
		Title:  "Magic Covers",
		Author: *author,
		Comments: []Comment{
			Comment{Body: "First review comment"},
		},
	}
	// 2nd item in Tutorial list
	author = &Author{Name: "JK. Rowling", Tutorials: []int{1}}
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
	fmt.Println("Simple graphql example")

	// get the value of tutorials list
	tutorials := populate()

	// define the fields first ... this should match as per the Struct
	// So here we map the Go Struct to the graphQL object...
	var commentType = graphql.NewObject( // define object
		graphql.ObjectConfig{
			Name: "Comment", // Object Name
			Fields: graphql.Fields{ // object has fields and name and type
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	// define the fields first ... this should match as per the Struct
	// So here we map the Go Struct to the graphQL object...
	var authorType = graphql.NewObject( // define object
		graphql.ObjectConfig{
			Name: "Author", // object Name
			Fields: graphql.Fields{ // object has fields and name and type
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"tutorials": &graphql.Field{
					Type: graphql.NewList(graphql.Int),
				},
			},
		},
	)

	// define the fields first ... this should match as per the Struct
	// So here we map the Go Struct to the graphQL object...
	var tutorialtype = graphql.NewObject( // define object
		graphql.ObjectConfig{
			Name: "Tutorial", //object Name
			Fields: graphql.Fields{ // object has fields and name and type
				"id": &graphql.Field{
					Type: graphql.String,
				},
				"title": &graphql.Field{
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

	//now we define the fields on basis of which we will QUERY
	fields := graphql.Fields{
		"tutorial": &graphql.Field{
			Type:        tutorialtype, // Type of field
			Description: "get tutorial by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{ // Argument of field which is used while querying
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
			Type:        graphql.NewList(tutorialtype), // Type of field
			Description: "Returns the list of tutorials",
			// Argument if needed for field which is used while querying...since its LIST to return everything so we need no arguments here.
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return tutorials, nil
			},
		},
	}
	// define the rootQuety and then schemaConfig which takes rootQuery as input and then the scehma
	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new schema , err %v", err)
	}

	//****** TESTING *******
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
