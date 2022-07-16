package registry

import (
	"bytes"
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
	mutex         *sync.RWMutex
}

//used for registration
func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	r.mutex.Unlock()
	//registration is done Now update other registered services
	err := r.sendRequiredServices(reg)
	if err != nil {
		fmt.Printf("unable to notify registration %v", err)
	}
	//problem with only having sendRequiredServices is that a service gets notified
	//about the required service only iff it is already running at the startup
	//to notify it whenever a required service comes up following method is used
	r.notify(patch{
		Added: []patchEntry{
			patchEntry{Name: reg.ServiceName, URL: reg.ServiceURL},
		},
	})
	return nil
}

func (r *registry) sendRequiredServices(reg Registration) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch

	fmt.Printf("sending required service for %v\n", reg.ServiceName)
	//get all registrations from registry
	for _, serviceReg := range r.registrations {
		//get required services from current registrations
		for _, reqService := range reg.RequiredServices {

			if serviceReg.ServiceName == reqService {
				p.Added = append(p.Added, patchEntry{
					Name: serviceReg.ServiceName,
					URL:  serviceReg.ServiceURL,
				})
				fmt.Printf("found required service %v\n", serviceReg.ServiceName)
			}
		}
	}

	return r.sendPatch(p, reg.ServiceUpdateURL)
}

//only notify services which are interested
func (r registry) notify(fullpatch patch) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	for _, reg := range r.registrations {
		go func(reg Registration) {

			for _, reqServices := range reg.RequiredServices {
				p := patch{Added: []patchEntry{}, Removed: []patchEntry{}}
				sendUpdate := false
				for _, added := range fullpatch.Added {
					if added.Name == reqServices {
						p.Added = append(p.Added, added)
						sendUpdate = true
					}
				}
				for _, removed := range fullpatch.Removed {
					if removed.Name == reqServices {
						p.Removed = append(p.Removed, removed)
						sendUpdate = true
					}
				}

				if sendUpdate {
					err := r.sendPatch(p, reg.ServiceUpdateURL)
					if err != nil {
						fmt.Println(err)
						return
					}
				}
			}
		}(reg)
	}
}

func (r registry) sendPatch(p patch, url string) error {
	d, err := json.Marshal(p)
	if err != nil {
		fmt.Printf("error in marshal %v", err)
		return err
	}
	fmt.Printf("d: %v\n", d)
	fmt.Printf("p :%v\n", p)
	_, err = http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		fmt.Printf("error in post %v", err)
		//how abt retry

		req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(d))
		if err != nil {
			return err
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			fmt.Printf("error in retry post %v", err)
		}
		// req.Close = true

		fmt.Println("trying another time ")
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(d))
		var client = &http.Client{
			Transport: &http.Transport{},
		}

		_, err = client.Do(req)
		fmt.Printf("error in 2retry post %v", err)
		return err
	}
	return nil
}

//used for de-registration
func (r *registry) remove(url string) error {

	for index, registration := range r.registrations {

		if registration.ServiceURL == url {
			r.notify(patch{
				Removed: []patchEntry{
					patchEntry{
						Name: registration.ServiceName,
						URL:  registration.ServiceURL,
					}},
			})
			r.mutex.Lock()
			r.registrations = append(r.registrations[:index], r.registrations[index+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}

	return fmt.Errorf("service at URL %v not found", url)
}

//create a registry instance
var reg = registry{registrations: make([]Registration, 0), mutex: new(sync.RWMutex)}

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
