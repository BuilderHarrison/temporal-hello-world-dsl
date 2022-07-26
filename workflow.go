package app

import (
	"time"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"go.temporal.io/sdk/workflow"
)

type OperationState model.OperationState

// HelloWorld workflow definition
func HelloWorld(ctx workflow.Context, dslWorkflow *model.Workflow) ([]byte, error) {
	logger := workflow.GetLogger(ctx)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for _, state := range dslWorkflow.States {
		switch b := state.(type) {
		case *model.OperationState:
			err := executeOperation(ctx, b)
			if err != nil {
				logger.Error("DSL Workflow failed.", "Error", err)
				return nil, err
			}
		default:
			logger.Error("Unkown state")
		}
	}

	logger.Info("DSL Workflow completed.")

	return nil, nil
}

func executeOperation(ctx workflow.Context, b *model.OperationState) error {
	if b.ActionMode == "parallel" {
		err := executeParallel(ctx, b.Actions)
		if err != nil {
			return err
		}
	}

	if b.ActionMode == "sequential" {
		err := executeSequence(ctx, b.Actions)
		if err != nil {
			return err
		}
	}

	return nil
}

func executeSequence(ctx workflow.Context, actions []model.Action) error {
	for _, a := range actions {
		err := workflow.ExecuteActivity(ctx, a.FunctionRef.RefName, a.FunctionRef.Arguments).Get(ctx, nil)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeParallel(ctx workflow.Context, actions []model.Action) error {
	//
	// You can use the context passed in to activity as a way to cancel the activity like standard GO way.
	// Cancelling a parent context will cancel all the derived contexts as well.
	//

	// In the parallel block, we want to execute all of them in parallel and wait for all of them.
	// if one activity fails then we want to cancel all the rest of them as well.
	childCtx, cancelHandler := workflow.WithCancel(ctx)
	selector := workflow.NewSelector(ctx)
	var activityErr error
	for _, s := range actions {
		f := executeAsync(s, childCtx)
		selector.AddFuture(f, func(f workflow.Future) {
			err := f.Get(ctx, nil)
			if err != nil {
				// cancel all pending activities
				cancelHandler()
				activityErr = err
			}
		})
	}

	for i := 0; i < len(actions); i++ {
		selector.Select(ctx) // this will wait for one branch
		if activityErr != nil {
			return activityErr
		}
	}

	return nil
}

func executeAsync(action model.Action, ctx workflow.Context) workflow.Future {
	future, settable := workflow.NewFuture(ctx)
	workflow.Go(ctx, func(ctx workflow.Context) {
		err := workflow.ExecuteActivity(ctx, action.FunctionRef.RefName, action.FunctionRef.Arguments).Get(ctx, nil)
		settable.Set(nil, err)
	})
	return future
}
