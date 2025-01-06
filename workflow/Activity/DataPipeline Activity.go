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

func (act *Activity) ListFilesInFolder(sourceToken string, sourceConfig model.Config, folderID string,
	bathSize int, startIndex int) ([]*drive.File, error) {
	sourceClient, err := act.DriveClient.CreateConnection(sourceToken, sourceConfig)
	if err != nil {
		return nil, fmt.Errorf("invalid source client: %w", err)
	}

	query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
	nextPageToken := ""
	var FileList []*drive.File

	for {
		fileList, err := sourceClient.Files.List().
			Q(query).
			Fields("nextPageToken, files(id, name, mimeType)").
			PageToken(nextPageToken).
			Do()

		if err != nil {
			return nil, fmt.Errorf("unable to retrieve files: %w", err)
		}

		for _, item := range fileList.Files {
			if item.MimeType == "application/vnd.google-apps.folder" {

				subFolderResult, err := act.ListFilesInFolder(sourceToken, sourceConfig, item.Id, bathSize, 0)
				if err != nil {
					return nil, fmt.Errorf("failed to list files in subfolder %s: %w", item.Name, err)
				}

				FileList = append(FileList, subFolderResult...)
			} else {
				FileList = append(FileList, item)
			}
		}
		if fileList.NextPageToken == "" {
			break
		}

		nextPageToken = fileList.NextPageToken
	}
	log.Print("files fetched: ", len(FileList))

	return FileList, nil
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
	destinationConfig model.Config, sourceFolderID string, destinationFolderID string, fileList []*drive.File, counter model.FileCounter) (model.FileCounter, error) {

	// Create a connection to the source client
	sourceClient, err := act.DriveClient.CreateConnection(sourceToken, model.Config{})
	if err != nil {
		return counter, fmt.Errorf("invalid source client: %w", err)
	}

	totalFiles := len(fileList)
	log.Printf("Copying %d files", totalFiles)

	var lastProcessedIndex int = 0
	hb := activity.GetHeartbeatDetails(ctx, &lastProcessedIndex)
	if hb == nil {
		log.Printf("error while getting heartbeat  %v", hb)
	}
	log.Printf("Last processed file index: %d", lastProcessedIndex)

	for i := lastProcessedIndex; i < totalFiles; i++ {
		file := fileList[i]
		if file.MimeType == "application/vnd.google-apps.folder" {
			log.Printf("Skipping folder name %s", file.Name)
			continue
		}

		newFile := &drive.File{
			Name:    file.Name,
			Parents: []string{destinationFolderID},
		}

		_, err := sourceClient.Files.Copy(file.Id, newFile).Do()
		if err != nil {
			log.Printf("Failed to copy file name %s: %v", file.Name, err)
			counter.FailedCounter++

			activity.RecordHeartbeat(ctx, int(i))
			return counter, fmt.Errorf("failed to copy file name %s: %v", file.Name, err)
		} else {
			log.Printf("Successfully copied file name %s", file.Name)
			counter.NoOfFiles++
		}

		activity.RecordHeartbeat(ctx, int(i))
	}

	return counter, nil
}
