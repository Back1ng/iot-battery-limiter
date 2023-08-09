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
	YandexOauth yandex.OauthToken
	DeviceId    int
}

func NewConfiguration() Configuration {
	_ = godotenv.Load()

	pmin, _ := strconv.Atoi(os.Getenv("PERCENT_MIN"))
	pmax, _ := strconv.Atoi(os.Getenv("PERCENT_MAX"))
	devId, _ := strconv.Atoi(os.Getenv("DEVICE_ID"))

	return Configuration{
		Percent: Percent{
			Min: pmin,
			Max: pmax,
		},
		YandexOauth: yandex.OauthToken{
			Token: os.Getenv("YANDEX_OAUTH"),
		},
		DeviceId: devId,
	}
}

func (c Configuration) ToMap() map[string]string {
	return map[string]string{
		"PERCENT_MIN":  fmt.Sprint(c.Percent.Min),
		"PERCENT_MAX":  fmt.Sprint(c.Percent.Max),
		"YANDEX_OAUTH": c.YandexOauth.Token,
		"DEVICE_ID":    fmt.Sprint(c.DeviceId),
	}
}