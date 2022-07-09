package log

import (
	"io/ioutil"
	stlog "log"
	"net/http"
	"os"
)

/*
	This file implements this core logic for log service's server side i.e business logic
*/

//custom logger that will handle logging
var log *stlog.Logger

//type to handle writing to a filesystem
type fileLog string

//implement writer interface : why? refer :https://pkg.go.dev/log#New
func (fl fileLog) Write(data []byte) (int, error) {

	f, err := os.OpenFile(string(fl), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return 0, err
	}

	defer f.Close()

	return f.Write(data)
}

//instantiate this logger i.e point it to some file
func Run(destination string) {
	log = stlog.New(fileLog(destination), "rdlogger", stlog.LstdFlags)
}

//register the http endpoints for this service
func RegisterHandlers() {

	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		msg, err := ioutil.ReadAll(r.Body)
		if err != nil || len(msg) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		write(string(msg))

	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
