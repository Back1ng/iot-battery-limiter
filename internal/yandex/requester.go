package yandex

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/back1ng/iot-battery-limiter/internal/env"
)

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

func sendSafeRequest(req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)

	for err != nil {
		res, err = http.DefaultClient.Do(req)

		fmt.Println("Cannot send request to Yandex! Please, check your internet connection.")
		fmt.Println("Trying again after 10 second.")
		fmt.Println("")

		time.Sleep(time.Second * 10)
	}

	return res
}

func (api *OauthToken) addHeaders(header http.Header) {
	header.Add("Authorization", api.Token)
	header.Add("X-Request-Id", "ac851019-b7dd-45dc-aace-0d9a8e47fc66")
	header.Add("Content-Type", "application/json")
}

func (api *OauthToken) Info() ResponseGetDevices {
	req, err := http.NewRequest("GET", ApiUrl+"/v1.0/user/info", nil)

	req.Header.Add("Authorization", api.Token)

	if err != nil {
		// handle error
	}

	res := sendSafeRequest(req)

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		// handle error
	}

	var result ResponseGetDevices
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return result
}

func (api *OauthToken) GetBody(id, value string) []byte {
	return []byte(fmt.Sprintf(`{
		"devices": [{
			"id": "%s",
			"actions": [{
				"type": "devices.capabilities.on_off",
				"state": {
					"instance": "on",
					"value": %s
				}
			}]
		}]
	}`, id, value))
}

func (api *OauthToken) Enable(id string) {
	body := api.GetBody(id, "true")

	req, err := http.NewRequest("POST", ApiUrl+"/v1.0/devices/actions", bytes.NewBuffer(body))

	if err != nil {
		// throw error
	}

	api.addHeaders(req.Header)

	res := sendSafeRequest(req)

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
	}

	defer res.Body.Close()

	fmt.Println(string(resBody))

	var result ResponseChangeStatusDevice
	if err := json.Unmarshal(resBody, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
}

func (api *OauthToken) Disable(id string) {
	body := api.GetBody(id, "false")

	req, err := http.NewRequest("POST", ApiUrl+"/v1.0/devices/actions", bytes.NewBuffer(body))

	if err != nil {
		// handle error
	}

	api.addHeaders(req.Header)

	res := sendSafeRequest(req)

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		// handle error
	}

	defer res.Body.Close()

	var result ResponseChangeStatusDevice
	if err := json.Unmarshal(resBody, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}
}

func (api *OauthToken) PrintDevices() {
	devices := api.Info()

	fmt.Println("Please, choose your device number by number before \")\"")
	for i, d := range devices.Devices {
		fmt.Println(fmt.Sprintf("%d) Name: %s, Id: %s", i, d.Name, d.Id))
	}

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	id, err := strconv.Atoi(input.Text())

	if err != nil {
		// handle error
	}

	os.Setenv("DEVICE_ID", devices.Devices[id].Id)

	file, err := os.Open("./.env")

	defer file.Close()

	env.Write(file, "DEVICE_ID", devices.Devices[id].Id)

	fmt.Println("You select: " + input.Text())
}
