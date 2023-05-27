package file

import (
	"fmt"
	"mime/multipart"
	"path/filepath"

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

	path := fmt.Sprintf("./public/%s/", fs.Path) + fs.FileName

	gin.SaveUploadedFile(fs.File, path)

	return path, nil
}
