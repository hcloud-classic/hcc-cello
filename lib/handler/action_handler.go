package handler

import (
	"fmt"
	"hcc/cello/dao"
	"hcc/cello/lib/formatter"
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
	// formatter.VolObjectMap.PreLoad()
	return err
}
