package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/dravikant/go-distributed-app/grades"
	"github.com/dravikant/go-distributed-app/registry"
	"github.com/dravikant/go-distributed-app/service"
)

func main() {

	host, port := "localhost", "6000"
	serviceAddress := "http://" + host + ":" + port

	var r registry.Registration

	r.ServiceName = registry.GradingService
	r.ServiceURL = serviceAddress

	ctx, err := service.Start(context.Background(), r, host, port, grades.RegisterHandlersFunc)

	if err != nil {
		stlog.Fatal("err")
	}

	<-ctx.Done()

	fmt.Println("Shutting Down the Grading service")

}
