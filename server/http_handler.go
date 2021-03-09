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
	if diff.Hours() > 31*24 {
		return time.Time{}, time.Time{}, errors.New("date range cannot be greater than a month")
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
