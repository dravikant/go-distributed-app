package registry

import (
	"encoding/json"
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

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	return nil
}

//create a registry instance
var reg = registry{registrations: make([]Registration, 0), mutex: new(sync.Mutex)}

//create the registry service itself
type RegistryService struct{}

//this service should be able to accept new registration
func (s RegistryService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Received registration request")
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

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
}
