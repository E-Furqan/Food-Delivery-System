package activityPac

import (
	"context"
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	utils "github.com/E-Furqan/Food-Delivery-System/Utils"
	"go.temporal.io/sdk/activity"
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

func (act *Activity) CountFilesInFolder(sourceToken string, sourceConfig model.Config, folderID string) (int, error) {
	sourceClient, err := act.DriveClient.CreateConnection(sourceToken, sourceConfig)
	if err != nil {
		return 0, fmt.Errorf("invalid source client: %w", err)
	}

	log.Println("folderID", folderID)
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	fileList, err := sourceClient.Files.List().Q(query).Do()
	if err != nil {
		log.Println("Error listing files:", err)
		return 0, err
	}
	log.Printf("counting files in folder: %s", folderID)

	return len(fileList.Files), nil
}

func (act *Activity) MoveDataFromSourceToDestination(ctx context.Context, sourceToken string, destinationToken string,
	sourceFolderUrl string, destinationFolderUrl string, sourceConfig model.Config, batchSize int) (model.FileCounter, error) {

	var counter model.FileCounter

	sourceClient, err := act.DriveClient.CreateConnection(sourceToken, sourceConfig)
	if err != nil {
		return model.FileCounter{}, fmt.Errorf("invalid source client: %w", err)
	}

	activity.RecordHeartbeat(ctx, nil)
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

	totalFiles := len(fileList)

	for i := 0; i < totalFiles; i += batchSize {
		end := i + batchSize
		if end > totalFiles {
			end = totalFiles
		}

		batch := fileList[i:end]
		for _, file := range batch {

			_, err := sourceClient.Files.Update(file.Id, nil).
				AddParents(destinationFolderID).
				RemoveParents(sourceFolderID).
				Do()

			if err != nil {
				counter.FailedCounter += 1
				log.Printf("Failed to move file name %s: %v", file.Name, err)
			}
			counter.NoOfFiles += 1

		}

		log.Printf("Processed batch %d - %d", i+1, end)

		activity.RecordHeartbeat(ctx, nil)
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

func (act *Activity) MoveBatchActivity(ctx context.Context, sourceToken string, destinationToken string, sourceConfig model.Config,
	destinationConfig model.Config, sourceFolderID string, destinationFolderID string, startIndex int, endIndex int) error {

	sourceClient, err := act.DriveClient.CreateConnection(sourceToken, model.Config{})
	if err != nil {
		return fmt.Errorf("invalid source client: %w", err)
	}

	fileList, err := utils.ListFilesInFolder(&sourceClient, sourceFolderID)
	if err != nil {
		return err
	}

	batch := fileList[startIndex:endIndex]

	for _, file := range batch {
		_, err := sourceClient.Files.Update(file.Id, nil).
			AddParents(destinationFolderID).
			RemoveParents(sourceFolderID).
			Do()

		if err != nil {
			log.Printf("Failed to move file ID %s: %v", file.Id, err)
		}
	}

	return nil
}
