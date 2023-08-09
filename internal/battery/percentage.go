package battery

import (
	"github.com/distatus/battery"
)

func Get() (int, error) {
	battery, err := battery.Get(0)

	if err != nil {
		return 0, err
	}

	return int(battery.Current / battery.Full * 100), nil
}

func Charging() (bool, error) {
	battery, err := battery.Get(0)

	if err != nil {
		return false, err
	}

	return battery.State.String() == "Charging", nil
}

func Discharging() (bool, error) {
	battery, err := battery.Get(0)

	if err != nil {
		return false, err
	}

	return battery.State.String() == "Discharging", nil
}
