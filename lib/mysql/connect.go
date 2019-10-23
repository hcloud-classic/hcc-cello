package mysql

import (
	"database/sql"
	"hcc/cello/lib/config"
	"hcc/cello/lib/logger"

	_ "github.com/go-sql-driver/mysql" // Needed for connect mysql
)

// Db : Pointer of mysql connection
var Db *sql.DB

// Prepare : Connect to mysql and prepare pointer of mysql connection
func Prepare() error {
	var err error
	Db, err = sql.Open("mysql", config.MysqlID+":"+config.MysqlPassword+"@tcp("+
		config.MysqlAddress+":"+config.MysqlPort+")/"+config.MysqlDatabase+"?parseTime=true")
	if err != nil {
		logger.Logger.Println(err)
		return err
	}

	logger.Logger.Println("db is connected")

	err = Db.Ping()
	if err != nil {
		logger.Logger.Println(err.Error())
		return err
	}

	return nil
}
