package main

import "time"

const (
	MON = 1
	TUE = 2
	WED = 3
	THU = 4
	FRI = 5
	SAT = 6
	SUN = 0
)

func getTomorrowsDayNumber() uint8 {
	currentTime := time.Now()

	tomorrow := currentTime.Add(24 * time.Hour)

	return uint8(tomorrow.Weekday())
}

func convertDayNumberToShortString(dayNumber uint8) string {
	var dayShortName string
	switch dayNumber {
	case MON:
		dayShortName = "m"
	case TUE:
		dayShortName = "t"
	case WED:
		dayShortName = "w"
	case THU:
		dayShortName = "th"
	case FRI:
		dayShortName = "f"
	case SAT:
		dayShortName = "s"
	case SUN:
		dayShortName = "su"
	}

	return dayShortName
}
