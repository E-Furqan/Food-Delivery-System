package utils

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	model "github.com/E-Furqan/Food-Delivery-System/Models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v2"
)

func GetEnv(key string, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}

func GenerateResponse(httpStatusCode int, c *gin.Context, title1 string, message1 string, title2 string, input interface{}) {

	errorMessage := strings.TrimPrefix(message1, "ERROR: ")
	response := gin.H{
		title1: errorMessage,
	}

	if title2 != "" && input != nil {
		response[title2] = input
	}

	c.JSON(httpStatusCode, response)
}

func CreateAuthObj(config model.Configuration) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  config.RedirectUris,
		Scopes:       []string{drive.DriveScope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthUri,
			TokenURL: config.TokenUri,
		},
	}
}

func OpenBrowser(url string) error {
	var cmd *exec.Cmd
	switch {
	case isWindows():
		cmd = exec.Command("start", url)
	case isMac():
		cmd = exec.Command("open", url)
	case isLinux():
		cmd = exec.Command("xdg-open", url)
	default:
		return fmt.Errorf("unsupported operating system")
	}
	return cmd.Start()
}

func isWindows() bool {
	return runtime.GOOS == "windows"
}

func isMac() bool {
	return runtime.GOOS == "darwin"
}

func isLinux() bool {
	return runtime.GOOS == "linux"
}
