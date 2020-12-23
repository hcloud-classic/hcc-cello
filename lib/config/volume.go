package config

type volumeHandle struct {
	VOLUMEPOOL            string   `goconf:"volume:volume_pool"`             // For storage vol zpool
	ROOTUUID              string   `goconf:"root_uuid"`                      // For storage vol zpool
	ORIGINVOL             []string `goconf:"volume:origin_vol"`              // Original volume
	SupportOS             []string `goconf:"volume:support_os"`              // Original volume
	IscsiDiscoveryAddress []string `goconf:"volume:iscsi_discovery_address"` // For Iscsi Deploy
}

// VolumeConfig : cello volume config
var VolumeConfig volumeHandle
