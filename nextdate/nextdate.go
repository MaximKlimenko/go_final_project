package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MaximKlimenko/go_final_project/utils"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("пустое правило повторения")
	}
	//Проверка формата даты
	parsedDate, err := time.Parse(utils.TimeFormat, date)
	if err != nil {
		return "", fmt.Errorf("некорректная дата: %w", err)
	}

	for {
		if repeat == "y" {
			parsedDate = parsedDate.AddDate(1, 0, 0)
		} else if strings.HasPrefix(repeat, "d ") {
			daysStr := strings.TrimPrefix(repeat, "d ")
			days, err := strconv.Atoi(daysStr)
			if err != nil || days < 1 || days > 400 {
				return "", fmt.Errorf("некорректное правило: d %s", daysStr)
			}
			parsedDate = parsedDate.AddDate(0, 0, days)
		} else {
			return "", fmt.Errorf("недопустимый формат правила: %s", repeat)
		}
		//Проверка возвращаемой даты
		if parsedDate.After(now) {
			return parsedDate.Format(utils.TimeFormat), nil
		}
	}
}
