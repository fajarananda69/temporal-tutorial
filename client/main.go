package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	mo "temporal-tutorial/model"
	wo "temporal-tutorial/workflow"

	"github.com/segmentio/ksuid"
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
	loop := os.Args[1]
	value := os.Args[2]

	count, _ := strconv.Atoi(loop)
	var we client.WorkflowRun

	for i := 0; i < count; i++ {
		ID := ksuid.New().String()
		options := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("my-workflow-%s", ID),
			TaskQueue: mo.MyTaskQueue1,
		}
		we, err = c.ExecuteWorkflow(context.Background(), options, wo.MyWorkflow1, value)
		if err != nil {
			log.Fatalln("unable to complete Workflow", err)
		}

		// Get the results
		var greeting mo.Response
		err = we.Get(context.Background(), &greeting)
		if err != nil {
			log.Fatalln("unable to get Workflow result", err)
		}

		printResults(greeting, we.GetID(), we.GetRunID())
	}
}

func printResults(greeting mo.Response, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	jData, _ := json.Marshal(greeting)
	fmt.Println(string(jData))
}
