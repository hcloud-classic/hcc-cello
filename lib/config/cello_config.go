package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/cello/cello.conf"

type celloConfig struct {
	MysqlConfig  *goconf.Section
	GrpcConfig   *goconf.Section
	HTTPConfig   *goconf.Section
	CelloConfig  *goconf.Section
	VolumeConfig *goconf.Section
	SupportOS    *goconf.Section
}
