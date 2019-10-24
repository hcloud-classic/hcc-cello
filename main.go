package main

import (
	"hcc/cello/action/graphql"
	"hcc/cello/lib/config"
	"hcc/cello/lib/logger"
	"hcc/cello/lib/mysql"
	"hcc/cello/lib/syscheck"
	"net/http"
	"strconv"
)

func main() {
	if !syscheck.CheckRoot() {
		return
	}

	if !logger.Prepare() {
		return
	}
	defer func() {
		_ = logger.FpLog.Close()
	}()

	config.Parser()

	err := mysql.Prepare()
	if err != nil {
		return
	}
	defer func() {
		_ = mysql.Db.Close()
	}()
	// err = logger.CreateDirIfNotExist("/root/boottp/HCC/" + "XXXXXXXXX")
	// logger.Logger.Println(err)
	// if err != nil {
	// 	return
	// }
	http.Handle("/graphql", graphql.Handler)

	logger.Logger.Println("Server is running on port " + strconv.Itoa(int(config.HTTP.Port)))
	err = http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println("Failed to prepare http server!")
	}
}
