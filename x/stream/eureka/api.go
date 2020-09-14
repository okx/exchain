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
func register(instance *Instance, serverURL string, appName string) error {
	b, err := json.Marshal(&InstanceInfo{Instance: instance})
	if err != nil {
		return err
	}

	urlAction := serverURL + "/eureka/apps" + "/" + appName
	req, err := http.NewRequest(http.MethodPost, urlAction, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("status code is %d, not in range [200, 300)", resp.StatusCode)
	}
	return nil
}

// DELETE /eureka/v2/apps/appID/instanceID
func unRegister(serverURL string, appName string, instanceID string) error {
	urlAction := serverURL + "/eureka/apps" + "/" + appName + "/" + instanceID
	req, err := http.NewRequest(http.MethodDelete, urlAction, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d, require 200", resp.StatusCode)
	}
	return nil
}

// nolint
func getAllInstance(serverURL string) (*Applications, error) {
	urlAction := serverURL + "/eureka/apps"
	req, err := http.NewRequest(http.MethodGet, urlAction, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accpet", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is %d, require 200", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := &Applications{}
	err = xml.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetOneInstance(serverURL string, appName string) (*Application, error) {
	urlAction := serverURL + "/eureka/apps" + "/" + appName
	req, err := http.NewRequest(http.MethodGet, urlAction, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accpet", "application/xml")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is %d, require 200", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &Application{}
	err = xml.Unmarshal(b, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// PUT /eureka/apps/appID/instanceID
func heartbeat(serverURL string, appName string, instanceID string) error {
	params := url.Values{
		"status": {"UP"},
	}

	urlAction := serverURL + "/eureka/apps" + "/" + appName + "/" + instanceID + "?" + params.Encode()
	req, err := http.NewRequest(http.MethodPut, urlAction, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status code is %d, require 200", resp.StatusCode)
	}
	return nil
}
