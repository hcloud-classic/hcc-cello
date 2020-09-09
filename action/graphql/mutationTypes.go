package graphql

import (
	"fmt"
	graphqlType "hcc/cello/action/graphql/type"
	"hcc/cello/dao"
	"hcc/cello/driver"
	"hcc/cello/lib/formatter"
	"hcc/cello/lib/logger"

	"github.com/graphql-go/graphql"
)

var mutationTypes = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		// volume DB
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
				"lun_num": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"pool": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return driver.CreateVolActionHandler(params)
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
		// volume_attachment DB
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
		"update_volume_attachment": &graphql.Field{
			Type:        graphqlType.VolumeAttachmentType,
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
			Type:        graphqlType.VolumeAttachmentType,
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
		"testql": &graphql.Field{
			Type:        graphqlType.VolumeType,
			Description: "test None input val",
			Args: graphql.FieldConfigArgument{
				"page": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"row": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {

				// formatter.LoadDB()
				// handler.PrepareNReloadVolInfo()
				// qwe := formatter.New()
				// qwe.PutVal("codex")
				// fmt.Println("test1", qwe.GetVal("codex"))
				// asd := formatter.New()
				// fmt.Println("test2", asd.GetVal("codex"))
				// formatter.VolObjectMap.PutDomain("codex")
				// formatter.VolObjectMap.PutDomain("2020")
				// formatter.VolObjectMap.PutDomain("asd")
				// formatter.VolObjectMap.PutDomain("7890")
				for i, args := range formatter.GlobalVolumesDB {
					fmt.Println(i, " > ", args)
					// tpl, _ := i
					if formatter.GlobalVolumesDB[i].ServerUUID == "7890" {
						formatter.VolObjectMap.PutDomain(formatter.GlobalVolumesDB[i].ServerUUID)
						formatter.VolObjectMap.SetIscsiLun(formatter.GlobalVolumesDB[i], formatter.GlobalVolumesDB[i].Pool+"asd")
					} else {
						formatter.VolObjectMap.PutDomain(formatter.GlobalVolumesDB[i].ServerUUID)
					}

					// if tpl < len(formatter.GlobalVolumesDB) {
					// 	formatter.VolObjectMap.SetIscsiLun(formatter.GlobalVolumesDB[tpl], "/root/qwe/zxc")
					// }
				}
				// formatter.VolObjectMap.SetIscsiLun(formatter.GlobalVolumesDB[0], "/root/qwe/zxc")
				// formatter.VolObjectMap.SetIscsiLun(formatter.GlobalVolumesDB[1], "/root/zxc/aaaa")
				for _, args := range formatter.VolObjectMap.Domain {
					// tpl, _ := strconv.Atoi(i)
					// fmt.Println("Before", i, "<=>", args.TargetName, args.Lun[tpl].Path)
					for i, qwe := range args.Lun {
						// tpl, _ := strconv.Atoi(i)
						fmt.Println("before", i, "<=>", qwe, qwe.Path)
					}
				}

				volume := formatter.GlobalVolumesDB[0]
				formatter.VolObjectMap.RemoveIscsiLun(volume, 1)
				// formatter.VolObjectMap.RemoveDomain(volume.ServerUUID)
				for _, args := range formatter.VolObjectMap.Domain {
					for i, qwe := range args.Lun {
						// tpl, _ := strconv.Atoi(i)
						fmt.Println("After", i, "<=>", qwe, qwe.Path)
					}
				}

				// for i, args := range formatter.GlobalVolumesDB {
				// 	tpl := formatter.VolObjectMap.GetIscsiData(args.ServerUUID)
				// 	fmt.Println(i, "-", tpl)
				// }
				return dao.ReadVolumeAll(params.Args)
			},
		},
	},
})
