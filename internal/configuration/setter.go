package configuration

import (
	"bufio"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/pkg/browser"
)

func GenerateEnv() {
	_ = godotenv.Load()

	conf := make(map[string]string, 4)

	if os.Getenv("PERCENT_MAX") == "" {
		fmt.Println("Please, choose maximum percentage of battery")

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		conf["PERCENT_MAX"] = input.Text()
	} else {
		conf["PERCENT_MAX"] = os.Getenv("PERCENT_MAX")
	}

	if os.Getenv("PERCENT_MIN") == "" {
		fmt.Println("Please, choose minimum percentage of battery")

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		conf["PERCENT_MIN"] = input.Text()
	} else {
		conf["PERCENT_MIN"] = os.Getenv("PERCENT_MIN")
	}

	if os.Getenv("YANDEX_OAUTH") == "" {
		fmt.Println("Next, you need login at yandex on the next link and write code.")
		fmt.Println("Press any key, when you ready to open browser.")

		input := bufio.NewScanner(os.Stdin)
		input.Scan()

		browser.OpenURL("https://oauth.yandex.ru/authorize?response_type=token&client_id=1e9f9d503a8749619d4a0539c1e7bb0c")

		fmt.Println("Please, enter token from browser.")
		input.Scan()
		conf["YANDEX_OAUTH"] = input.Text()
	} else {
		conf["YANDEX_OAUTH"] = os.Getenv("YANDEX_OAUTH")
	}

	if os.Getenv("DEVICE_ID") != "" {
		conf["DEVICE_ID"] = os.Getenv("DEVICE_ID")
	}

	godotenv.Write(conf, "./.env")
	godotenv.Load()
}
