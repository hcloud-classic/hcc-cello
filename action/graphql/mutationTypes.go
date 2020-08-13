package graphql

import (
	"hcc/cello/action/driver"
	graphqlType "hcc/cello/action/graphql/type"
	"hcc/cello/lib/logger"

	"github.com/graphql-go/graphql"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		"create_volume": &graphql.Field{
			Type:        graphqlType.VolumeType,
			Description: "Create new volume",
			Args: graphql.FieldConfigArgument{
				"size": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"filesystem": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"server_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"network_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"use_type": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"user_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"gateway_ip": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return driver.CreatePxeActionHandler(params)
			},
		},
		"update_volume": &graphql.Field{
			Type:        graphqlType.VolumeType,
			Description: "Update volume",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
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
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_volume")
				return dao.UpdateVolume(params.Args)
			},
		},
		"delete_volume": &graphql.Field{
			Type:        graphqlType.VolumeType,
			Description: "Delete volume by uuid",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: delete_volume")
				return dao.DeleteVolume(params.Args)
			},
		},
	},
	"create_volume_attachment": &graphql.Field{
		Type:        graphqlType.VolumeAttachmentType,
		Description: "Create new volume_attachment",
		Args: graphql.FieldConfigArgument{
			"volume_uuid": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
			"server_uuid": &graphql.ArgumentConfig{
				Type: graphql.String,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return dao.CreateVolumeAttachment(params.Args)
		},
	},
})
