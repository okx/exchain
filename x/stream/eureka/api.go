package eureka

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// POST /eureka/v2/apps/appID
func register(instance *Instance, serverUrl string, appName string) error {
	b, err := json.Marshal(&InstanceInfo{Instance: instance})
	if err != nil {
		return err
	}

	urlAction := serverUrl + "/eureka/apps" + "/" + appName
	req, err := http.NewRequest(http.MethodPost, urlAction, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("status code is %d, not in range [200, 300)", resp.StatusCode)
	}
	return nil
}

// DELETE /eureka/v2/apps/appID/instanceID
func unRegister(serverUrl string, appName string, instanceId string) error {
	urlAction := serverUrl + "/eureka/apps" + "/" + appName + "/" + instanceId
	req, err := http.NewRequest(http.MethodDelete, urlAction, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d, require 200", resp.StatusCode)
	}
	return nil
}

// GET /eureka/v2/apps
func getAllInstance(serverUrl string) (*Applications, error) {
	urlAction := serverUrl + "/eureka/apps"
	req, err := http.NewRequest(http.MethodGet, urlAction, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accpet", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is %d, require 200", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &Applications{}
	err = xml.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetOneInstance(serverUrl string, appName string) (*Application, error) {
	urlAction := serverUrl + "/eureka/apps" + "/" + appName
	req, err := http.NewRequest(http.MethodGet, urlAction, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accpet", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is %d, require 200", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := &Application{}
	err = xml.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PUT /eureka/apps/appID/instanceID
func heartbeat(serverUrl string, appName string, instanceId string) error {
	params := url.Values{
		"status": {"UP"},
	}

	urlAction := serverUrl + "/eureka/apps" + "/" + appName + "/" + instanceId + "?" + params.Encode()
	req, err := http.NewRequest(http.MethodPut, urlAction, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d, require 200", resp.StatusCode)
	}
	return nil
}
