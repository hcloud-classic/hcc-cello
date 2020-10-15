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

func parseGrpc() {
	config.GrpcConfig = conf.Get("grpc")
	if config.GrpcConfig == nil {
		logger.Logger.Panicln("no grpc section")
	}

	Grpc.Port, err = config.GrpcConfig.Int("port")
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

	//Will be Deprecated
	VolumeConfig.VOLUMEPOOL, err = config.VolumeConfig.String("volume_pool")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	VolumeConfig.SupportOS, err = config.VolumeConfig.Strings("support_os", " ")
	logger.Logger.Println("Support Os List[", VolumeConfig.SupportOS, "] : ")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	VolumeConfig.ORIGINVOL, err = config.VolumeConfig.Strings("origin_vol", " ")
	logger.Logger.Println("Support Volume List[", VolumeConfig.ORIGINVOL, "]")
	if err != nil {
		logger.Logger.Panicln(err)
	}
	VolumeConfig.IscsiDiscoveryAddress, err = config.VolumeConfig.Strings("iscsi_discovery_address", " ")
	logger.Logger.Println("IscsiDiscoveryAddress[", VolumeConfig.IscsiDiscoveryAddress, "] : ")
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
	parseGrpc()
	parseHTTP()
	parseVolumeHandle()
}
