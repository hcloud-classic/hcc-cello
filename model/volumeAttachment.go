package model

import "time"

type VolumeAttachment struct {
	UUID       string    `json:"uuid"`
	VolumeUUID string    `json:"volume_uuid"`
	ServerUUID string    `json:"server_uuid"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type VolumeAttachments struct {
	VolumeAttachments []VolumeAttachment `json:"volumeAttachment"`
}
