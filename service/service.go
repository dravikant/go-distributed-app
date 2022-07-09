package service

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dravikant/go-distributed-app/registry"
)

/*
	This file has service startup logic
*/

//public function to kick-start the web service
//it is kept generic to use it for multiple services
func Start(ctx context.Context, reg registry.Registration, host, port string, registerHandlersFunc func()) (context.Context, error) {

	//call the register handler func
	registerHandlersFunc()

	//create new context
	ctx = startService(ctx, string(reg.ServiceName), host, port)

	//register the service
	err := registry.RegisterService(reg)

	if err != nil {
		return ctx, nil
	}

	return ctx, nil

}

func startService(ctx context.Context, serviceName, host, port string) context.Context {

	ctx, cancel := context.WithCancel(ctx)

	//create a server instance
	var srv http.Server
	srv.Addr = ":" + port

	//goroutine to start the server
	go func() {
		log.Println(srv.ListenAndServe())
		//if the above stmt returns, that means we ran into error, so cancel the context
		cancel()
	}()

	//another go routine to provide user an option to stop the service i.e. server
	go func() {
		log.Printf("%v started. Press any key to stop\n", serviceName)
		var s string
		fmt.Scan(&s)
		//if we receive input from user, gracefully shutdown the server i.e service and remove from registry
		err := registry.ShutdownService(fmt.Sprintf("http://%v:%v", host, port))
		if err != nil {
			log.Println(err)
		}
		srv.Shutdown(ctx)
		//also cancel the context
		cancel()
	}()

	return ctx
}
