package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/cello/cello.conf"

type celloConfig struct {
	MysqlConfig *goconf.Section
	HTTPConfig  *goconf.Section
	CelloConfig *goconf.Section
}
