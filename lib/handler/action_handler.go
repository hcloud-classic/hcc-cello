package handler

import (
	"errors"
	"fmt"
)

//ActionHandle : action handler
func ActionHandle(args map[string]interface{}) (bool, interface{}) {

	if args["use_type"].(string) == "os" {
		actionstatus, err := PreparePxeSetting(args["server_uuid"].(string), args["use_type"].(string), args["network_ip"].(string))
		if !actionstatus {
			strerr := "create_volume action status=>actionstatus " + fmt.Sprintln(err)
			return false, errors.New("[Cello]Can't Create Volume in false: " + strerr)
		}
		createstatus, err := CreateVolume(args["filesystem"].(string), args["server_uuid"].(string), args["use_type"].(string), args["size"].(int))
		if !createstatus {
			strerr := "create_volume action status=>createstatus " + fmt.Sprintln(err)
			return false, errors.New("[Cello]Can't Create Volume in false: " + strerr)
		}

	}

	if args["use_type"].(string) == "data" {
		createstatus, err := CreateVolume(args["filesystem"].(string), args["server_uuid"].(string), args["use_type"].(string), args["size"].(int))
		if !createstatus {
			strerr := "create_volume action status=> " + fmt.Sprintln(err)
			return false, errors.New("[Cello]Can't Create Volume in true: " + strerr)
		}

	}
	return true, "Complete Action"
}
