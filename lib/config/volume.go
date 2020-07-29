package config

type volumeHandle struct {
	VOLUMEPOOL string `goconf:"volume:volume_pool"` // For storage vol zpool
	ORIGINVOL  string `goconf:"volume:origin_vol"`  // Original volume
}

// Mysql : mysql config structure
var VolumeConfig volumeHandle
