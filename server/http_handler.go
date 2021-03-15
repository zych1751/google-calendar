package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func parseTime(startTimeStr string, endTimeStr string) (time.Time, time.Time, error) {
	startTime, err := time.Parse(time.RFC3339, startTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("startTime is not valid (time format have to be RFC3339)")
	}
	endTime, err := time.Parse(time.RFC3339, endTimeStr)
	if err != nil {
		return time.Time{}, time.Time{}, errors.New("endTime is not valid (time format have to be RFC3339)")
	}

	diff := endTime.Sub(startTime)
	if diff.Hours() < 0 {
		return time.Time{}, time.Time{}, errors.New("date range cannot be negative")
	}

	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()

	firstDayOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	diff = startTime.Sub(firstDayOfMonth)
	if diff.Hours() < -15*24 {
		return time.Time{}, time.Time{}, errors.New("startDate cannot be lower than two weeks before the first day of current Month")
	}
	diff = endTime.Sub(lastDayOfMonth)
	if diff.Hours() > 15*24 {
		return time.Time{}, time.Time{}, errors.New("endDate cannot be greater than two weeks after the last day of current Month")
	}

	return startTime, endTime, nil
}

func getScheduleHandler(googleClient *GoogleClient) func(echo.Context) error {
	return func(c echo.Context) error {
		startTimeStr := c.QueryParam("startTime")
		endTimeStr := c.QueryParam("endTime")

		startTime, endTime, err := parseTime(startTimeStr, endTimeStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		items, err := googleClient.GetSchedule(startTime, endTime)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		return c.JSON(http.StatusOK, items)
	}
}
