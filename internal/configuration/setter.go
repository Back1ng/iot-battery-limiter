package configuration

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"github.com/back1ng/iot-battery-limiter/internal/yandex"
	"github.com/joho/godotenv"
	"github.com/pkg/browser"
)

func ReadInput() string {
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	return input.Text()
}

func GenerateEnv() {
	_ = godotenv.Load()

	conf := NewConfiguration()

	if conf.Percent.Max == 0 {
		fmt.Println("Please, choose maximum percentage of battery")

		max, _ := strconv.Atoi(ReadInput())
		conf.Percent.Max = max
	}

	if conf.Percent.Min == 0 {
		fmt.Println("Please, choose minimum percentage of battery")

		min, _ := strconv.Atoi(ReadInput())
		conf.Percent.Min = min
	}

	if len(conf.YandexOauth.Token) == 0 {
		fmt.Println("Next, you need login at yandex on the next link and write code.")
		fmt.Println("Press any key, when you ready to open browser.")

		ReadInput()

		browser.OpenURL("https://oauth.yandex.ru/authorize?response_type=token&client_id=1e9f9d503a8749619d4a0539c1e7bb0c")

		conf.YandexOauth = yandex.OauthToken{
			Token: ReadInput(),
		}
	}

	godotenv.Write(conf.ToMap(), "./.env")
	godotenv.Load()
}
