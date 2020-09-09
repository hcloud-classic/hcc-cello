package handler

import (
	"errors"
	"fmt"
	"hcc/cello/dao"
	"hcc/cello/lib/formatter"
	"hcc/cello/lib/logger"
	"hcc/cello/model"
)

// ReloadAllofVolInfo : Reload All of volume info
func ReloadAllofVolInfo() error {

	celloParams := make(map[string]interface{})
	celloParams["row"] = 254
	celloParams["page"] = 1
	dbVol, err := dao.ReadVolumeAll(celloParams)
	if err != nil {
		fmt.Println("Error")
	}
	formatter.GlobalVolumesDB = dbVol.([]model.Volume)
	formatter.VolObjectMap.PreLoad()
	// For Debug, What is in "formatter.GlobalVolumesDB"
	// logger.Logger.Println("Reload Volume DB")
	// for _, args := range formatter.GlobalVolumesDB {
	// 	qwe, _ := json.Marshal(args)
	// 	var obj map[string]interface{}
	// 	json.Unmarshal([]byte(qwe), &obj)
	// 	f := colorjson.NewFormatter()
	// 	f.Indent = 4
	// 	s, _ := f.Marshal(obj)
	// 	logger.Logger.Println(string(s))
	// }

	return err
}

//ActionHandle : action handler
func ActionHandle(args map[string]interface{}) error {

	volume := model.Volume{
		Size:       args["size"].(int),
		Filesystem: args["filesystem"].(string),
		ServerUUID: args["server_uuid"].(string),
		UseType:    args["use_type"].(string),
		UserUUID:   args["user_uuid"].(string),
		NetworkIP:  args["network_ip"].(string),
		GatewayIP:  args["gateway_ip"].(string),
	}

	if args["use_type"].(string) == "os" {
		logger.Logger.Println("ActionHandle: Creating OS volume")

		actionstatus, err := PreparePxeSetting(args["server_uuid"].(string), args["use_type"].(string), args["network_ip"].(string), args["gateway_ip"].(string))
		if !actionstatus {
			strerr := "create_volume action status=>actionstatus " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Prepare Setting (" + strerr + ")")
		}
		logger.Logger.Println("after ActionHandle")

		createstatus, err := CreateVolume(volume)
		if !createstatus {
			strerr := "create_volume action status=>createstatus " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Create Volume ( " + strerr + ")")
		}
		logger.Logger.Println("after CreateVolume")

		iscsistatus, err := WriteIscsiConfigObject(volume)
		if !iscsistatus {
			strerr := "create_volume action status=>iscsistatus " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Create Iscsi Setting ( " + strerr + ")")
		}
		logger.Logger.Println("[Action Result]  WriteIscsiConfigObject : ", actionstatus, " , CreateVolume : ", createstatus, "PrepareIscsiSetting : ", iscsistatus)

	}

	if args["use_type"].(string) == "data" {
		logger.Logger.Println("ActionHandle: Creating data volume")

		createstatus, err := CreateVolume(volume)
		if !createstatus {
			strerr := "create_volume action status=> " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Create Volume ( " + strerr + ")")
		}

		iscsistatus, err := WriteIscsiConfigObject(volume)
		if !iscsistatus {
			strerr := "create_volume action status=>iscsistatus " + fmt.Sprintln(err)
			return errors.New("[Cello]Can't Create Iscsi Setting ( " + strerr + ")")
		}

		logger.Logger.Println("[Action Result] : ", createstatus)
	}

	return nil
}

func updateVolData() {

}
