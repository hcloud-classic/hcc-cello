package dao

import (
	"database/sql"
	"errors"
	"hcc/cello/lib/logger"
	"hcc/cello/lib/mysql"
	"hcc/cello/model"
	"strconv"
	"time"

	"innogrid.com/hcloud-classic/hcc_errors"
)

// ReadVolume : Single Volume info
func ReadVolume(in *model.Volume) (model.Volume, uint64, string) {
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
		return volume, hcc_errors.CelloSQLOperationFail, err.Error()
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

	return volume, 0, ""
}

func checkVolumePageRow(args map[string]interface{}) bool {
	_, rowOk := args["row"].(int)
	_, pageOk := args["page"].(int)

	return !rowOk || !pageOk
}

// ReadVolumeList - cgs
func ReadVolumeList(in *model.Volume, row int, page int) ([]model.Volume, uint64, string) {
	var err error
	var volumes []model.Volume
	var requestUUID string
	var createdAt time.Time

	var size int
	var filesystem string
	var serverUUID string
	var useType string
	var userUUID string
	var lunNum int
	var pool string

	if in.ServerUUID != "" {
		serverUUID = in.ServerUUID
	}
	if in.UserUUID != "" {
		userUUID = in.UserUUID
	}
	// sql := "select * from volume "
	sql := "select * from volume where 1=1"
	// if in.Size != "" {
	// sql += " and size = '1'"
	// }
	// if filesystemOk {
	// 	sql += " and filesystem = '" + filesystem + "'"
	// }
	// if serverUUIDOk {
	sql += " and server_uuid = '" + serverUUID + "'"
	// }
	// if useTypeOk {
	// 	sql += " and use_type = '" + useType + "'"
	// }
	// if userUUIDOk {
	// sql += " and user_uuid = '" + userUUID + "'"
	// }

	sql += " order by created_at asc limit ? offset ?"
	// sql += " order by created_at "

	stmt, err := mysql.Db.Query(sql, row, 0)
	// stmt, err := mysql.Db.Query(sql)

	// logger.Logger.Println(err.Error())
	if err != nil {
		logger.Logger.Println(err)
		return volumes, hcc_errors.CelloSQLOperationFail, err.Error()
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {

		err := stmt.Scan(
			&requestUUID,
			&size,
			&filesystem,
			&serverUUID,
			&useType,
			&userUUID,
			&lunNum,
			&pool,
			&createdAt)

		if err != nil {
			logger.Logger.Println(sql, err.Error())
			return volumes, hcc_errors.CelloSQLOperationFail, err.Error()
		}
		volume := model.Volume{
			UUID:       requestUUID,
			Size:       size,
			Filesystem: filesystem,
			ServerUUID: serverUUID,
			UseType:    useType,
			UserUUID:   userUUID,
			CreatedAt:  createdAt,
			Pool:       pool,
			LunNum:     lunNum,
		}
		logger.Logger.Println(volume)
		volumes = append(volumes, volume)
	}
	return volumes, 0, ""
}

// ReadVolumeAll - cgs
func ReadVolumeAll(args map[string]interface{}) (interface{}, uint64, string) {
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
	var stmt *sql.Rows
	row, _ := args["row"].(int)
	page, _ := args["page"].(int)
	if checkVolumePageRow(args) {
		return nil, hcc_errors.CelloGrpcArgumentError, errors.New("need row and page arguments").Error()
	}

	if row == 0 && page == 0 {
		sql := "select * from cello.volume order by created_at desc"
		stmt, err = mysql.Db.Query(sql)
	} else {
		sql := "select * from cello.volume order by created_at desc limit ? offset ?"
		stmt, err = mysql.Db.Query(sql, row, row*(page-1))
	}
	// stmt, err := mysql.Db.Query(sql, row, 0)

	if err != nil {

		logger.Logger.Println(err.Error())
		goto ERROR
	}

	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		//err := stmt.Scan(&uuid, &subnetUUID, &os, &serverName, &serverDesc, &cpu, &memory, &diskSize, &status, &userUUID, &createdAt)
		err := stmt.Scan(&requestUUID, &size, &filesystem, &serverUUID, &useType, &userUUID, &lunNum, &pool, &createdAt)

		if err != nil {
			logger.Logger.Println(err)
			goto ERROR

		}
		volume := model.Volume{UUID: requestUUID, Size: size, Filesystem: filesystem, ServerUUID: serverUUID, UseType: useType, UserUUID: userUUID, LunNum: lunNum, Pool: pool, CreatedAt: createdAt}
		volumes = append(volumes, volume)
	}
	return volumes, 0, ""

ERROR:
	return nil, hcc_errors.CelloSQLOperationFail, err.Error()
}

// ReadVolumeNum : The number of Volumes
func ReadVolumeNum() (model.VolumeNum, uint64, string) {
	var volumeNum model.VolumeNum
	var volumeNr int
	var err error

	sql := "select count(*) from volume"
	err = mysql.Db.QueryRow(sql).Scan(&volumeNr)
	if err != nil {
		logger.Logger.Println(err)
		return volumeNum, hcc_errors.CelloSQLOperationFail, err.Error()
	}
	volumeNum.Number = volumeNr

	return volumeNum, 0, ""
}

// CreateVolume - cgs
func CreateVolume(in *model.Volume) (uint64, string) {
	sql := "insert into volume(uuid, size, filesystem, server_uuid, use_type, user_uuid,lun_num , pool,created_at) values (?, ?, ?, ?, ?, ?, ?, ?, now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return hcc_errors.CelloSQLOperationFail, err.Error()
	}
	defer func() {
		_ = stmt.Close()
	}()
	_, err = stmt.Exec(in.UUID, in.Size, in.Filesystem, in.ServerUUID, in.UseType, in.UserUUID, in.LunNum, in.Pool)
	if err != nil {
		errStr := "[volumeDao]Can't Update DB: " + err.Error()
		return hcc_errors.CelloSQLOperationFail, errStr
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
func UpdateVolume(args map[string]interface{}) (interface{}, uint64, string) {
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
			return nil, hcc_errors.CelloGrpcArgumentError, errors.New("need some arguments").Error()
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
			goto ERROR
		}
		defer func() {
			_ = stmt.Close()
		}()

		result, err2 := stmt.Exec(volume.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			goto ERROR
		}
		logger.Logger.Println(result.LastInsertId())
		return volume, 0, ""

	}

ERROR:
	return nil, hcc_errors.CelloSQLOperationFail, err.Error()
}

// DeleteVolume - cgs
func DeleteVolume(in *model.Volume) (uint64, string) {
	var err error

	if in.UUID != "" {
		sql := "delete from volume where uuid = ?"
		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			goto ERROR
		}
		defer func() {
			_ = stmt.Close()
		}()
		result, err2 := stmt.Exec(in.UUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			goto ERROR
		}
		logger.Logger.Println(result.RowsAffected())

		return 0, ""
	}
ERROR:
	return hcc_errors.CelloSQLOperationFail, err.Error()
}
