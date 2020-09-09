package driver

import (
	"errors"
	"fmt"
	"hcc/cello/dao"
	"hcc/cello/lib/handler"
	"hcc/cello/lib/logger"

	"github.com/graphql-go/graphql"
)

//CreateVolActionHandler : Only Create volume action,
func CreateVolActionHandler(params graphql.ResolveParams) (interface{}, error) {

	logger.Logger.Println("Resolving: create_volume")

	err := handler.ActionHandle(params.Args)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	volume, err := dao.CreateVolume(params.Args)
	if err != nil {
		strerr := "create_volume action status => (" + fmt.Sprintln(err) + ")"
		logger.Logger.Println(strerr)
		return nil, errors.New("[Cello]Can't Create Volume : " + strerr)
	}
	// To-do
	// update in memory volume data

	//cross check
	err = handler.ReloadAllofVolInfo()
	if err != nil {
		strerr := "Can't reload DB  => (" + fmt.Sprintln(err) + ")"
		return nil, errors.New("[Cello]Can't Create Volume : " + strerr)
	}

	logger.Logger.Println("[Create Volume] Success : ", volume)
	return volume, nil
}

//UpdateVolActionHandler : Update Volume
func UpdateVolActionHandler(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}

//DeleteVolActionHandler : Delete Volume
func DeleteVolActionHandler(params graphql.ResolveParams) (interface{}, error) {
	return nil, nil
}
