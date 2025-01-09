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

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (s *UnitTestSuite) Test_DataSyncWorkflow_Success() {

	datapipelineClientEnv := environmentVariable.ReadPipelineClientEnv()
	var DatapipelineClient datapipelineClient.DatapipelineClientInterface = datapipelineClient.NewClient(datapipelineClientEnv)
	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()

	activity := activityPac.NewController(nil, nil, nil, nil, DatapipelineClient, DriveClient)
	workflow := workflows.NewController(activity, model.WorkFlowEnv{BATCH_SIZE: 20})

	var source model.Source
	source.SourcesID = 1

	s.env.OnActivity(activity.FetchSourceConfiguration, source).Return(model.Config{
		ClientID:       "ClientID",
		ClientSecret:   "ClientSecret",
		TokenURI:       "TokenURI",
		RefreshToken:   "RefreshToken",
		FolderURL:      "https://drive.google.com/drive/u/0/folders/140UlUXL5QXfG4sSxmpBAE_c77gqeVPgD",
		SourcesID:      1,
		DestinationsID: 2,
	}, nil)

	var destination model.Destination
	destination.DestinationsID = 1

	s.env.OnActivity(activity.FetchDestinationConfiguration, destination).Return(model.Config{
		ClientID:       "ClientID",
		ClientSecret:   "ClientSecret",
		TokenURI:       "TokenURI",
		RefreshToken:   "RefreshToken",
		FolderURL:      "https://drive.google.com/drive/u/0/folders/1HhBcRTdBWgcR5lPlslJucpyuaGODDTXc",
		SourcesID:      1,
		DestinationsID: 2,
	}, nil)

	s.env.OnActivity(activity.CreateSourceToken, mock.Anything).Return("sourceToken", nil)
	s.env.OnActivity(activity.CreateDestinationToken, mock.Anything).Return("destinationToken", nil)
	s.env.OnActivity(activity.AddLogs, model.FileCounter{NoOfFiles: 5, FailedCounter: 1}, 2).Return(nil)

	s.env.OnWorkflow(workflow.MoveDataWorkflow, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(model.FileCounter{NoOfFiles: 5, FailedCounter: 1}, nil)

	s.env.ExecuteWorkflow(workflow.DataSyncWorkflow, model.Pipeline{PipelineID: 2, SourcesID: 1, DestinationsID: 1})

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	s.NoError(err)
}

func (s *UnitTestSuite) Test_DataSyncWorkflow_FetchSourceConfiguration_Failure() {
	datapipelineClientEnv := environmentVariable.ReadPipelineClientEnv()
	var DatapipelineClient datapipelineClient.DatapipelineClientInterface = datapipelineClient.NewClient(datapipelineClientEnv)
	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()

	activity := activityPac.NewController(nil, nil, nil, nil, DatapipelineClient, DriveClient)
	workflow := workflows.NewController(activity, model.WorkFlowEnv{BATCH_SIZE: 20})
	s.env.OnActivity(activity.FetchSourceConfiguration, mock.Anything, mock.Anything).Return(model.Config{}, errors.New("FetchSourceConfiguration failed"))

	pipeline := model.Pipeline{SourcesID: 1, DestinationsID: 1, PipelineID: 2}
	s.env.ExecuteWorkflow(workflow.DataSyncWorkflow, pipeline)

	assert.True(s.T(), s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "FetchSourceConfiguration failed")

}

func (s *UnitTestSuite) Test_MoveDataWorkflow_Success() {

	datapipelineClientEnv := environmentVariable.ReadPipelineClientEnv()
	var DatapipelineClient datapipelineClient.DatapipelineClientInterface = datapipelineClient.NewClient(datapipelineClientEnv)
	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()

	activity := activityPac.NewController(nil, nil, nil, nil, DatapipelineClient, DriveClient)
	workflow1 := workflows.NewController(activity, model.WorkFlowEnv{BATCH_SIZE: 20})

	s.env.OnActivity(activity.ListFilesInFolder, mock.Anything, mock.Anything, mock.Anything, 2, 0).
		Return([]*drive.File{
			{Name: "pipelineModel.go"},
			{Name: "LogModel.go"},
		}, nil)

	s.env.OnActivity(activity.CopyBatchActivity, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		[]*drive.File{
			{Name: "pipelineModel.go"},
			{Name: "LogModel.go"},
		}, model.FileCounter{}).
		Return(model.FileCounter{NoOfFiles: 5, FailedCounter: 1}, nil)

	sourceConfig := model.Config{
		ClientID:       "ClientID",
		ClientSecret:   "ClientSecret",
		TokenURI:       "TokenURI",
		RefreshToken:   "RefreshToken",
		FolderURL:      "https://drive.google.com/drive/u/0/folders/140UlUXL5QXfG4sSxmpBAE_c77gqeVPgD",
		SourcesID:      1,
		DestinationsID: 2,
	}

	destConfig := model.Config{
		ClientID:       "ClientID",
		ClientSecret:   "ClientSecret",
		TokenURI:       "TokenURI",
		RefreshToken:   "RefreshToken",
		FolderURL:      "https://drive.google.com/drive/u/0/folders/1HhBcRTdBWgcR5lPlslJucpyuaGODDTXc",
		SourcesID:      1,
		DestinationsID: 2,
	}

	s.env.ExecuteWorkflow(workflow1.MoveDataWorkflow, "source token", "Destination Token", "140UlUXL5QXfG4sSxmpBAE_c77gqeVPgD", "1HhBcRTdBWgcR5lPlslJucpyuaGODDTXc",
		sourceConfig, destConfig, 2, model.FileCounter{NoOfFiles: 0, FailedCounter: 0}, 0)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "continue as new")
}

func (s *UnitTestSuite) Test_MoveDataWorkflow_Failure() {

	datapipelineClientEnv := environmentVariable.ReadPipelineClientEnv()
	var DatapipelineClient datapipelineClient.DatapipelineClientInterface = datapipelineClient.NewClient(datapipelineClientEnv)
	var DriveClient driveClient.DriveClientInterface = driveClient.NewClient()

	activity := activityPac.NewController(nil, nil, nil, nil, DatapipelineClient, DriveClient)
	workflow1 := workflows.NewController(activity, model.WorkFlowEnv{BATCH_SIZE: 20})

	s.env.OnActivity(activity.ListFilesInFolder, mock.Anything, mock.Anything, mock.Anything, 2, 0).
		Return([]*drive.File{}, errors.New("Error While Listing Files"))

	sourceConfig := model.Config{
		ClientID:       "ClientID",
		ClientSecret:   "ClientSecret",
		TokenURI:       "TokenURI",
		RefreshToken:   "RefreshToken",
		FolderURL:      "https://drive.google.com/drive/u/0/folders/140UlUXL5QXfG4sSxmpBAE_c77gqeVPgD",
		SourcesID:      1,
		DestinationsID: 2,
	}

	destConfig := model.Config{
		ClientID:       "ClientID",
		ClientSecret:   "ClientSecret",
		TokenURI:       "TokenURI",
		RefreshToken:   "RefreshToken",
		FolderURL:      "https://drive.google.com/drive/u/0/folders/1HhBcRTdBWgcR5lPlslJucpyuaGODDTXc",
		SourcesID:      1,
		DestinationsID: 2,
	}

	s.env.ExecuteWorkflow(workflow1.MoveDataWorkflow, "source token", "Destination Token", "140UlUXL5QXfG4sSxmpBAE_c77gqeVPgD", "1HhBcRTdBWgcR5lPlslJucpyuaGODDTXc",
		sourceConfig, destConfig, 2, model.FileCounter{NoOfFiles: 0, FailedCounter: 0}, 0)

	s.True(s.env.IsWorkflowCompleted())
	err := s.env.GetWorkflowError()
	assert.Error(s.T(), err)
	assert.Contains(s.T(), err.Error(), "Error While Listing Files")
}
