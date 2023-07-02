package handler

import (
	"BE-6/src/api/request"
	"BE-6/src/api/response"
	"BE-6/src/config/database"
	"BE-6/src/model"
	"BE-6/src/util"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetAllProdiHandler(c echo.Context) error {
	idFakultas := c.QueryParam("fakultas")

	db := database.DB
	ctx := c.Request().Context()
	result := []response.Prodi{}
	condition := ""

	if idFakultas != "" && idFakultas != "0" {
		condition = "id_fakultas = " + idFakultas
	}

	if err := db.WithContext(ctx).Preload("Fakultas").Where(condition).Order("id_fakultas").Find(&result).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func GetProdiByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	result := &response.Prodi{}

	if err := db.WithContext(ctx).Preload("Fakultas").First(result, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func InsertProdiHandler(c echo.Context) error {
	request := &request.Prodi{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()

	data := request.MapRequest()

	if err := db.WithContext(ctx).Create(data).Error; err != nil {
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "fakultas sudah ada"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusCreated, data.ID)
}

func EditProdiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	request := &request.Prodi{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Prodi), id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Where("id", id).Updates(request.MapRequest()).Error; err != nil {
		if err != nil {
			if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "prodi sudah ada"})
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteProdiHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.Prodi), id)
	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
