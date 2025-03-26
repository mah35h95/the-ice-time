package workflows

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"the-ice-time/auth"
	"the-ice-time/dice"
	"the-ice-time/model"
	"the-ice-time/utils"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func JobCMDWorkflow(
	ctx workflow.Context,
	index int,
	dataSourceId,
	metaSvcUrl,
	body,
	bearer,
	cmd string,
) (string, error) {
	// Define the activity options, including the retry policy
	activityOptions := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute * 5,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second, // amount of time that must elapse before the first retry occurs
			MaximumInterval:    time.Minute, // maximum interval between retries
			BackoffCoefficient: 2,           // how much the retry interval increases
			MaximumAttempts:    5,           // Uncomment this if you want to limit attempts
		},
	}
	ctx = workflow.WithActivityOptions(ctx, activityOptions)

	err := cmdActivity(ctx, index, dataSourceId, metaSvcUrl, body, bearer, cmd)
	if err.Error() == "403" {
		bearer = auth.GetIdentityToken()
		err = cmdActivity(ctx, index, dataSourceId, metaSvcUrl, body, bearer, cmd)
	}

	if err != nil {
		return "", fmt.Errorf("failed to execute DICE CMD: %+v", err)
	}

	return fmt.Sprintf("%s DICE Job %s Completed", dataSourceId, cmd), nil
}

func cmdActivity(ctx workflow.Context, index int, dataSourceId, metaSvcUrl, body, bearer, cmd string) error {
	// body = `{}`
	err := error(nil)

	switch cmd {
	case model.Pause:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			model.Pause,
			body,
		).Get(ctx, nil)

	case model.Resume:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			model.Resume,
			body,
		).Get(ctx, nil)

	case model.Stop:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			model.Stop,
			body,
		).Get(ctx, nil)

	case model.Load:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			model.Load,
			body,
		).Get(ctx, nil)

	case model.Lock:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			model.Lock,
			body,
		).Get(ctx, nil)

	case model.Unlock:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			model.Unlock,
			body,
		).Get(ctx, nil)

	case model.Reload:
		body = `{"keepFoundryDataset": true,"retainData": false}`
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			model.Reload,
			body,
		).Get(ctx, nil)

	case model.EditGCPTarget:
		// body = `{"targetProjectIds": ["prep-2134-entdatalake-969cbf","qa-2134-entdatalake-d057be"],"jdbcTargets": []}`
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			model.Edit,
			body,
		).Get(ctx, nil)

	case model.Delete:
		err = workflow.ExecuteActivity(
			ctx,
			dice.DeleteJob,
			dataSourceId,
			metaSvcUrl,
			bearer,
		).Get(ctx, nil)

	case model.EditCron:
		cron, cronTimeZone := utils.GetCron(index)
		err = workflow.ExecuteActivity(
			ctx,
			dice.EditCronSchedule,
			dataSourceId,
			metaSvcUrl,
			bearer,
			cron,
			cronTimeZone,
		).Get(ctx, nil)

	case model.DeleteHydratedRes:
		err = workflow.ExecuteActivity(
			ctx,
			dice.DeleteHydratedResources,
			dataSourceId,
			metaSvcUrl,
			bearer,
		).Get(ctx, nil)

	default:
		err = errors.New("CMD provided does not match with predefined cases")
		log.Fatalln("unable to run DICE CMD activity", err)
	}

	return err
}
