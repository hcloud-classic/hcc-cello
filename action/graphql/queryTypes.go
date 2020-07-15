package graphql

import (
	graphqlType "hcc/cello/action/graphql/type"
	"hcc/cello/dao"
	"hcc/cello/lib/logger"
	"hcc/cello/model"

	"github.com/graphql-go/graphql"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"volume": &graphql.Field{
				Type:        graphqlType.VolumeType,
				Description: "Get volume by uuid",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: volume")
					return dao.ReadVolume(params.Args)
				},
			},
	})
