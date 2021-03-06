package utils

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/maddevsio/comedian/model"
)

//SplitUser divides full user object to name & id
func SplitUser(user string) (string, string) {
	userFull := strings.Split(user, "|")
	userID := strings.Replace(userFull[0], "<@", "", -1)
	userName := strings.Replace(userFull[1], ">", "", -1)
	return userID, userName
}

// FormatTime returns hour and minutes from string
func FormatTime(t string) (hour, min int, err error) {
	newErr := errors.New("time format error")
	ts := strings.Split(t, ":")
	if len(ts) != 2 {
		return 0, 0, newErr
	}

	hour, err = strconv.Atoi(ts[0])
	if err != nil {
		return 0, 0, newErr
	}
	min, err = strconv.Atoi(ts[1])
	if err != nil {
		return 0, 0, newErr
	}

	if hour < 0 || hour > 23 || min < 0 || min > 59 {
		return 0, 0, newErr
	}
	return hour, min, nil
}

//SplitTimeTalbeCommand returns set of substrings
func SplitTimeTalbeCommand(t, on, at string) (string, string, int64, error) {

	a := strings.Split(t, on)
	if len(a) != 2 {
		return "", "", int64(0), errors.New("Sorry, could not understand where are the standupers and where is the rest of the command. Please, check the text for mistakes and try again")
	}
	users := strings.TrimSpace(a[0])
	b := strings.Split(a[1], at)
	if len(a) != 2 {
		return "", "", int64(0), errors.New("Sorry, could not understand where are the weekdays and where is the time. Please, check the text for mistakes and try again")
	}
	weekdays := strings.ToLower(strings.TrimSpace(b[0]))
	timeText := strings.ToLower(strings.TrimSpace(b[1]))
	time, err := ParseTimeTextToInt(timeText)
	if err != nil {
		return "", "", int64(0), err
	}
	reg, err := regexp.Compile("[^а-яА-Яa-zA-Z0-9]+")
	if err != nil {
		return "", "", int64(0), err
	}
	weekdays = reg.ReplaceAllString(weekdays, " ")
	return users, weekdays, time, nil
}

func ParseTimeTextToInt(timeText string) (int64, error) {
	if timeText == "0" {
		return int64(0), nil
	}
	matchHourMinuteFormat, _ := regexp.MatchString("[0-9][0-9]:[0-9][0-9]", timeText)
	matchAMPMFormat, _ := regexp.MatchString("[0-9][0-9][a-z]", timeText)

	if matchHourMinuteFormat {
		t := strings.Split(timeText, ":")
		hours, _ := strconv.Atoi(t[0])
		munites, _ := strconv.Atoi(t[1])

		if hours > 23 || munites > 59 {
			return int64(0), errors.New("Wrong time! Please, check the time format and try again!")
		}
		return time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), hours, munites, 0, 0, time.Local).Unix(), nil
	}

	if matchAMPMFormat {
		return int64(0), errors.New("Seems like you used short time format, please, use 24:00 hour format instead!")
	}

	return int64(0), errors.New("Could not understand how you mention time. Please, use 24:00 hour format and try again!")

}

func PrepareTimeTable(tt model.TimeTable, weekdays string, timeInt int64) model.TimeTable {
	if strings.Contains(weekdays, "mon") || strings.Contains(weekdays, "пн") {
		tt.Monday = timeInt
	}
	if strings.Contains(weekdays, "tue") || strings.Contains(weekdays, "вт") {
		tt.Tuesday = timeInt
	}
	if strings.Contains(weekdays, "wed") || strings.Contains(weekdays, "ср") {
		tt.Wednesday = timeInt
	}
	if strings.Contains(weekdays, "thu") || strings.Contains(weekdays, "чт") {
		tt.Thursday = timeInt
	}
	if strings.Contains(weekdays, "fri") || strings.Contains(weekdays, "пт") {
		tt.Friday = timeInt
	}
	if strings.Contains(weekdays, "sat") || strings.Contains(weekdays, "сб") {
		tt.Saturday = timeInt
	}
	if strings.Contains(weekdays, "sun") || strings.Contains(weekdays, "вс") {
		tt.Sunday = timeInt
	}
	return tt
}

//SetupDays gets dates and returns their differense in days
func SetupDays(dateFrom, dateTo time.Time) (time.Time, int, error) {
	if dateTo.Before(dateFrom) {
		return time.Now(), 0, errors.New("date to is less than date from")
	}
	if dateTo.After(time.Now()) {
		return time.Now(), 0, errors.New("date to is in the future")
	}
	dateFromRounded := time.Date(dateFrom.Year(), dateFrom.Month(), dateFrom.Day(), 0, 0, 0, 0, time.UTC)
	dateToRounded := time.Date(dateTo.Year(), dateTo.Month(), dateTo.Day(), 0, 0, 0, 0, time.UTC)
	numberOfDays := int(dateToRounded.Sub(dateFromRounded).Hours() / 24)
	return dateFromRounded, numberOfDays, nil
}
