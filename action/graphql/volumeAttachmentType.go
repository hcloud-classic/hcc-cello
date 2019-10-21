package graphql

import "github.com/graphql-go/graphql"

var volumeAttachmentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "VolumeAttachment",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"volume_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"server_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
			"updated_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
