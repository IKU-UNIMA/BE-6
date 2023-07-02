package handler

import (
	"BE-6/src/api/request"
	"BE-6/src/api/response"
	"BE-6/src/config/database"
	"BE-6/src/model"
	"BE-6/src/util"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const getRektorQuery = "SELECT rektor.id, nama, nip, akun.email, bagian FROM rektor JOIN akun where rektor.id = rektor.id"

func GetAllRektorHandler(c echo.Context) error {
	db := database.DB
	ctx := c.Request().Context()
	result := []response.Rektor{}

	if err := db.WithContext(ctx).Raw(getRektorQuery).Find(&result).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func GetRektorByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	result := &response.Rektor{}

	condition := getRektorQuery + fmt.Sprintf(" AND rektor.id = %d", id)
	if err := db.WithContext(ctx).Raw(condition).First(result).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

func InsertRektorHandler(c echo.Context) error {
	request := &request.Rektor{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	tx := db.Begin()
	ctx := c.Request().Context()
	akun := &model.Akun{}
	akun.Email = request.Email
	akun.Role = string(util.REKTOR)
	password := util.GeneratePassword()
	akun.Password = util.HashPassword(password)

	if err := tx.WithContext(ctx).Create(akun).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "email sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	rektor := request.MapRequest()
	rektor.ID = akun.ID

	if err := tx.WithContext(ctx).Create(rektor).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NIP sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusCreated, map[string]string{"password": password})
}

func EditRektorHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	request := &request.Rektor{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	tx := db.Begin()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Rektor), id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Table("akun").Where("id", id).Update("email", request.Email).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "email sudah digunakan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	result := request.MapRequest()
	if err := tx.WithContext(ctx).Where("id", id).Omit("password").Updates(result).Error; err != nil {
		if err != nil {
			tx.Rollback()
			if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "NIP sudah digunakan"})
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteRektorHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.Akun), id)
	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
