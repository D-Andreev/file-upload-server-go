package file_service

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type FileService struct {
	path string
}

func New(path string) *FileService {
	return &FileService{
		path: path,
	}
}

func (fs *FileService) getPath(fileName string) string {
	return filepath.Join(fs.path, fileName)
}

func (fs *FileService) UploadFile(fileName string, content []byte) error {
	path := fs.getPath(fileName)
	if _, err := os.Stat(path); err == nil {
		return fmt.Errorf("%s already exists", fileName)
	}
	err := os.WriteFile(path, content, 0644)
	return err
}

func (fs *FileService) DownloadFile(fileName string) ([]byte, error) {
	path := fs.getPath(fileName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fs *FileService) ListFiles() []string {
	files, err := ioutil.ReadDir(fs.path)
	if err != nil {
		log.Fatal(err)
	}

	var result []string
	for _, file := range files {
		if !file.IsDir() {
			result = append(result, file.Name())
		}
	}

	return result
}
