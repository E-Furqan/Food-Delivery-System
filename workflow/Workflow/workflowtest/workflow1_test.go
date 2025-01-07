package workflows_test

import (
	"errors"
	"testing"

	activityPac "github.com/E-Furqan/Food-Delivery-System/Activity"
	datapipelineClient "github.com/E-Furqan/Food-Delivery-System/Client/DatapipelineClient"
	driveClient "github.com/E-Furqan/Food-Delivery-System/Client/DriveClient"
	environmentVariable "github.com/E-Furqan/Food-Delivery-System/EnviormentVariable"
	model "github.com/E-Furqan/Food-Delivery-System/Models"
	workflows "github.com/E-Furqan/Food-Delivery-System/Workflow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.temporal.io/sdk/testsuite"
	"google.golang.org/api/drive/v3"
)

type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite
	env *testsuite.TestWorkflowEnvironment
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func TestDataSyncWorkflowTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_DataSyncWorkflow_Success() {
	// Initialize the necessary environment and clients
	datapipelineClientEnv := environmentVariable.ReadPipelineClientEnv()
	var DatapipelineClient datapipelineClient.DatapipelineClientInterface = datapipelineClient.NewClient(datapipelineClientEnv)
	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()

	// Create activity and workflow environment
	activity := activityPac.NewController(nil, nil, nil, nil, DatapipelineClient, DriveClient)
	workflow := workflows.NewController(activity, model.WorkFlowEnv{BATCH_SIZE: 20})

	// Mock the activities with correct return types based on your models
	s.env.OnActivity(activity.FetchSourceConfiguration, mock.Anything, mock.Anything).Return(model.Config{
		ClientID:       "clientID",
		ClientSecret:   "clientSecret",
		TokenURI:       "tokenURI",
		RefreshToken:   "refreshToken",
		FolderURL:      "folderURL",
		SourcesID:      1,
		DestinationsID: 1,
	}, nil)

	s.env.OnActivity(activity.FetchDestinationConfiguration, mock.Anything, mock.Anything).Return(model.Config{
		ClientID:       "destClientID",
		ClientSecret:   "destClientSecret",
		TokenURI:       "destTokenURI",
		RefreshToken:   "destRefreshToken",
		FolderURL:      "destFolderURL",
		SourcesID:      1,
		DestinationsID: 2,
	}, nil)

	s.env.OnActivity(activity.CreateSourceToken, mock.Anything, mock.Anything).Return("sourceToken", nil)
	s.env.OnActivity(activity.CreateDestinationToken, mock.Anything, mock.Anything).Return("destinationToken", nil)

	s.env.OnActivity(activity.ListFilesInFolder, mock.Anything, mock.Anything).Return([]*drive.File{
		&drive.File{Name: "file1"},
		&drive.File{Name: "file2"},
	}, nil)

	s.env.OnActivity(activity.CopyBatchActivity, mock.Anything, mock.Anything).Return(nil)
	s.env.OnActivity(activity.AddLogs, mock.Anything, mock.Anything).Return(nil)

	// Execute the workflow
	s.env.ExecuteWorkflow(workflow.DataSyncWorkflow, "test_success")

	// Verify if the workflow completed successfully
	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.NoError(err)
}

func (s *UnitTestSuite) Test_DataSyncWorkflow_FetchSourceConfiguration_Failure() {
	// Initialize clients
	datapipelineClientEnv := environmentVariable.ReadPipelineClientEnv()
	var DatapipelineClient datapipelineClient.DatapipelineClientInterface = datapipelineClient.NewClient(datapipelineClientEnv)
	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()

	// Create activity and workflow environment
	activity := activityPac.NewController(nil, nil, nil, nil, DatapipelineClient, DriveClient)
	workflow := workflows.NewController(activity, model.WorkFlowEnv{BATCH_SIZE: 20})
	s.env.OnActivity(activity.FetchSourceConfiguration, mock.Anything, mock.Anything).Return("", errors.New("FetchSourceConfiguration failed"))

	// Execute the workflow
	pipeline := model.Pipeline{SourcesID: 1, DestinationsID: 1, PipelineID: 2}
	s.env.ExecuteWorkflow(workflow.DataSyncWorkflow, pipeline)

	// Assert error
	assert.True(s.T(), s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assert.Error(s.T(), err)
	assert.Equal(s.T(), "FetchSourceConfiguration failed", err.Error())
}
