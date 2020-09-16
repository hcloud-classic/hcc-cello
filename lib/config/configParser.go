package config

import (
	"hcc/cello/lib/logger"

	"github.com/Terry-Mao/goconf"
)

// var iscsiconf = fomatter.New()
var conf = goconf.New()
var config = celloConfig{}
var err error

// func loadIscsiDB(){
// 	celloParams := make(map[string]interface{})
// 	celloParams["row"] = 0
// 	celloParams["page"] = 0
// 	var volumes []model.Volume
// 	volumes := ReadVolumeAll(celloParams)
// 	body, _ := json.Marshal(volumes)
// 	var obj map[string]interface{}
// 	json.Unmarshal([]byte(body), &obj)

// 	// Make a custom formatter with indent set
// 	f := colorjson.NewFormatter()
// 	f.Indent = 4

// 	// Marshall the Colorized JSON
// 	s, _ := f.Marshal(obj)
// 	// fmt.Println(string(s))
// 	logger.Logger.Println("doHcc Action [", string(s), "]")

// }
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

	//Will be Deprecated
	VolumeConfig.VOLUMEPOOL, err = config.VolumeConfig.String("volume_pool")
	if err != nil {
		logger.Logger.Panicln(err)
	}

	VolumeConfig.SupportOS, err = config.VolumeConfig.Strings("support_os", " ")
	logger.Logger.Println("Support Os List[", VolumeConfig.SupportOS, "] : ")
	// for _, args := range VolumeConfig.SupportOS {
	// 	logger.Logger.Println("Support Os List[", VolumeConfig.SupportOS, args, "]")

	// }
	if err != nil {
		logger.Logger.Panicln(err)
	}
	VolumeConfig.ORIGINVOL, err = config.VolumeConfig.Strings("origin_vol", " ")
	logger.Logger.Println("Support Volume List[", VolumeConfig.ORIGINVOL, "]")
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
