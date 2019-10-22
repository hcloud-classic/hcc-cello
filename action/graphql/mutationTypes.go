package graphql

import (
	"github.com/graphql-go/graphql"
	"hcc/cello/dao"
	"hcc/cello/lib/logger"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// volume DB
		"create_volume": &graphql.Field{
			Type:        volumeType,
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
				"use_type": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"user_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return dao.CreateVolume(params.Args)
			},
		},
		"update_volume": &graphql.Field{
			Type:        volumeType,
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
			Type:        volumeType,
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
		// volume_attachment DB
		"create_volume_attachment": &graphql.Field{
			Type:        volumeAttachmentType,
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
		"update_volume_attachment": &graphql.Field{
			Type:        volumeAttachmentType,
			Description: "Update volume",
			Args: graphql.FieldConfigArgument{
				"volume_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"server_uuid": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: update_volume_attachment")
				return dao.UpdateVolumeAttachment(params.Args)
			},
		},
		"delete_volume_attachment": &graphql.Field{
			Type:        volumeAttachmentType,
			Description: "Delete volume_attachment by uuid",
			Args: graphql.FieldConfigArgument{
				"uuid": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				logger.Logger.Println("Resolving: delete_volume_attachment")
				return dao.DeleteVolumeAttachment(params.Args)
			},
		},
	},
})
