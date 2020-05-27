package dao

import (
	"hcc/cello/lib/logger"
	"hcc/cello/lib/mysql"
	"hcc/cello/model"
	"strconv"
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

func ReadVolumeList(args map[string]interface{}) (interface{}, error) {
	var err error
	var volumes []model.Volume
	var requestUUID string
	var createdAt time.Time

	size, sizeOk := args["size"].(int)
	filesystem, filesystemOk := args["filesystem"].(string)
	serverUUID, serverUUIDOk := args["server_uuid"].(string)
	useType, useTypeOk := args["use_type"].(string)
	userUUID, userUUIDOk := args["user_uuid"].(string)

	if !userUUIDOk {
		return nil, err
	}
	row, _ := args["row"].(int)
	page, _ := args["page"].(int)
	if checkReadVolumeListPageRow(args) {
		return nil, err
	}

	sql := "select * from volume where 1=1"
	if sizeOk {
		sql += " and size = " + strconv.Itoa(size)
	}
	if filesystemOk {
		sql += " and filesystem = '" + filesystem + "'"
	}
	if serverUUIDOk {
		sql += " and server_uuid = '" + serverUUID + "'"
	}
	if useTypeOk {
		sql += " and use_type = '" + useType + "'"
	}

	sql += " and user_uuid = ? order by created_at desc limit ? offset ?"

	stmt, err := mysql.Db.Query(sql, userUUID, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&requestUUID, &size, &filesystem, &serverUUID, &useType, &userUUID, &createdAt)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		volume := model.Volume{UUID: requestUUID, Size: size, Filesystem: filesystem, ServerUUID: serverUUID, UseType: useType, UserUUID: userUUID, CreatedAt: createdAt}
		logger.Logger.Println(volume)
		volumes = append(volumes, volume)
	}
	return volumes, nil
}

func ReadVolumeAll(args map[string]interface{}) (interface{}, error) {
	var err error
	var volumes []model.Volume
	var requestUUID string
	var size int
	var filesystem string
	var serverUUID string
	var useType string
	var userUUID string
	var createdAt time.Time
	row, rowOk := args["row"].(int)
	page, pageOk := args["page"].(int)
	if !rowOk || !pageOk {
		return nil, err
	}

	sql := "select * from volume order by created_at desc limit ? offset ?"

	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&requestUUID, &size, &filesystem, &serverUUID, &useType, &userUUID, &createdAt)

		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		volume := model.Volume{UUID: requestUUID, Size: size, Filesystem: filesystem, ServerUUID: serverUUID, UseType: useType, UserUUID: userUUID, CreatedAt: createdAt}
		volumes = append(volumes, volume)
	}

	return volumes, nil
}

func ReadVolumeNum() (model.VolumeNum, error) {
	var volumeNum model.VolumeNum
	var volumeNr int
	var err error

	sql := "select count(*) from volume"
	err = mysql.Db.QueryRow(sql).Scan(&volumeNr)
	if err != nil {
		logger.Logger.Println(err)
		return volumeNum, err
	}
	volumeNum.Number = volumeNr

	return volumeNum, nil
}
