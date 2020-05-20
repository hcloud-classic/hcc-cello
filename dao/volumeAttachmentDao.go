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
