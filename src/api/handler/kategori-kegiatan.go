package handler

import (
	"BE-6/src/api/request"
	"BE-6/src/api/response"
	"BE-6/src/config/database"
	"BE-6/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetAllKategoriKegiatan(c echo.Context) error {
	db := database.DB
	ctx := c.Request().Context()
	data := []response.JenisKategoriKegiatan{}

	if err := db.WithContext(ctx).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func InsertKategoriKegiatan(c echo.Context) error {
	req := &request.JenisKategoriKegiatan{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.DB
	ctx := c.Request().Context()
	data := req.MapRequest()

	if err := db.WithContext(ctx).Create(data); err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data.ID)
}
