package birthday

import "time"

func CalculateNextBirthday(birthday time.Time) time.Time {
	now := time.Now()
	thisYearBirthday := time.Date(now.Year(), birthday.Month(), birthday.Day(), 0, 0, 0, 0, time.Local)
	if now.Before(thisYearBirthday) {
		return thisYearBirthday
	}
	return thisYearBirthday.AddDate(1, 0, 0)
}
