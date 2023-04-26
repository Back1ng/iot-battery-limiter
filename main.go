package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OauthToken struct {
	token string
}

var ApiUrl = "https://api.iot.yandex.net"

type ResponseGetDevices struct {
	Devices []struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		Type         string `json:"type"`
		Capabilities []struct {
			Reportable  bool   `json:"reportable"`
			Retrievable bool   `json:"retrievable"`
			Type        string `json:"type"`
			State       struct {
				Instance string `json:"instance"`
				Value    bool   `json:"value"`
			} `json:"state"`
		} `json:"capabilities"`
	} `json:"devices"`
}

type RequestEnableDevice struct {
	Devices []struct {
		Id      string `json="id"`
		Actions []struct {
			Type  string `json="type"`
			State struct {
				Instance string      `json="instance"`
				Value    interface{} `json="value"`
			} `json:"state"`
		} `json:"actions"`
	} `json:"devices"`
}

type ResponseChangeStatusDevice struct {
	Status    string `json:"status"`
	RequestID string `json:"request_id"`
	Devices   []struct {
		ID           string `json:"id"`
		Capabilities []struct {
			Type  string `json:"type"`
			State struct {
				Instance     string `json:"instance"`
				ActionResult struct {
					Status string `json:"status"`
				} `json:"action_result"`
			} `json:"state"`
		} `json:"capabilities"`
	} `json:"devices"`
}

func (api *OauthToken) Info() ResponseGetDevices {

	req, err := http.NewRequest("GET", ApiUrl+"/v1.0/user/info", nil)
	req.Header.Add("Authorization", api.token)

	if err != nil {
		// handle error
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle error
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		// handle error
	}

	var result ResponseGetDevices
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	res.Body.Close()

	return result
}

func (api *OauthToken) Enable(id string) {
	body := []byte(fmt.Sprintf(`{
		"devices": [{
			"id": "%s",
			"actions": [{
				"type": "devices.capabilities.on_off",
				"state": {
					"instance": "on",
					"value": true
				}
			}]
		}]
	}`, id))

	req, err := http.NewRequest("POST", ApiUrl+"/v1.0/devices/actions", bytes.NewBuffer(body))

	if err != nil {
		// handle error
	}

	req.Header.Add("Authorization", api.token)
	req.Header.Add("X-Request-Id", "ac851019-b7dd-45dc-aace-0d9a8e47fc66")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		// handle error
	}

	resBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	fmt.Println(string(resBody))

	var result ResponseChangeStatusDevice
	if err := json.Unmarshal(resBody, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
}

func (api *OauthToken) Disable(id string) {
	body := []byte(fmt.Sprintf(`{
		"devices": [{
			"id": "%s",
			"actions": [{
				"type": "devices.capabilities.on_off",
				"state": {
					"instance": "on",
					"value": false
				}
			}]
		}]
	}`, id))

	req, err := http.NewRequest("POST", ApiUrl+"/v1.0/devices/actions", bytes.NewBuffer(body))

	if err != nil {
		// handle error
	}

	req.Header.Add("Authorization", api.token)
	req.Header.Add("X-Request-Id", "ac851019-b7dd-45dc-aace-0d9a8e47fc66")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		// handle error
	}

	resBody, err := io.ReadAll(res.Body)

	defer res.Body.Close()

	var result ResponseChangeStatusDevice
	if err := json.Unmarshal(resBody, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
}

func main() {
	str := ""
	token := OauthToken{"Bearer " + str}

	devices := token.Info()

	id := devices.Devices[0].Id

	// token.Enable(id)
	// token.Disable(id)
}
