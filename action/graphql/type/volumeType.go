package graphqlType

import "github.com/graphql-go/graphql"

var VolumeType = graphql.NewObject(
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
			"network_ip": &graphql.Field{
				Type: graphql.String,
			},
			"use_type": &graphql.Field{
				Type: graphql.String,
			},
			"user_uuid": &graphql.Field{
				Type: graphql.String,
			},
			"lun_num": &graphql.Field{
				Type: graphql.Int,
			},
			"pool": &graphql.Field{
				Type: graphql.Int,
			},
			"created_at": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
