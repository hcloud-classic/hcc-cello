package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryTypes,
		Mutation: mutationTypes,
	},
)

var GraphqlHandler = handler.New(&handler.Config{
	Schema:   &Schema,
	Pretty:   true,
	GraphiQL: true,
})
