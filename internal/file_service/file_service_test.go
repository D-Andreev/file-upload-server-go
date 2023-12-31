package file_service

import (
	"path/filepath"
	"testing"
	"testing/fstest"
)

func TestUploadFile(t *testing.T) {
	fileSys := fstest.MapFS{
		"cat.txt": {
			Data: []byte("cat"),
		},
	}
	fs := New(filepath.Join("./", "files"))
	err := fs.UploadFile("test")
}
