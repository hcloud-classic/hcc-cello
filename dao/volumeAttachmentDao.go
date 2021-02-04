package dao

import (
	"hcc/cello/lib/logger"
	"hcc/cello/lib/mysql"
	"hcc/cello/model"
	"time"
)

func ReadVolumeAttachment(args map[string]interface{}) (interface{}, error) {
	var volumeAttachment model.VolumeAttachment
	var err error
	uuid := args["uuid"].(string)

	var volumeUUID string
	var serverUUID string
	var createdAt time.Time
	var updatedAt time.Time

	sql := "select * from volume_attachment where uuid = ? order by created_at desc"
	err = mysql.Db.QueryRow(sql, uuid).Scan(
		&uuid,
		&volumeUUID,
		&serverUUID,
		&createdAt,
		&updatedAt)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}

	volumeAttachment.UUID = uuid
	volumeAttachment.VolumeUUID = volumeUUID
	volumeAttachment.ServerUUID = serverUUID
	volumeAttachment.CreatedAt = createdAt
	volumeAttachment.UpdatedAt = updatedAt

	return volumeAttachment, nil
}

func ReadVolumeAttachmentList(args map[string]interface{}) (interface{}, error) {
	var err error
	var volumeAttachments []model.VolumeAttachment
	var uuid string
	var createdAt time.Time
	var updatedAt time.Time

	volumeUUID, volumeUUIDOk := args["volume_uuid"].(string)
	serverUUID, serverUUIDOk := args["server_uuid"].(string)

	sql := "select * from volume_attachment where 1 = 1 and "
	if volumeUUIDOk {
		sql += " volume_uuid = '" + volumeUUID + "' order by created_at desc"
	} else if serverUUIDOk {
		sql += " server_uuid = '" + serverUUID + "' order by created_at desc"
	} else if volumeUUIDOk && serverUUIDOk {
		sql += " volume_uuid = '" + volumeUUID + "and server_uuid = '" + serverUUID + "' order by created_at desc"
	}

	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &volumeUUID, &serverUUID, &createdAt, &updatedAt)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		volumeAttachment := model.VolumeAttachment{UUID: uuid, VolumeUUID: volumeUUID, ServerUUID: serverUUID, CreatedAt: createdAt, UpdatedAt: updatedAt}
		volumeAttachments = append(volumeAttachments, volumeAttachment)
	}
	return volumeAttachments, nil
}

func ReadVolumeAttachmentAll(args map[string]interface{}) (interface{}, error) {

	var err error
	var volumeAttachments []model.VolumeAttachment
	var uuid string
	var volumeUUID string
	var serverUUID string
	var createdAt time.Time
	var updatedAt time.Time

	sql := "select * from volume_attachment order by created_at desc"

	stmt, err := mysql.Db.Query(sql)
	if err != nil {
		logger.Logger.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()

	for stmt.Next() {
		err := stmt.Scan(&uuid, &volumeUUID, &serverUUID, &createdAt, &updatedAt)

		if err != nil {
			logger.Logger.Println(err)
			return nil, err
		}
		volumeAttachment := model.VolumeAttachment{UUID: uuid, VolumeUUID: volumeUUID, ServerUUID: serverUUID, CreatedAt: createdAt, UpdatedAt: updatedAt}
		volumeAttachments = append(volumeAttachments, volumeAttachment)
	}

	return volumeAttachments, nil
}

func CreateVolumeAttachment(args map[string]interface{}) (interface{}, error) {
	out, err := gouuid.NewV4()
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	uuid := out.String()

	volumeAttachment := model.VolumeAttachment{
		UUID:       uuid,
		VolumeUUID: args["volume_uuid"].(string),
		ServerUUID: args["server_uuid"].(string),
	}

	sql := "insert into volume_attachment(uuid, volume_uuid, server_uuid, created_at, updated_at) values (?, ?, ?, now(), now())"
	stmt, err := mysql.Db.Prepare(sql)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	defer func() {
		_ = stmt.Close()
	}()
	result, err := stmt.Exec(volumeAttachment.UUID, volumeAttachment.VolumeUUID, volumeAttachment.ServerUUID)
	if err != nil {
		logger.Logger.Println(err)
		return nil, err
	}
	logger.Logger.Println(result.LastInsertId())

	return volumeAttachment, nil
}

func UpdateVolumeAttachment(args map[string]interface{}) (interface{}, error) {

	var err error
	var volumeAttachment model.VolumeAttachment
	volumeUUID, volumeUUIDOk := args["volume_uuid"].(string)
	serverUUID, serverUUIDOk := args["server_uuid"].(string)

	if volumeUUIDOk && serverUUIDOk {

		sql := "update volume_attachment set server_uuid = " + serverUUID + " where volume_uuid = ?"

		stmt, err := mysql.Db.Prepare(sql)
		if err != nil {
			logger.Logger.Println(err.Error())
			return nil, err
		}
		defer func() {
			_ = stmt.Close()
		}()

		result, err2 := stmt.Exec(volumeUUID)
		if err2 != nil {
			logger.Logger.Println(err2)
			return nil, err
		}
		logger.Logger.Println(result.RowsAffected())

		return volumeAttachment, nil
	}

	return nil, err
}
func DeleteVolumeAttachment(args map[string]interface{}) (interface{}, error) {
	var err error

	requestedUUID, ok := args["uuid"].(string)
	if ok {
		sql := "delete from volume_attachment where uuid = ?"
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
