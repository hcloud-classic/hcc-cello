package graphqlType

import "github.com/graphql-go/graphql"

var VolumeNum = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "VolumeNum",
		Fields: graphql.Fields{
			"number": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
