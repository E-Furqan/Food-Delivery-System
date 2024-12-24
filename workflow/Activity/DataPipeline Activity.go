package activity

import (
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
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

func (act *Activity) CreateSourceToken(source model.Config) (string, error) {

	sourceToken, err := act.DriveClient.CreateToken(source)
	if err != nil {
		return "", err
	}
	return sourceToken, nil
}

func (act *Activity) CreateDestinationToken(destination model.Config) (string, error) {

	destinationToken, err := act.DriveClient.CreateToken(destination)
	if err != nil {
		return "", err
	}

	return destinationToken, nil
}

func (act *Activity) MoveDataFromSourceToDestination(sourceToken string, destinationToken string,
	sourceFolderUrl string, destinationFolderUrl string, sourceConfig model.Config) (model.FileCounter, error) {

	var counter model.FileCounter

	sourceClient, err := act.DriveClient.CreateConnection(sourceToken, sourceConfig)
	if err != nil {
		return model.FileCounter{}, fmt.Errorf("invalid source client: %w", err)
	}

	sourceFolderID, err := utils.ExtractFolderID(sourceFolderUrl)
	if err != nil {
		return model.FileCounter{}, fmt.Errorf("invalid source folder URL: %w", err)
	}

	destinationFolderID, err := utils.ExtractFolderID(destinationFolderUrl)
	if err != nil {
		return model.FileCounter{}, fmt.Errorf("invalid destination folder URL: %w", err)
	}

	fileList, err := utils.ListFilesInFolder(&sourceClient, sourceFolderID)
	if err != nil {
		return model.FileCounter{}, fmt.Errorf("failed to list files in source folder: %w", err)
	}

	for _, file := range fileList {
		_, err := sourceClient.Files.Update(file.Id, nil).
			AddParents(destinationFolderID).
			RemoveParents(sourceFolderID).
			Do()
		if err != nil {
			counter.FailedCounter += 1
		}
		counter.NoOfFiles += 1
	}
	log.Print("failed counter; ", counter.FailedCounter)
	log.Print("files; ", counter.NoOfFiles)

	return counter, nil
}

func (act *Activity) AddLogs(counter model.FileCounter, PipelinesID int) error {
	var log model.Log
	FilesMovedSuccessfully := counter.NoOfFiles - counter.FailedCounter

	if counter.FailedCounter != 0 {
		log.LogMessage = fmt.Sprintf("the data sync failed to move %v files but successfully moved %v files", counter.FailedCounter, FilesMovedSuccessfully)
	} else {
		log.LogMessage = fmt.Sprintf("the data sync successfully moved %v files", FilesMovedSuccessfully)
	}
	log.PipelinesID = PipelinesID

	err := act.DatapipelineClient.AddLogs(log)
	if err != nil {
		return err
	}

	return nil
}
