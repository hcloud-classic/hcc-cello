package main

import (
	"hcc/cello/action/graphql"
	celloEnd "hcc/cello/end"
	celloInit "hcc/cello/init"
	"hcc/cello/lib/config"
	"hcc/cello/lib/logger"
	"net/http"
	"strconv"
)

func init() {
	err := celloInit.MainInit()
	if err != nil {
		panic(err)
	}
}
func main() {
	defer func() {
		celloEnd.MainEnd()
	}()

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println(err)
		logger.Logger.Println("Failed to prepare http server!")
		return
	}
}
