package utils

import "the-ice-time/model"

// GetCron - returns an increasing cron string
func GetCron(value int) (string, string) {
	cronTimeZone := "America/Chicago"

	cronRanges := []model.CronRange{
		{Min: 1, Max: 50, Cron: "0 0 * * *"},
		{Min: 51, Max: 100, Cron: "30 0 * * *"},
		{Min: 101, Max: 150, Cron: "0 1 * * *"},
		{Min: 151, Max: 200, Cron: "30 1 * * *"},
		{Min: 201, Max: 250, Cron: "0 2 * * *"},
		{Min: 251, Max: 300, Cron: "30 2 * * *"},
		{Min: 301, Max: 350, Cron: "0 3 * * *"},
		{Min: 351, Max: 400, Cron: "30 3 * * *"},
		{Min: 401, Max: 450, Cron: "0 4 * * *"},
		{Min: 451, Max: 500, Cron: "30 4 * * *"},
		{Min: 501, Max: 550, Cron: "0 5 * * *"},
		{Min: 551, Max: 600, Cron: "30 5 * * *"},
		{Min: 601, Max: 650, Cron: "0 6 * * *"},
		{Min: 651, Max: 700, Cron: "30 6 * * *"},
		{Min: 701, Max: 750, Cron: "0 7 * * *"},
		{Min: 751, Max: 800, Cron: "30 7 * * *"},
		{Min: 801, Max: 850, Cron: "0 8 * * *"},
		{Min: 851, Max: 900, Cron: "30 8 * * *"},
	}

	for _, cronRange := range cronRanges {
		if value >= cronRange.Min && value <= cronRange.Max {
			return cronRange.Cron, cronTimeZone
		}
	}

	return "0 0 * * *", cronTimeZone
}
