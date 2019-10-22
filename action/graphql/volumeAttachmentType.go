<<<<<<< HEAD:action/graphql/volumeAttachmentType.go
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
=======
package graphql

import "github.com/graphql-go/graphql"

var volumeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Volume",
		Fields: graphql.Fields{
			"uuid": &graphql.Field{
				Type: graphql.String,
			},
			"size": &graphql.Field{
				Type: graphql.Int,
			},
			"filesystem": &graphql.Field{
				Type: graphql.String,
			},
			"server_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"use_type": &graphql.Field{
				Type: graphql.String,
			},
			"user_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
>>>>>>> 6d18aed (update architecture):action/graphql/volumeType.go
