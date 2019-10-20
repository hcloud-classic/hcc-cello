package graphql

import "github.com/graphql-go/graphql"

var volumeNum = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "VolumeNum",
		Fields: graphql.Fields{
			"number": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)
