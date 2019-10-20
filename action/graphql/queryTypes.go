package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/cello/dao"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
)

var queryTypes = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"volume": &graphql.Field{
				Type:        volumeType,
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
			"list_volume": &graphql.Field{
				Type:        graphql.NewList(volumeType),
				Description: "Get volume list",
				Args: graphql.FieldConfigArgument{
					"uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"size": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"filesystem": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"server_uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"use_type": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"user_uuid": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: list_volume")
					return dao.ReadVolumeList(params.Args)
				},
			},
			"all_volume": &graphql.Field{
				Type:        graphql.NewList(volumeType),
				Description: "Get all volume list",
				Args: graphql.FieldConfigArgument{
					"row": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
					"page": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: all_volume")
					return dao.ReadVolumeAll(params.Args)
				},
			},
			"num_volume": &graphql.Field{
				Type:        volumeNum,
				Description: "Get the number of volume",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					logger.Logger.Println("Resolving: num_volume")
					var volumeNum model.VolumeNum
					var err error
					volumeNum, err = dao.ReadVolumeNum()

					return volumeNum, err
				},
			},
		},
	})
