package model

import "time"

// DefaultPXEdir : Default PXE directory
var DefaultPXEdir = "/root/boottp/HCC"

// OSDiskSize : Disk size for OS use
var OSDiskSize = 100

// Volume - cgs
type Volume struct {
	UUID       string    `json:"uuid"`
	Size       int       `json:"size"`
	Filesystem string    `json:"filesystem"` //os
	ServerUUID string    `json:"server_uuid"`
	UseType    string    `json:"use_type"` //
	UserUUID   string    `json:"user_uuid"`
	LunNum     int       `json:"lun_num"`
	Pool       string    `json:"pool"`
	CreatedAt  time.Time `json:"created_at"`
	NetworkIP  string    `json:"network_ip"`
	GatewayIP  string    `json:"gateway_ip"`
}

// Volumes - cgs
type Volumes struct {
	Volumes []Volume `json:"volume"`
}

// VolumeNum - cgs
type VolumeNum struct {
	Number int `json:"number"`
}
