package battery

import (
	"github.com/distatus/battery"
)

func Get() (int, error) {
	b, err := battery.Get(0)

	if err != nil {
		return 0, err
	}

	return int(b.Current / b.Full * 100), nil
}

func Charging() (bool, error) {
	b, err := battery.Get(0)

	if err != nil {
		return false, err
	}

	return b.State.String() == "Charging", nil
}

func Discharging() (bool, error) {
	b, err := battery.Get(0)

	if err != nil {
		return false, err
	}

	return b.State.String() == "Discharging", nil
}
