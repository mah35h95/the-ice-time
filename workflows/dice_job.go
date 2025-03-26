package workflows

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"the-ice-time/auth"
	"the-ice-time/dice"
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
	case dice.Pause:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			dice.Pause,
			body,
		).Get(ctx, nil)

	case dice.Resume:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			dice.Resume,
			body,
		).Get(ctx, nil)

	case dice.Stop:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			dice.Stop,
			body,
		).Get(ctx, nil)

	case dice.Load:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			dice.Load,
			body,
		).Get(ctx, nil)

	case dice.Lock:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			dice.Lock,
			body,
		).Get(ctx, nil)

	case dice.Unlock:
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			dice.Unlock,
			body,
		).Get(ctx, nil)

	case dice.Reload:
		body = `{"keepFoundryDataset": true,"retainData": false}`
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			dice.Reload,
			body,
		).Get(ctx, nil)

	case dice.EditGCPTarget:
		// body = `{"targetProjectIds": ["prep-2134-entdatalake-969cbf","qa-2134-entdatalake-d057be"],"jdbcTargets": []}`
		err = workflow.ExecuteActivity(
			ctx,
			dice.ExecuteJobCmd,
			dataSourceId,
			metaSvcUrl,
			bearer,
			http.MethodPost,
			dice.Edit,
			body,
		).Get(ctx, nil)

	case dice.Delete:
		err = workflow.ExecuteActivity(
			ctx,
			dice.DeleteJob,
			dataSourceId,
			metaSvcUrl,
			bearer,
		).Get(ctx, nil)

	case dice.EditCron:
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

	case dice.DeleteHydratedRes:
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
