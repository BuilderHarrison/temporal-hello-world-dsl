package main

import (
	"context"
	"flag"
	"log"

	"github.com/serverlessworkflow/sdk-go/v2/parser"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"

	"temporal-poc-hello-world-dsl/app"
)

func main() {
	var dslConfig string
	flag.StringVar(&dslConfig, "dslConfig", "sequential.json", "dslConfig specify the json file for the dsl workflow.")
	flag.Parse()

	workflowDSL, err := parser.FromFile(dslConfig)
	if err != nil {
		log.Fatalln("Failed to parse dsl config", err)
	}

	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}

	workflowOptions := client.StartWorkflowOptions{
		ID:        "dsl_" + uuid.New(),
		TaskQueue: "dsl",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, app.HelloWorld, workflowDSL)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

}
