package config

type volumeHandle struct {
	VOLUMEPOOL string   `goconf:"volume:volume_pool"` // For storage vol zpool
	ORIGINVOL  []string `goconf:"volume:origin_vol"`  // Original volume
	SupportOS  []string `goconf:"volume:support_os"`  // Original volume
}

// VolumeConfig : cello volume config
var VolumeConfig volumeHandle
