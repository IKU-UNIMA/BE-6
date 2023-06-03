package handler

import (
	"BE-6/src/config/database"
	"BE-6/src/config/env"
	"BE-6/src/config/storage"
	"BE-6/src/model"
	"BE-6/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func EditDokumenKerjasamaHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	dokumen, _ := c.FormFile("file")
	if dokumen == nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "file tidak boleh kosong"})
	}

	if err := util.CheckFileIsPDF(dokumen); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	tx := db.Begin()
	ctx := c.Request().Context()

	oldDokumen := ""
	if err := db.WithContext(ctx).Model(new(model.Kerjasama)).Select("dokumen").First(&oldDokumen, "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "kerjasama tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	dDokumen, err := storage.CreateFile(dokumen, env.GetDokumenFolderId())
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Table("kerjasama").
		Where("id", id).Update("dokumen", util.CreateFileUrl(dDokumen.Id)).Error; err != nil {
		tx.Rollback()
		storage.DeleteFile(dDokumen.Id)
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := storage.DeleteFile(util.GetFileIdFromUrl(oldDokumen)); err != nil {
		tx.Rollback()
		storage.DeleteFile(dDokumen.Id)
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.FailedResponse(http.StatusOK, nil)
}
