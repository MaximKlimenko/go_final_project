package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	stringsRepeat := strings.Split(repeat, " ")
	var endDate string
	dateTime, err := time.Parse("20240229", date)
	if err != nil {
		return "", fmt.Errorf("Ошибка преобразования даты")
	}
	switch stringsRepeat[0] {
	case "":
		return "", fmt.Errorf("There's no repeat")
	case "d":
		days, err := strconv.Atoi(stringsRepeat[1])
		if err != nil {
			return "", fmt.Errorf("You have to write a number of days!! *angry smile*")
		}

		if days > 400 {
			return "", fmt.Errorf("Out of range days > 400")
		}
		endDate = dateTime.AddDate(0, 0, days).String()
	case "y":
		endDate = dateTime.AddDate(1, 0, 0).String()
	default:
		return "", fmt.Errorf("Something went wrong!")
	}
	if now.After(dateTime) {
		return "", fmt.Errorf("Возвращаемая дата должна быть больше даты, указанной в переменной now.")
	}
	return endDate, nil
}
