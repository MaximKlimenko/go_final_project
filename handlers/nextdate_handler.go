package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/MaximKlimenko/go_final_project/nextdate"
	"github.com/MaximKlimenko/go_final_project/utils"
)

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	now, err := time.Parse(utils.TimeFormat, nowStr)
	if err != nil {
		http.Error(w, "некорректная дата 'now'", http.StatusBadRequest)
		return
	}

	nextDate, err := nextdate.NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, nextDate)
}
