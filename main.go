package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"the-ice-time/auth"
	"the-ice-time/model"
	"the-ice-time/utils"
	"the-ice-time/workflows"

	"go.temporal.io/sdk/client"
)

func main() {
	environment := ""
	if environment == "" {
		log.Fatalln("Must specify a environment to perform DICE CMD")
	}

	cmd := ""
	if cmd == "" {
		log.Fatalln("Must specify a cmd to perform on DICE Job")
	}

	body := `{}`
	// body = `{"targetProjectIds": ["prep-2134-entdatalake-969cbf","qa-2134-entdatalake-d057be"],"jdbcTargets": []}`

	metaSvcUrl := ""
	chunkSize := 15

	switch environment {
	case model.Dev:
		metaSvcUrl = "https://dice-meta-svc-dot-dev-2367-entdataingst-5a9bf0.appspot.com"
	case model.QA:
		metaSvcUrl = "https://dice-meta-svc-dot-qa-2367-entdataingst-c1271b.appspot.com"
	case model.Preprod:
		metaSvcUrl = "https://dice-meta-svc-dot-prep-2367-entdataingst-804660.appspot.com"
	case model.Prod:
		metaSvcUrl = "https://dice-meta-svc-dot-prod-2367-entdataingst-7010d5.appspot.com"
	}

	if len(os.Args) <= 1 {
		log.Fatalln("Must specify a jobs.txt file path as the command-line argument")
	}
	filePath := os.Args[1]

	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("%s was not able to be read", filePath)
	}

	if len(fileData) == 0 {
		log.Fatalf("%s no data to read", filePath)
	}

	jobs := strings.Split(strings.ReplaceAll(string(fileData), "\r", ""), "\n")

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	chunkJobIDs := utils.ChunkJobs(jobs, chunkSize)

	for i := range chunkJobIDs {
		jobIDs := chunkJobIDs[i]

		bearer := auth.GetIdentityToken()

		wg := sync.WaitGroup{}
		wg.Add(len(jobIDs))

		for j := range jobIDs {
			index := (chunkSize * i) + j + 1
			dataSourceId := jobIDs[j]

			workflowID := fmt.Sprintf("%s-%s-%s", environment, cmd, dataSourceId)
			options := client.StartWorkflowOptions{
				ID:        workflowID,
				TaskQueue: model.DICEJobQueue,
			}

			go func() {
				workflowExecutor, err := temporalClient.ExecuteWorkflow(
					context.Background(),
					options,
					workflows.JobCMDWorkflow,
					index,
					dataSourceId,
					metaSvcUrl,
					body,
					bearer,
					cmd,
				)
				if err != nil {
					log.Fatalln("Unable to execute workflow", err)
				}

				result := ""
				err = workflowExecutor.Get(context.Background(), &result)
				if err != nil {
					log.Fatalln("Unable get workflow result", err)
				}

				fmt.Println(result)
				wg.Done()
			}()
		}

		wg.Wait()
	}

	fmt.Println("All jobs execution complete!")
}
