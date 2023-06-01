package env

import "os"

func GetDokumenFolderId() string {
	return os.Getenv("GOOGLE_DRIVE_DOKUMEN_FOLDER_ID")
}
