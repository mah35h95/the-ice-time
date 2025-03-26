package model

const (
	Dev     string = "dev"
	QA      string = "qa"
	Preprod string = "preprod"
	Prod    string = "prod"

	Pause             string = "pause"
	Stop              string = "stop"
	Resume            string = "resume"
	Load              string = "load"
	Lock              string = "lock"
	Unlock            string = "unlock"
	Reload            string = "reload"
	Delete            string = "delete"
	DeleteHydratedRes string = "delete_hydrated_resources"
	Edit              string = "edit"
	EditCron          string = "edit_cron"
	EditGCPTarget     string = "edit_gcp_target"
)

// CronRange - Has a min max value for a cron string
type CronRange struct {
	Min  int
	Max  int
	Cron string
}
