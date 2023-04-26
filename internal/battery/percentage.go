package battery

import "github.com/distatus/battery"

func Get() int {
	battery, err := battery.Get(0)

	if err != nil {
		// handle error
	}

	return int(battery.Current / battery.Full * 100)
}
