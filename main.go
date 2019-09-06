package main

import (
	"GraphQL_Cello/cellocheckroot"
	"GraphQL_Cello/celloconfig"
	"GraphQL_Cello/cellographql"
	"GraphQL_Cello/cellologger"
	"GraphQL_Cello/cellomysql"
	"net/http"
)

func main() {
	if !cellocheckroot.CheckRoot() {
		return
	}

	if !cellologger.Prepare() {
		return
	}
	defer cellologger.FpLog.Close()

	err := cellomysql.Prepare()
	if err != nil {
		return
	}
	defer cellomysql.Db.Close()

	http.Handle("/graphql", cellographql.GraphqlHandler)

	cellologger.Logger.Println("Server is running on port " + celloconfig.HTTPPort)
	err = http.ListenAndServe(":"+celloconfig.HTTPPort, nil)
	if err != nil {
		cellologger.Logger.Println("Failed to prepare http server!")
	}
}
