package driver

import (
	"errors"
	"fmt"
	"hcc/cello/dao"
	"hcc/cello/lib/handler"
	"hcc/cello/lib/logger"

	"github.com/graphql-go/graphql"
)

func CreatePxeActionHandler(params graphql.ResolveParams) (interface{}, error) {

	logger.Logger.Println("Resolving: create_volume")

	err := handler.ActionHandle(params.Args)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	volume, err := dao.CreateVolume(params.Args)
	if err != nil {
		strerr := "create_volume action status => (" + fmt.Sprintln(err) + ")"
		return nil, errors.New("[Cello]Can't Create Volume : " + strerr)
	}
	logger.Logger.Println("[Create Volume] Success : ", volume)
	return volume, nil
}
