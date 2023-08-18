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
)

type Token string

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

func sendSafeRequest(req *http.Request) *http.Response {
	res, err := http.DefaultClient.Do(req)

	for err != nil {
		fmt.Println("Cannot send request to Yandex! Please, check your internet connection.")
		fmt.Println("Trying again after 10 second.")
		fmt.Println("")

		time.Sleep(time.Second * 10)

		res, err = http.DefaultClient.Do(req)
	}

	return res
}

func (api Token) addHeaders(header http.Header) {
	header.Add("Authorization", "Bearer "+string(api))
	header.Add("X-Request-Id", "ac851019-b7dd-45dc-aace-0d9a8e47fc66")
	header.Add("Content-Type", "application/json")
}

func (api Token) Info() ResponseGetDevices {
	var result ResponseGetDevices

	req, err := http.NewRequest("GET", ApiUrl+"/v1.0/user/info", nil)
	if err != nil {
		fmt.Printf("OauthToken.Info - http.NewRequest: %v", err)
		return result
	}

	api.addHeaders(req.Header)

	res := sendSafeRequest(req)
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("OauthToken.Info - io.ReadAll: %v", err)
		return result
	}

	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Can not unmarshal JSON")
	}

	return result
}

func (api *Token) GetBody(id, value string) []byte {
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

func (api *Token) Enable(id string) {
	body := api.GetBody(id, "true")

	req, err := http.NewRequest("POST", ApiUrl+"/v1.0/devices/actions", bytes.NewBuffer(body))

	if err != nil {
		fmt.Printf("OauthToken.Enable - http.NewRequest: %v", err)
		return
	}

	api.addHeaders(req.Header)

	res := sendSafeRequest(req)

	resBody, err := io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("OauthToken.Enable - io.ReadAll: %v", err)
		return
	}

	defer res.Body.Close()

	fmt.Println(string(resBody))
}

func (api *Token) Disable(id string) {
	body := api.GetBody(id, "false")

	req, err := http.NewRequest("POST", ApiUrl+"/v1.0/devices/actions", bytes.NewBuffer(body))

	if err != nil {
		fmt.Printf("OauthToken.Disable - http.NewRequest: %v", err)
		return
	}

	api.addHeaders(req.Header)

	res := sendSafeRequest(req)

	_, err = io.ReadAll(res.Body)

	if err != nil {
		fmt.Printf("OauthToken.Disable - io.ReadAll: %v", err)
		return
	}
}

func (api *Token) PrintDevices() string {
	devices := api.Info()

	fmt.Println("Please, choose your device number by number before \")\"")
	for i, d := range devices.Devices {
		fmt.Printf("%d) Name: %s, Id: %s\n", i, d.Name, d.Id)
	}

	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	id, err := strconv.Atoi(input.Text())

	if err != nil {
		fmt.Printf("OauthToken.PrintDevices - strconv.Atoi: %v", err)
		return ""
	}

	return devices.Devices[id].Id
}
