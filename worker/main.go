package main

import (
	"log"

	"the-ice-time/dice"
	"the-ice-time/model"
	"the-ice-time/workflows"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Create the client object just once per process
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer temporalClient.Close()

	// This worker hosts both Workflow and Activity functions
	diceWorker := worker.New(temporalClient, model.DICEJobQueue, worker.Options{})
	diceWorker.RegisterWorkflow(workflows.JobCMDWorkflow)
	diceWorker.RegisterActivity(dice.ExecuteJobCmd)
	diceWorker.RegisterActivity(dice.DeleteJob)
	diceWorker.RegisterActivity(dice.EditCronSchedule)
	diceWorker.RegisterActivity(dice.DeleteHydratedResources)

	// Start listening to the Task Queue
	err = diceWorker.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
