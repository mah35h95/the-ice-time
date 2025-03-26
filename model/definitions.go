package model

const (
	Dev     string = "dev"
	QA      string = "qa"
	Preprod string = "preprod"
	Prod    string = "prod"
)

// CronRange - Has a min max value for a cron string
type CronRange struct {
	Min  int
	Max  int
	Cron string
}
