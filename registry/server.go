package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

/*this file would be defining the service definitions*/

const ServerPort = ":3000"
const ServiceURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

//used for registration
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

//used for de-registration
func (r *registry) remove(url string) error {

	for index, registration := range r.registrations {

		if registration.ServiceURL == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:index], r.registrations[index+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}

	return fmt.Errorf("Service at URL %v not found", url)
}

//create a registry instance
var reg = registry{registrations: make([]Registration, 0), mutex: new(sync.Mutex)}

//create the registry service itself
type RegistryService struct{}

//this service should be able to accept new registration
func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Received  request")
	switch r.Method {
	case http.MethodPost:
		dec := json.NewDecoder(r.Body)
		var reqReg Registration
		err := dec.Decode(&reqReg)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//add to the registry
		log.Printf("Adding service : %v with URL : %v\n", reqReg.ServiceName, reqReg.ServiceURL)

		err = reg.add(reqReg)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodDelete:

		payload, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := string(payload)

		log.Printf("Removing service with URL : %v\n", url)

		err = reg.remove(url)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
}
