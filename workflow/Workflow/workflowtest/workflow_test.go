package workflows_test

import (
	"testing"
	"time"

	activityPac "github.com/E-Furqan/Food-Delivery-System/Activity"
	datapipelineClient "github.com/E-Furqan/Food-Delivery-System/Client/DatapipelineClient"
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
)

func TestDataSyncWorkflow(t *testing.T) {

	datapipelineClientEnv := environmentVariable.ReadPipelineClientEnv()
	var DatapipelineClient datapipelineClient.DatapipelineClientInterface = datapipelineClient.NewClient(datapipelineClientEnv)
	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()

	// Create activity and workflow environment
	activity := activityPac.NewController(nil, nil, nil, nil, DatapipelineClient, DriveClient)
	workflow := workflows.NewController(activity, model.WorkFlowEnv{BATCH_SIZE: 20})
	suite := &testsuite.WorkflowTestSuite{}
	env := suite.NewTestWorkflowEnvironment()

	env.SetTestTimeout(60 * time.Minute)
	env.SetWorkflowRunTimeout(60 * time.Minute)

	// Register activities and workflow
	env.RegisterActivity(activity.FetchSourceConfiguration)
	env.RegisterActivity(activity.FetchDestinationConfiguration)
	env.RegisterActivity(activity.CreateSourceToken)
	env.RegisterActivity(activity.CreateDestinationToken)
	env.RegisterActivity(activity.ListFilesInFolder)
	env.RegisterActivity(activity.CopyBatchActivity)
	env.RegisterActivity(activity.AddLogs)
	env.RegisterWorkflow(workflow.DataSyncWorkflow)
	env.RegisterWorkflow(workflow.MoveDataWorkflow)

	// Execute the workflow
	pipeline := model.Pipeline{SourcesID: 1, DestinationsID: 1, PipelineID: 2}
	env.ExecuteWorkflow(workflow.DataSyncWorkflow, pipeline)

	// Assert no error
	assert.True(t, env.IsWorkflowCompleted())
	assert.NoError(t, env.GetWorkflowError())

}
