package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterService(r Registration) error {

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(r)

	if err != nil {
		return err
	}

	res, err := http.Post(ServiceURL, "application/json", buf)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register service, got response code %v from registry service", res.StatusCode)
	}

	return nil

}

func ShutdownService(serviceURL string) error {

	req, err := http.NewRequest(http.MethodDelete, ServiceURL, bytes.NewBuffer([]byte(serviceURL)))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "text/plain")
	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK || err != nil {
		return fmt.Errorf("de-registration failed, registry service returned with code %v", res.StatusCode)
	}

	return nil

}
