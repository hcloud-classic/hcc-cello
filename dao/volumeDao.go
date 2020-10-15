package dao

import (
	"errors"
	hccerr "hcc/cello/lib/errors"
	"hcc/cello/lib/logger"
	"hcc/cello/lib/mysql"
	"hcc/cello/model"
	"strconv"
	"time"
)

// ReadVolume : Single Volume info
func ReadVolume(in *model.Volume) (interface{}, error) {
	var volume model.Volume
	var err error
	uuid := in.UUID

	var size int
	var filesystem string
	var serverUUID string
	var useType string
	var userUUID string
	var lunNum int
	var pool string
	var createdAt time.Time

	sql := "select * from volume where uuid = ?"
	err = mysql.Db.QueryRow(sql, uuid).Scan(
		&uuid,
		&size,
		&filesystem,
		&serverUUID,
		&useType,
		&userUUID,
		&lunNum,
		&pool,
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
	volume.LunNum = lunNum
	volume.Pool = pool
	volume.CreatedAt = createdAt

	return volume, nil
}

func checkVolumePageRow(args map[string]interface{}) bool {
	_, rowOk := args["row"].(int)
	_, pageOk := args["page"].(int)

	return !rowOk || !pageOk
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

	row, _ := args["row"].(int)
	page, _ := args["page"].(int)
	if checkVolumePageRow(args) {
		return nil, errors.New("need row and page arguments")
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
	if userUUIDOk {
		sql += " and user_uuid = '" + userUUID + "'"
	}

	sql += " order by created_at desc limit ? offset ?"

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
	var lunNum int
	var pool string
	var createdAt time.Time
	row, _ := args["row"].(int)
	page, _ := args["page"].(int)
	if checkVolumePageRow(args) {
		return nil, errors.New("need row and page arguments")
	}

	sql := "select * from cello.volume order by created_at desc limit ? offset ?"

	stmt, err := mysql.Db.Query(sql, row, row*(page-1))
	// stmt, err := mysql.Db.Query(sql, row, 0)

	if err != nil {

		logger.Logger.Println(err.Error())
		return nil, err
	}

	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		//err := stmt.Scan(&uuid, &subnetUUID, &os, &serverName, &serverDesc, &cpu, &memory, &diskSize, &status, &userUUID, &createdAt)
		err := stmt.Scan(&requestUUID, &size, &filesystem, &serverUUID, &useType, &userUUID, &lunNum, &pool, &createdAt)

		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		volume := model.Volume{UUID: requestUUID, Size: size, Filesystem: filesystem, ServerUUID: serverUUID, UseType: useType, UserUUID: userUUID, LunNum: lunNum, Pool: pool, CreatedAt: createdAt}
		volumes = append(volumes, volume)
	}

	return volumes, nil
}

// ReadVolumeNum : The number of Volumes
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
func CreateVolume(in *model.Volume) (uint64, string) {
	sql := "insert into volume(uuid, size, filesystem, server_uuid, use_type, user_uuid,lun_num , pool,created_at) values (?, ?, ?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(in.UUID, in.Size, in.Filesystem, in.ServerUUID, in.UseType, in.UserUUID, in.LunNum, in.Pool)
	if err != nil {
		errStr := "[volumeDao]Can't Update DB: " + err.Error()
		return hccerr.CelloSQLOperationFail, errStr
	}
	return 0, ""
}

func checkUpdateVolumeArgs(args map[string]interface{}) bool {
	_, sizeOk := args["size"].(int)
	_, filesystemOk := args["filesystem"].(string)
	_, serverUUIDOk := args["server_uuid"].(string)
	_, useTypeOk := args["use_type"].(string)

	return !sizeOk && !filesystemOk && !serverUUIDOk && !useTypeOk
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
		if checkUpdateVolumeArgs(args) {
			return nil, errors.New("need some arguments")
		}

		sql := "update volume set"
		var updateSet = ""
		if sizeOk {
			updateSet += " size = " + strconv.Itoa(volume.Size) + ", "
		}
		if filesystemOk {
			updateSet += " filesystem = '" + volume.Filesystem + "', "
		}
		if serverUUIDOk {
			updateSet += " server_uuid = '" + volume.ServerUUID + "', "
		}
		if useTypeOk {
			updateSet += " use_type = '" + volume.UseType + "', "
		}
		if userUUIDOk {
			updateSet += " user_uuid = '" + volume.UserUUID + "', "
		}
		sql += updateSet[0:len(updateSet)-2] + " where uuid = ?"

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
func DeleteVolume(in *model.Volume) (uint64, error) {
	var err error

	if in.UUID != "" {
		sql := "delete from volume where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return hccerr.CelloSQLOperationFail, err
		}
		defer func() {
			_ = stmt.Close()
		}()
		result, err2 := stmt.Exec(in.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return hccerr.CelloSQLOperationFail, err
		}
		logger.Logger.Println(result.RowsAffected())

		return 0, nil
	}

	return hccerr.CelloSQLOperationFail, err
}
