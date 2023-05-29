package file

import (
	"fmt"
	"jd_workout_golang/lib/helper"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type FileStore interface {
	Valdate() error
	Store() (string, error)
}

type GinFileStore struct {
	File     *multipart.FileHeader
	Path     string
	FileName string
}

func (fs GinFileStore) Valdate() error {

	ext := filepath.Ext(fs.File.Filename)

	if ext != ".jpg" && ext != ".png" {
		return fmt.Errorf("file type error")
	}

	return nil
}

func (fs GinFileStore) Store() (string, error) {

	gin := gin.Context{}

	storePath := fmt.Sprintf("./public/%s/%s%s", fs.Path, helper.RandString(10), filepath.Ext(fs.File.Filename))

	gin.SaveUploadedFile(fs.File, storePath)

	return strings.Replace(storePath, "/public", "", 1), nil
}
