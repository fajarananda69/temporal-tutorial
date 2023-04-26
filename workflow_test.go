package test

import (
	ac "temporal-tutorial/activity"
	mo "temporal-tutorial/model"
	wo "temporal-tutorial/workflow"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	// Set up the test suite and testing execution environment
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// Mock activity implementation
	env.OnActivity(ac.MyActivity1, mock.Anything, "World").Return("Hello World!", nil)

	env.ExecuteWorkflow(wo.MyWorkflow1, "1", "World")
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())

	var greeting mo.Response
	require.NoError(t, env.GetWorkflowResult(&greeting))
	require.Equal(t, mo.Response{Status: 200, Data: "Hello World!"}, greeting)
}
