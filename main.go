package main

import (
	"encoding/json"
	"fmt"
	"internal/file_service"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

/*
1. Develop a Golang application serving as a File Upload Service.
2. Handle HTTP requests for file uploads and downloads.
3. Implement an endpoint for uploading files to the server.
4. Create an endpoint to retrieve a list of uploaded files.
5. Implement an endpoint to download a specific file.
6. Write comprehensive unit tests to cover the critical functionalities of your code.

	POST /files - Upload file
	GET /files/:name Download file by name
	GET /files List files
*/

const (
	MAX_SIZE_BYTES = 30 * 1024 * 1024
)

func onError(res http.ResponseWriter, err error, message string) {
	fmt.Println(message)
	res.WriteHeader(500)
	res.Write([]byte(err.Error()))
}

func uploadFile(res http.ResponseWriter, req *http.Request, fs *file_service.FileService) {
	req.Body = http.MaxBytesReader(res, req.Body, MAX_SIZE_BYTES)
	file, header, err := req.FormFile("file")
	fmt.Println(header)
	if err != nil {
		onError(res, err, "Error Retrieving the File")
		return
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		onError(res, err, "Error Reading File")
		return
	}
	err = fs.UploadFile(header.Filename, data)
	if err != nil {
		onError(res, err, "Error writing file")
		return
	}
	res.WriteHeader(http.StatusOK)
}

func downloadFile(res http.ResponseWriter, req *http.Request, fs *file_service.FileService, fileName string) {
	data, err := fs.DownloadFile(fileName)
	if err != nil {
		onError(res, err, "Error reading file")
		return
	}

	res.Header().Set("Content-Type", "application/octet-stream")
	res.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	res.WriteHeader(http.StatusOK)
	res.Write(data)
}

func listFiles(res http.ResponseWriter, req *http.Request, fs *file_service.FileService) {
	files := fs.ListFiles()
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	json.NewEncoder(res).Encode(files)
}

func routerFunc(res http.ResponseWriter, req *http.Request, fs *file_service.FileService) {
	if req.Method == "POST" {
		uploadFile(res, req, fs)
	} else if req.Method == "GET" {
		path := strings.Split(req.URL.Path, "/")
		if len(path) >= 2 && path[2] != "" {
			downloadFile(res, req, fs, path[2])
		} else {
			listFiles(res, req, fs)
		}
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func main() {
	fs := file_service.New(filepath.Join("./", "files"))
	http.HandleFunc("/files/", func(res http.ResponseWriter, req *http.Request) {
		routerFunc(res, req, fs)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
