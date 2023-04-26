package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	mo "temporal-tutorial/model"
	wo "temporal-tutorial/workflow"

	"go.temporal.io/sdk/client"
)

func main() {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{
		HostPort:  client.DefaultHostPort,
		Namespace: client.DefaultNamespace,
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// Start the Workflow
	service := os.Args[1]
	value := os.Args[2:]

	var we client.WorkflowRun
	if service == "3" {
		var val1, val2 string
		if len(value) > 1 {
			val1 = value[0]
			val2 = value[1]
		}
		options := client.StartWorkflowOptions{
			ID:        "my-workflow-2",
			TaskQueue: mo.MyTaskQueue2,
		}
		we, err = c.ExecuteWorkflow(context.Background(), options, wo.MyWorkflow2, service, val1, val2)
		if err != nil {
			log.Fatalln("unable to complete Workflow", err)
		}
	} else {
		options := client.StartWorkflowOptions{
			ID:        "my-workflow-1",
			TaskQueue: mo.MyTaskQueue1,
		}
		we, err = c.ExecuteWorkflow(context.Background(), options, wo.MyWorkflow1, service, value[0])
		if err != nil {
			log.Fatalln("unable to complete Workflow", err)
		}
	}

	// Get the results
	var greeting mo.Response
	err = we.Get(context.Background(), &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(greeting, we.GetID(), we.GetRunID())
}

func printResults(greeting mo.Response, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	jData, _ := json.Marshal(greeting)
	fmt.Println(string(jData))
}
