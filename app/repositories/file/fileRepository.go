package file

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jd_workout_golang/lib/helper"
	"mime/multipart"
	"path/filepath"
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

	path := fmt.Sprintf("./public/%s/%s%s", fs.Path, helper.RandString(10), filepath.Ext(fs.File.Filename))

	gin.SaveUploadedFile(fs.File, path)

	return path, nil
}
