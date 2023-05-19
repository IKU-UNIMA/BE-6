package util

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func CheckFileIsExcel(file *multipart.FileHeader) error {
	fileType := file.Header["Content-Type"][0]
	if fileType == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
		return nil
	}

	return FailedResponse(http.StatusBadRequest, map[string]string{"message": "unsupported file type for " + file.Filename})
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
