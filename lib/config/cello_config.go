package config

import "github.com/Terry-Mao/goconf"

var configLocation = "/etc/hcc/cello/cello.conf"

type celloConfig struct {
	MysqlConfig  *goconf.Section
	HTTPConfig   *goconf.Section
	CelloConfig  *goconf.Section
	VolumeConfig *goconf.Section
}

/*-----------------------------------
         Config File Example
/*-----------------------------------
[mysql]
id root
password qwe1212!Q
address 192.168.110.240
port 3306
database cello

[http]
port 7200

-----------------------------------*/
