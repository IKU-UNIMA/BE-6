package util

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func CheckFileType(file *multipart.FileHeader) error {
	fileType := file.Header["Content-Type"][0]
	switch fileType {
	case "application/pdf":
		return nil
	case "image/png":
		return nil
	case "image/jpeg":
		return nil
	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return nil
	case "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":
		return nil
	case "text/plain":
		return nil
	}

	return errors.New(fmt.Sprint("unsupported file type for ", file.Filename))
}

func CheckFileIsExcel(file *multipart.FileHeader) error {
	fileType := file.Header["Content-Type"][0]
	if fileType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		return nil
	}

	return FailedResponse(http.StatusBadRequest, map[string]string{"message": "unsupported file type for " + file.Filename})
}

func CreateFileUrl(fileId string) string {
	return "https://drive.google.com/file/d/" + fileId + "/view"
}

func CheckFileIsPDF(file *multipart.FileHeader) error {
	fileType := file.Header["Content-Type"][0]
	if fileType == "application/pdf" {
		return nil
	}

	return errors.New(fmt.Sprint("unsupported file type for ", file.Filename))
}

func WriteFile(file *multipart.FileHeader) error {
	if err := CheckFileIsExcel(file); err != nil {
		return err
	}

	src, err := file.Open()
	if err != nil {
		return FailedResponse(http.StatusInternalServerError, nil)
	}
	defer src.Close()

	dst, err := os.Create(file.Filename)
	if err != nil {
		return FailedResponse(http.StatusInternalServerError, nil)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return FailedResponse(http.StatusInternalServerError, nil)
	}

	return nil
}

func GetFileIdFromUrl(url string) string {
	arr := strings.Split(url, "/")
	return arr[len(arr)-2]
}
