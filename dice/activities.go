package dice

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/sjson"
)

const (
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

// ValidateToken - validates token with dice meta api
func ValidateToken(metaSvcUrl, bearer string) error {
	path := fmt.Sprintf("%v/", metaSvcUrl)

	request, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %v", err)
	}

	request.Header = http.Header{
		"Authorization": {bearer},
		"Content-Type":  {"application/json"},
		"tyson-user":    {"qtpie"},
	}

	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do: %v", err)
	}
	if response.StatusCode == 403 {
		fmt.Println("Token is Invalid")
		return fmt.Errorf("403")
	}
	defer response.Body.Close()

	fmt.Println("Token is Valid")
	return nil
}

// ExecuteJobCmd - Executes the dice api call for the give http method, cmd and body
func ExecuteJobCmd(dataSourceId, metaSvcUrl, bearer, httpMethod, cmd, stringBody string) error {
	parts := strings.Split(dataSourceId, ".")
	if len(parts) != 5 {
		return errors.New("invalid dataSourceId " + dataSourceId)
	}

	path := fmt.Sprintf(
		"%s/sources/%s/technologies/%s/databases/%s/jobs/%s.%s/%s",
		metaSvcUrl,
		parts[0],
		parts[1],
		parts[2],
		parts[3],
		parts[4],
		cmd,
	)

	body := strings.NewReader(stringBody)

	request, err := http.NewRequest(httpMethod, path, body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %v", err)
	}

	request.Header = http.Header{
		"Authorization": {bearer},
		"Content-Type":  {"application/json"},
	}

	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do: %v", err)
	}
	if response.StatusCode == 403 {
		return fmt.Errorf("403")
	}
	defer response.Body.Close()

	fmt.Printf("Job %s has been triggered to be %s.\n", dataSourceId, cmd)

	return nil
}

// DeleteJob - Deletes DICE job
func DeleteJob(dataSourceId string, metaSvcUrl string, bearer string) error {
	parts := strings.Split(dataSourceId, ".")
	if len(parts) != 5 {
		return errors.New("invalid dataSourceId " + dataSourceId)
	}

	path := fmt.Sprintf(
		"%s/sources/%s/technologies/%s/databases/%s/jobs/%s.%s",
		metaSvcUrl,
		parts[0],
		parts[1],
		parts[2],
		parts[3],
		parts[4],
	)

	body := strings.NewReader(`{}`)

	request, err := http.NewRequest(http.MethodDelete, path, body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %v", err)
	}

	request.Header = http.Header{
		"Authorization": {bearer},
		"Content-Type":  {"application/json"},
	}

	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do: %v", err)
	}
	if response.StatusCode == 403 {
		return fmt.Errorf("403")
	}
	defer response.Body.Close()

	fmt.Printf("Job %s has been triggered to be deleted.\n", dataSourceId)

	return nil
}

// EditCronSchedule - edits jobs to have the given cron schedule
func EditCronSchedule(dataSourceId, metaSvcUrl, bearer, cron, cronTimeZone string) error {
	parts := strings.Split(dataSourceId, ".")
	if len(parts) != 5 {
		return errors.New("invalid dataSourceId " + dataSourceId)
	}

	path := fmt.Sprintf("%s/sources/%s/technologies/%s/databases/%s/jobs/%s.%s",
		metaSvcUrl,
		parts[0],
		parts[1],
		parts[2],
		parts[3],
		parts[4],
	)

	client := &http.Client{}

	fmt.Printf("Getting job data of %s\n", dataSourceId)

	request1, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest get: %v", err)
	}

	request1.Header = http.Header{
		"Authorization": {bearer},
	}
	response1, err := client.Do(request1)
	if err != nil {
		return fmt.Errorf("client.Do get: %v", err)
	}
	if response1.StatusCode == 403 {
		return fmt.Errorf("403")
	}

	body, err := io.ReadAll(response1.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body. %v", err)
	}
	defer response1.Body.Close()

	newScheduleValue, err := sjson.Set(string(body), "schedule", cron)
	if err != nil {
		return fmt.Errorf("failed to update json value. %v", err)
	}

	newScheduleTimeZone, err := sjson.Set(newScheduleValue, "cronTimezone", cronTimeZone)
	if err != nil {
		return fmt.Errorf("failed to update json value. %v", err)
	}

	request2, err := http.NewRequest(http.MethodPost, path+"/edit", bytes.NewBuffer([]byte(newScheduleTimeZone)))
	if err != nil {
		return fmt.Errorf("http.NewRequest post: %v", err)
	}

	request2.Header = http.Header{
		"Authorization": {bearer},
		"Content-Type":  {"application/json"},
		"Accept":        {"*/*"},
	}
	response2, err := client.Do(request2)
	if err != nil {
		return fmt.Errorf("client.Do post: %v", err)
	}
	if response2.StatusCode == 403 {
		return fmt.Errorf("403")
	}
	defer response2.Body.Close()

	fmt.Printf("Job %s cron has been triggered to be changed.\n", dataSourceId)

	return nil
}

// DeleteHydratedResources - Deletes Hydrated tables
func DeleteHydratedResources(dataSourceId string, metaSvcUrl string, bearer string) error {
	parts := strings.Split(dataSourceId, ".")
	if len(parts) != 5 {
		return errors.New("invalid dataSourceId " + dataSourceId)
	}

	path := fmt.Sprintf(
		"%s/sources/%s/technologies/%s/databases/%s/jobs/%s.%s/delete_hydrated_resources",
		metaSvcUrl,
		parts[0],
		parts[1],
		parts[2],
		parts[3],
		parts[4],
	)

	body := strings.NewReader(`{}`)

	request, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %v", err)
	}

	request.Header = http.Header{
		"Authorization": {bearer},
		"Content-Type":  {"application/json"},
	}

	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("client.Do: %v", err)
	}
	if response.StatusCode == 403 {
		return fmt.Errorf("403")
	}
	defer response.Body.Close()

	fmt.Printf("Job %s has been triggered to clean up the hydrated resources.\n", dataSourceId)

	return nil
}
