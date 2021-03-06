package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/dravikant/go-distributed-app/log"
	"github.com/dravikant/go-distributed-app/registry"
	"github.com/dravikant/go-distributed-app/service"
)

/*
	Following the go's convention to house the binaries/commands in the cmd package
	This file will provide log service executable command
*/

func main() {
	// this is going to set the destination i.e log file
	log.Run("./app.log")

	//instantiate host and port
	//TODO: read these from config file
	host, port := "localhost", "4000"
	// serviceAddress := fmt.Sprintf("http://%v:%v", host, port)
	serviceAddress := "http://" + host + ":" + port

	var r registry.Registration

	r.ServiceName = registry.LogService
	r.ServiceURL = serviceAddress
	r.RequiredServices = make([]registry.ServiceName, 0)
	r.ServiceUpdateURL = r.ServiceURL + "/services"

	ctx, err := service.Start(context.Background(), r, host, port, log.RegisterHandlers)

	if err != nil {
		stlog.Fatal(err)
	}

	//wait for the context.Done signal
	//this can be executed in case of cancel call from either of two go routines
	// i.e. in case server fails to start or user provides input to shutdown the service
	<-ctx.Done()

	//print the exit message
	fmt.Println("Shutting down Log service ")
}
