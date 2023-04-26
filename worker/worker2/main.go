package main

import (
	"log"
	ac "temporal-tutorial/activity"
	mo "temporal-tutorial/model"
	wo "temporal-tutorial/workflow"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	logger.Info("Zap logger created")

	// Create the client object just once per process
	c, err := client.Dial(client.Options{
		HostPort:  client.DefaultHostPort,
		Namespace: client.DefaultNamespace,
	})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	// This worker hosts both Workflow and Activity functions
	w := worker.New(c, mo.MyTaskQueue2, worker.Options{})

	// register workflow
	w.RegisterWorkflowWithOptions(wo.MyWorkflow2, workflow.RegisterOptions{
		Name: "Test Workflow 2",
	})

	// register activity
	w.RegisterActivity(ac.MyActivity3)

	// Start listening to the Task Queue
	err = w.Run(worker.InterruptCh())
	if err != nil {
		logger.Fatal("unable to start Worker", zap.Error(err))
	}
}
