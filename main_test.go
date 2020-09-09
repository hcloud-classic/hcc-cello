package main

import (
	"fmt"
	"hcc/cello/action/graphql"
	"hcc/cello/lib/config"
	"hcc/cello/lib/logger"
	"net/http"
	"strconv"
	"testing"
)

func TestPrepare(t *testing.T) {

	if !logger.Prepare() {
		fmt.Println("error occurred while preparing logger")
	}

	config.Parser()

	http.Handle("/graphql", graphql.GraphqlHandler)
	logger.Logger.Println("Opening server on port " + strconv.Itoa(int(config.HTTP.Port)) + "...")
	err := http.ListenAndServe(":"+strconv.Itoa(int(config.HTTP.Port)), nil)
	if err != nil {
		logger.Logger.Println(err)
		logger.Logger.Println("Failed to prepare http server!")
		return
	}

}
