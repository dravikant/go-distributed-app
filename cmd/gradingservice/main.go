package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/dravikant/go-distributed-app/grades"
	"github.com/dravikant/go-distributed-app/log"
	"github.com/dravikant/go-distributed-app/registry"
	"github.com/dravikant/go-distributed-app/service"
)

func main() {

	host, port := "localhost", "6000"
	serviceAddress := "http://" + host + ":" + port

	var r registry.Registration

	r.ServiceName = registry.GradingService
	r.ServiceURL = serviceAddress

	//add the required services
	r.RequiredServices = []registry.ServiceName{registry.LogService}
	r.ServiceUpdateURL = r.ServiceURL + "/services"

	ctx, err := service.Start(context.Background(), r, host, port, grades.RegisterHandlersFunc)

	if err != nil {
		stlog.Fatal("err")
	}

	//check if we got the required services
	if logProvider, err := registry.GetProvider(registry.LogService); err == nil {
		fmt.Printf("Logging service found at %v", logProvider)
		log.SetClientLogger(logProvider, r.ServiceName)
	}

	<-ctx.Done()

	fmt.Println("Shutting Down the Grading service")

}
