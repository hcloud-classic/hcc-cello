package config

import (
	"hcc/cello/lib/logger"

	"github.com/Terry-Mao/goconf"
)

var conf = goconf.New()
var config = celloConfig{}
var err error

func parseMysql() {
	config.MysqlConfig = conf.Get("mysql")
	if config.MysqlConfig == nil {
		logger.Logger.Panicln("no mysql section")
	}

	Mysql = mysql{}
	Mysql.ID, err = config.MysqlConfig.String("id")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Password, err = config.MysqlConfig.String("password")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Address, err = config.MysqlConfig.String("address")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Port, err = config.MysqlConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	Mysql.Database, err = config.MysqlConfig.String("database")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseHTTP() {
	config.HTTPConfig = conf.Get("http")
	if config.HTTPConfig == nil {
		logger.Logger.Panicln("no http section")
	}

	HTTP = http{}
	HTTP.Port, err = config.HTTPConfig.Int("port")
	if err != nil {
		logger.Logger.Panicln(err)
	}
}

func parseVolumeHandle() {
	config.VolumeConfig = conf.Get("volumeHandle")
	if config.VolumeConfig == nil {
		logger.Logger.Panicln("no volumeHandle section")
	}

	VolumeConfig = volumeHandle{}
	VolumeConfig.VOLUMEPOOL, err = config.VolumeConfig.String("volume_pool")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	VolumeConfig.ORIGINVOL, err = config.VolumeConfig.String("origin_vol")
	logger.Logger.Println("asdasdasdasdasdas[", VolumeConfig.ORIGINVOL, "]")
	if err != nil {
		logger.Logger.Panicln(err)
	}

}

// Parser : Parse config file
func Parser() {
	if err = conf.Parse(configLocation); err != nil {
		logger.Logger.Panicln(err)
	}

	parseMysql()
	parseHTTP()
	parseVolumeHandle()
}
