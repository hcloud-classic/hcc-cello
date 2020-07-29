package model

import "time"

type Volume struct {
	UUID       string    `json:"uuid"`
	Size       int       `json:"size"`
	Filesystem string    `json:"filesystem"` //os
	ServerUUID string    `json:"server_uuid"`
	UseType    string    `json:"use_type"` //
	UserUUID   string    `json:"user_uuid"`
	CreatedAt  time.Time `json:"created_at"`
	NetworkIP  string    `json:"network_ip"`
	GatewayIP  string    `json:"gateway_ip"`
}

type Volumes struct {
	Volumes []Volume `json:"volume"`
}

type VolumeNum struct {
	Number int `json:"number"`
}
