package xtime

import "time"

const layoutHoursMinutes = "15:04"

func ParseHoursMinutesFromString(strTime string) (time.Time, error) {
	t, err := time.Parse(layoutHoursMinutes, strTime)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
