package app

import (
	"context"
	"fmt"
)

type DslActivities struct {
}

func (a *DslActivities) PrintHelloWorld1(ctx context.Context, args map[string]interface{}) (bool, error) {
	fmt.Printf("\n%s\n", args["arg1"])
	return true, nil
}

func (a *DslActivities) PrintHelloWorld2(ctx context.Context, args map[string]interface{}) (bool, error) {
	fmt.Printf("\n%s\n", args["arg1"])
	return true, nil
}
