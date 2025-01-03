package activityPac

import (
	"context"
	"fmt"
	"log"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"go.temporal.io/sdk/activity"
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

func (act *Activity) ListFilesInFolder(sourceToken string, sourceConfig model.Config, folderID string) ([]*drive.File, error) {
	sourceClient, err := act.DriveClient.CreateConnection(sourceToken, sourceConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid source client: %w", err)
	}

	log.Println("folderID", folderID)
	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	fileList, err := sourceClient.Files.List().Q(query).Do()
	if err != nil {
		log.Println("Error listing files:", err)
		return nil, err
	}
	log.Printf("counting files in folder: %s", folderID)

	return fileList.Files, nil
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

func (act *Activity) CopyBatchActivity(ctx context.Context, sourceToken string, destinationToken string, sourceConfig model.Config,
	destinationConfig model.Config, sourceFolderID string, destinationFolderID string, fileList []*drive.File, counter model.FileCounter, startIndex int, endIndex int) (model.FileCounter, error) {

	sourceClient, err := act.DriveClient.CreateConnection(sourceToken, model.Config{})
	if err != nil {
		return counter, fmt.Errorf("invalid source client: %w", err)
	}

	if startIndex < 0 || endIndex > len(fileList) || startIndex > endIndex {
		return counter, fmt.Errorf("invalid startIndex or endIndex: startIndex=%d, endIndex=%d, fileList length=%d", startIndex, endIndex, len(fileList))
	}

	// Iterate through the batch
	for i := startIndex; i < endIndex; i++ {
		file := fileList[i]

		newFile := &drive.File{
			Name:    file.Name,
			Parents: []string{destinationFolderID},
		}

		_, err := sourceClient.Files.Copy(file.Id, newFile).Do()
		activity.RecordHeartbeat(ctx, i)
		if err != nil {
			log.Printf("Failed to copy file name %s: %v", file.Name, err)
			counter.FailedCounter++
		} else {
			log.Printf("Successfully copied file name %s", file.Name)
			counter.NoOfFiles++
		}
	}

	return counter, nil
}
