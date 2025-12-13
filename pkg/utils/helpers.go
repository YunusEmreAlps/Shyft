package utils

import (
	"shyft/internal/models"
	"strings"
	"time"
)

func UrlStringToOptions(url string) (string, string, string, string, string, string) {
	// parse redis url to username password and host and port and db
	// app_name://username:password@host:port/db

	// split all options (app_name, username, password, host, port, db)
	options := strings.Split(url, "://")
	// split protocol and info
	protocol := options[0]
	info := options[1]
	// split info to username, password, host, port, db
	infoOptions := strings.Split(info, "@")
	// split username and password
	usernamePassword := infoOptions[0]
	hostPortDb := infoOptions[1]
	// split username and password
	usernamePasswordOptions := strings.Split(usernamePassword, ":")
	username := usernamePasswordOptions[0]
	password := usernamePasswordOptions[1]
	// split host, port and db
	hostPortDbOptions := strings.Split(hostPortDb, "/")
	hostPort := hostPortDbOptions[0]
	db := hostPortDbOptions[1]
	// split host and port
	hostPortOptions := strings.Split(hostPort, ":")
	host := hostPortOptions[0]
	port := hostPortOptions[1]

	return protocol, username, password, host, port, db
}

// The function calculates and sets the start and end dates for a series of shift based on user input.
func CalculateSentryShiftTimes(startDate time.Time, frequency int, shifts []models.Shift) []models.Shift {
	start := startDate
	end := start.AddDate(0, 0, frequency)

	for _, shift := range shifts {
		// set start and end time for each shift
		shift.Start = start.Format("2006-01-02 15:04:05")
		shift.End = end.Format("2006-01-02 15:04:05")

		// set start and end time for next shift
		start = end
		end = start.AddDate(0, 0, frequency)
	}

	return shifts
}
