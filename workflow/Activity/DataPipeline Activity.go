package activity

import (
	"fmt"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"google.golang.org/api/drive/v3"
)

func (act *Activity) FetchSourceConfiguration(source model.Source) (model.Config, error) {

	configuration, err := act.DatapipelineClient.FetchSourceConfiguration(source)
	if err != nil {
		return model.Config{}, err
	}

	return configuration, nil
}

func (act *Activity) FetchDestinationConfiguration(destination model.Destination) (model.Config, error) {

	configuration, err := act.DatapipelineClient.FetchDestinationConfiguration(destination)
	if err != nil {
		return model.Config{}, err
	}

	return configuration, nil
}

func (act *Activity) CreateSourceConnection(source model.Config) (*drive.Service, error) {

	SourceConnection, err := act.DriveClient.CreateConnection(source)
	if err != nil {
		return &drive.Service{}, err
	}

	return SourceConnection, nil
}

func (act *Activity) CreateDestinationConnection(destination model.Config) (*drive.Service, error) {

	destinationConnection, err := act.DriveClient.CreateConnection(destination)
	if err != nil {
		return &drive.Service{}, err
	}

	return destinationConnection, nil
}

func (act *Activity) MoveDataFromSourceToDestination(sourceClient *drive.Service, destinationClient *drive.Service,
	sourceFolderUrl string, destinationFolderUrl string) (int, error) {

	var failedCounter int = 0

	sourceFolderID, err := utils.ExtractFolderID(sourceFolderUrl)
	if err != nil {
		return failedCounter, fmt.Errorf("invalid source folder URL: %w", err)
	}
	destinationFolderID, err := utils.ExtractFolderID(destinationFolderUrl)
	if err != nil {
		return failedCounter, fmt.Errorf("invalid destination folder URL: %w", err)
	}

	fileList, err := utils.ListFilesInFolder(sourceClient, sourceFolderID)
	if err != nil {
		return failedCounter, fmt.Errorf("failed to list files in source folder: %w", err)
	}

	for _, file := range fileList {
		_, err := sourceClient.Files.Update(file.Id, nil).
			AddParents(destinationFolderID).
			RemoveParents(sourceFolderID).
			Do()
		if err != nil {
			failedCounter += 1
		}
	}

	return failedCounter, nil
}
