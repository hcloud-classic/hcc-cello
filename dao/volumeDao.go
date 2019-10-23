package dao

import (
	"hcc/cello/lib/logger"
	"hcc/cello/lib/mysql"
	"hcc/cello/lib/uuidgen"
	"hcc/cello/model"
	"strconv"
	"time"
)

// ReadVolume - cgs
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

// ReadVolumeList - cgs
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
	row, rowOk := args["row"].(int)
	page, pageOk := args["page"].(int)
	if !rowOk || !pageOk {
		return nil, err
	}

	sql := "select * from volume where 1 = 1 and "
	if sizeOk {
		sql += " size = '" + strconv.Itoa(size) + "'"
		if filesystemOk || serverUUIDOk || useTypeOk {
			sql += " and"
		}
	}
	if filesystemOk {
		sql += " filesystem = '" + filesystem + "'"
		if serverUUIDOk || useTypeOk {
			sql += " and"
		}
	}
	if serverUUIDOk {
		sql += " server_uuid = '" + serverUUID + "'"
		if useTypeOk {
			sql += " and"
		}
	}
	if useTypeOk {
		sql += " use_type = '" + useType + "' and"
	}

	sql += " user_uuid = ? order by created_at desc limit ? offset ?"

	stmt, err := mysql.Db.Query(sql, userUUID, row, row*(page-1))
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

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

// ReadVolumeAll - cgs
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
	defer stmt.Close()

	for stmt.Next() {
		//err := stmt.Scan(&uuid, &subnetUUID, &os, &serverName, &serverDesc, &cpu, &memory, &diskSize, &status, &userUUID, &createdAt)
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

// ReadVolumeNum - cgs
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

// CreateVolume - cgs
func CreateVolume(args map[string]interface{}) (interface{}, error) {
	uuid, err := uuidgen.UUIDgen()
	if err != nil {
		logger.Logger.Println("Failed to generate uuid!")
		return nil, err
	}

	volume := model.Volume{
		UUID:       uuid,
		Size:       args["size"].(int),
		Filesystem: args["filesystem"].(string),
		ServerUUID: args["server_uuid"].(string),
		UseType:    args["use_type"].(string),
		UserUUID:   args["user_uuid"].(string),
	}

	sql := "insert into volume(uuid, size, filesystem, server_uuid, use_type, user_uuid, created_at) values (?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err := stmt.Exec(volume.UUID, volume.Size, volume.Filesystem, volume.ServerUUID, volume.UseType, volume.UserUUID)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return volume, nil
}

// UpdateVolume - cgs
func UpdateVolume(args map[string]interface{}) (interface{}, error) {
	var err error

	requestedUUID, requestedUUIDOk := args["uuid"].(string)
	size, sizeOk := args["size"].(int)
	filesystem, filesystemOk := args["filesystem"].(string)
	serverUUID, serverUUIDOk := args["server_uuid"].(string)
	useType, useTypeOk := args["use_type"].(string)
	userUUID, userUUIDOk := args["user_uuid"].(string)

	volume := new(model.Volume)
	volume.UUID = requestedUUID
	volume.Size = size
	volume.Filesystem = filesystem
	volume.ServerUUID = serverUUID
	volume.UseType = useType
	volume.UserUUID = userUUID

	if requestedUUIDOk {

		if !sizeOk && !filesystemOk && !serverUUIDOk && !useTypeOk {
			return nil, nil
		}

		sql := "update volume set"
		if sizeOk {
			sql += " size = '" + strconv.Itoa(volume.Size) + "'"
			if filesystemOk || serverUUIDOk || useTypeOk || userUUIDOk {
				sql += ", "
			}
		}
		if filesystemOk {
			sql += " filesystem = '" + volume.Filesystem + "'"
			if serverUUIDOk || useTypeOk || userUUIDOk {
				sql += ", "
			}
		}
		if serverUUIDOk {
			sql += " server_uuid = '" + volume.ServerUUID + "'"
			if useTypeOk || userUUIDOk {
				sql += ", "
			}
		}
		if useTypeOk {
			sql += " use_type = '" + volume.UseType + "'"
			if userUUIDOk {
				sql += ", "
			}
		}
		if userUUIDOk {
			sql += " user_uuid = " + volume.UserUUID
		}
		sql += " where uuid = ?"

		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()

		result, err2 := stmt.Exec(volume.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.LastInsertId())
		return volume, nil
	}

	return nil, err
}

// DeleteVolume - cgs
func DeleteVolume(args map[string]interface{}) (interface{}, error) {
	var err error

	requestedUUID, ok := args["uuid"].(string)
	if ok {
		sql := "delete from volume where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()
		result, err2 := stmt.Exec(requestedUUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.RowsAffected())

		return requestedUUID, nil
	}

	return requestedUUID, err
}
