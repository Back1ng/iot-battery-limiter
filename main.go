package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/back1ng/iot-battery-limiter/internal/battery"
	"github.com/back1ng/iot-battery-limiter/internal/configuration"
	"github.com/back1ng/iot-battery-limiter/internal/yandex"
	"github.com/ztrue/shutdown"
)

func getPercentages() (int, int) {
	if os.Getenv("PERCENT_MAX") == "" || os.Getenv("PERCENT_MIN") == "" {
		panic("Percentage MAX and MIN must be set")
	}

	percentageMin := os.Getenv("PERCENT_MIN")
	percentageMax := os.Getenv("PERCENT_MAX")

	pmin, err := strconv.Atoi(percentageMin)

	if err != nil {
		// handle error
	}

	pmax, err := strconv.Atoi(percentageMax)

	if err != nil {
		// handle error
	}

	if pmin >= pmax {
		panic("Minimum percentage must be not great than Maximum percentage")
	}

	return pmin, pmax
}

func main() {
	configuration.GenerateEnv()

	if os.Getenv("YANDEX_OAUTH") == "" {
		panic("Yandex oauth key must be not empty!")
	}

	pmin, pmax := getPercentages()

	token := yandex.OauthToken{Token: "Bearer " + os.Getenv("YANDEX_OAUTH")}

	device := os.Getenv("DEVICE_ID")

	if device == "" {
		token.PrintDevices()
		device = os.Getenv("DEVICE_ID")
	}

	shutdown.Add(func() {
		token.Disable(device)
	})

	go func() {
		for {
			percent, err := battery.Get()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(percent)

			if percent >= pmax {
				fmt.Println("got disable percentage")

				isCharging, err := battery.Charging()

				if err != nil {
					fmt.Println(err)
				}

				if isCharging {
					token.Disable(device)
				}
			} else if percent <= pmin {
				fmt.Println("got enable percentage")

				isDischarging, err := battery.Discharging()

				if err != nil {
					fmt.Println(err)
				}

				if isDischarging {
					token.Enable(device)
				}
			}

			time.Sleep(time.Second * 30)
		}
	}()

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
