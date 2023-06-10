package file

import (
	"fmt"
	"jd_workout_golang/lib/helper"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileStore interface {
	Validate() error
	Store() (*string, error)
	GetPath() *string
}

type GinFileStore struct {
	File     *multipart.FileHeader
	Path     *string
	FileName string
}

func (fs GinFileStore) Validate() error {

	ext := strings.ToLower(filepath.Ext(fs.File.Filename))

	if ext != ".jpg" && ext != ".png" {
		return fmt.Errorf("file type error")
	}

	return nil
}

func (fs GinFileStore) Store() (*string, error) {

	gin := gin.Context{}

	storePath := fmt.Sprintf("./public/%s/%s%s", *fs.Path, helper.RandString(10), filepath.Ext(fs.File.Filename))

	gin.SaveUploadedFile(fs.File, storePath)

	path := strings.Replace(storePath, "./public", "", 1)

	return &path, nil
}

func (fs GinFileStore) GetPath() *string {
	url := os.Getenv("APP_URL")

	if fs.Path == nil {
		return nil
	}

	path := fmt.Sprintf("%s%s", url, *fs.Path)

	return &path
}
