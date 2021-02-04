package driver

import (
	"github.com/graphql-go/graphql"
)

//Deprecate
//CreateVolActionHandler : Only Create volume action,
func CreateVolActionHandler(params graphql.ResolveParams) (interface{}, error) {
	// logger.Logger.Println("Resolving: create_volume")
	// if params.Args["use_type"] == "" || params.Args["server_uuid"] == "" {
	// 	strerr := "Param args missing => (" + ")"
	// 	return nil, errors.New("[Cello]Can't Create Volume : " + strerr)
	// }

	// err := handler.ReloadPoolObject()
	// if err != nil {
	// 	strerr := "Can't reload Zpool  => (" + fmt.Sprintln(err) + ")"
	// 	return nil, errors.New("[Cello]Can't Create Volume : " + strerr)
	// }

	// out, err := gouuid.NewV4()
	// if err != nil {
	// 	logger.Logger.Println("[volumeDao]Can't Create Volume UUID : ", err)
	// 	return nil, err
	// }
	// uuid := out.String()

	// params.Args["uuid"] = uuid

	// tempStruct, err := handler.ActionHandle(params.Args)
	// tempvolume := tempStruct.(model.Volume)
	// if err != nil {
	// 	logger.Logger.Println(err)
	// 	return nil, err
	// }
	// params.Args["lun_num"] = tempvolume.LunNum
	// volume, err := dao.CreateVolume(params.Args)
	// if err != nil {
	// 	strerr := "create_volume action status => (" + fmt.Sprintln(err) + ")"
	// 	logger.Logger.Println(strerr)
	// 	return nil, errors.New("[Cello]Can't Create Volume : " + strerr)
	// }

	// logger.Logger.Println("[Create Volume] Success : ", volume)
	return nil, nil
}

//UpdateVolActionHandler : Update Volume
func UpdateVolActionHandler(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

//DeleteVolActionHandler : Delete Volume
func DeleteVolActionHandler(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}
