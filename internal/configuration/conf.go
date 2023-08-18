package configuration

import (
	"fmt"
	"os"
	"strconv"

	"github.com/back1ng/iot-battery-limiter/internal/yandex"
	"github.com/joho/godotenv"
)

type Percent struct {
	Min, Max int
}

type Configuration struct {
	Percent     Percent
	YandexOauth yandex.Token
	DeviceId    string
}

func NewConfiguration() Configuration {
	_ = godotenv.Load()

	pmin, _ := strconv.Atoi(os.Getenv("PERCENT_MIN"))
	pmax, _ := strconv.Atoi(os.Getenv("PERCENT_MAX"))

	return Configuration{
		Percent: Percent{
			Min: pmin,
			Max: pmax,
		},
		YandexOauth: yandex.Token(os.Getenv("YANDEX_OAUTH")),
		DeviceId:    os.Getenv("DEVICE_ID"),
	}
}

func (c Configuration) ToMap() map[string]string {
	return map[string]string{
		"PERCENT_MIN":  fmt.Sprint(c.Percent.Min),
		"PERCENT_MAX":  fmt.Sprint(c.Percent.Max),
		"YANDEX_OAUTH": string(c.YandexOauth),
		"DEVICE_ID":    c.DeviceId,
	}
}
