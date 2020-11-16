package main

import (
	"hcc/cello/action/grpc/server"
	celloEnd "hcc/cello/end"
	celloInit "hcc/cello/init"

	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	err := celloInit.MainInit()
	if err != nil {
		panic(err)
	}
}

func main() {
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		celloEnd.MainEnd()
		fmt.Println("Exiting violin module...")
		os.Exit(0)
	}()

	server.Init()

}
