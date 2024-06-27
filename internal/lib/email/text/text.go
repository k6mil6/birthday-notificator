package text

import (
	"fmt"
	"time"
)

func BuildEmailHeader(name string) string {
	return fmt.Sprintf("%s birthday's notification", name)
}

func BuildEmailBody(name string, birthday, notificationDate time.Time) string {
	return fmt.Sprintf(
		"You've subscribed to %s's birthday notifications, we're reminding you that it's going to be on %02d.%02d.%d",
		name, birthday.Day(), birthday.Month(), notificationDate.Year(),
	)
}
