package main

import (
	"fmt"
	"syscall"
	"time"

	"github.com/back1ng/iot-battery-limiter/internal/battery"
	"github.com/back1ng/iot-battery-limiter/internal/configuration"
	"github.com/ztrue/shutdown"
)

func main() {
	conf := configuration.GenerateEnv()

	shutdown.Add(func() {
		conf.YandexOauth.Disable(conf.DeviceId)
	})

	go func() {
		for {
			percent, err := battery.Get()

			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(percent)

			if percent >= conf.Percent.Max {
				fmt.Println("got disable percentage")

				isCharging, err := battery.Charging()

				if err != nil {
					fmt.Println(err)
				}

				if isCharging {
					conf.YandexOauth.Disable(conf.DeviceId)
				}
			} else if percent <= conf.Percent.Min {
				fmt.Println("got enable percentage")

				isDischarging, err := battery.Discharging()

				if err != nil {
					fmt.Println(err)
				}

				if isDischarging {
					conf.YandexOauth.Enable(conf.DeviceId)
				}
			}

			time.Sleep(time.Second * 30)
		}
	}()

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM)
}
