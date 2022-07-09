package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dravikant/go-distributed-app/registry"
)

/*
	we wont be using service/service.go to start the service as it is specifically
	designed to handle client service
	and registration service wont be taking advantage of that

*/

func main() {
	http.Handle("/services", &registry.RegistryService{})

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	var srv http.Server

	srv.Addr = registry.ServerPort

	go func() {
		log.Println(srv.ListenAndServe())
		//if the above stmt returns, that means we ran into error, so cancel the context
		cancel()
	}()

	//another go routine to provide user an option to stop the service i.e. server
	go func() {
		log.Println("Registry service started. Press any key to stop")
		var s string
		fmt.Scan(&s)
		//if we receive input from user, gracefully shutdown the server i.e service
		srv.Shutdown(ctx)
		//also cancel the context
		cancel()
	}()

	<-ctx.Done()

	fmt.Println("Shutting down Registry service")

}
