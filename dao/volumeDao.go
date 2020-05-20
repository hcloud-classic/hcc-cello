package dao

import (
	"hcc/cello/lib/logger"
	"hcc/cello/lib/mysql"
	"hcc/cello/model"
	"time"
)

func ReadVolume(args map[string]interface{}) (interface{}, error) {
	var volume model.Volume
	var err error
	uuid := args["uuid"].(string)

	var size int
	var filesystem string
	var serverUUID string
	var useType string
	var userUUID string
	var createdAt time.Time

	sql := "select * from volume where uuid = ?"
	err = mysql.Db.QueryRow(sql, uuid).Scan(
		&uuid,
		&size,
		&filesystem,
		&serverUUID,
		&useType,
		&userUUID,
		&createdAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	volume.UUID = uuid
	volume.Size = size
	volume.Filesystem = filesystem
	volume.ServerUUID = serverUUID
	volume.UseType = useType
	volume.UserUUID = userUUID
	volume.CreatedAt = createdAt

	return volume, nil
}

func checkReadVolumeListPageRow(args map[string]interface{}) bool {
	_, rowOk := args["row"].(int)
	_, pageOk := args["page"].(int)

	return !rowOk || !pageOk
}
