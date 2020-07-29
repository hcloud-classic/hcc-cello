package handler

import (
	"errors"
	"fmt"
	"hcc/cello/lib/logger"
)

//ActionHandle : action handler
func ActionHandle(args map[string]interface{}) error {

	if args["use_type"].(string) == "os" {
		logger.Logger.Println("ActionHandle: Creating OS volume")

		actionstatus, err := PreparePxeSetting(args["server_uuid"].(string), args["use_type"].(string), args["network_ip"].(string), args["gateway_ip"].(string))
		if !actionstatus {
			strerr := "create_volume action status=>actionstatus " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Prepare Setting (" + strerr + ")")
		}

		createstatus, err := CreateVolume(args["filesystem"].(string), args["server_uuid"].(string), args["use_type"].(string), args["size"].(int))
		if !createstatus {
			strerr := "create_volume action status=>createstatus " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Create Volume ( " + strerr + ")")
		}

		iscsistatus, err := PrepareIscsiSetting(args["server_uuid"].(string), args["filesystem"].(string), args["use_type"].(string), args["size"].(int))
		if !iscsistatus {
			strerr := "create_volume action status=>iscsistatus " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Create Iscsi Setting ( " + strerr + ")")
		}
		logger.Logger.Println("[Action Result]  PreparePxeSetting : ", actionstatus, " , CreateVolume : ", createstatus, "PrepareIscsiSetting : ", iscsistatus)

	}

	if args["use_type"].(string) == "data" {
		logger.Logger.Println("ActionHandle: Creating data volume")

		createstatus, err := CreateVolume(args["filesystem"].(string), args["server_uuid"].(string), args["use_type"].(string), args["size"].(int))
		if !createstatus {
			strerr := "create_volume action status=> " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Create Volume ( " + strerr + ")")
		}
		logger.Logger.Println("[Action Result] : ", createstatus)
	}

	return nil
}
