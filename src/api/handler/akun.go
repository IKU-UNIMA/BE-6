package handler

import (
	"BE-6/src/api/request"
	"BE-6/src/api/response"
	"BE-6/src/config/database"
	"BE-6/src/model"
	"BE-6/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func LoginHandler(c echo.Context) error {
	request := &request.Login{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	data := &model.Akun{}

	if err := db.WithContext(ctx).First(data, "email", request.Email).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusUnauthorized, map[string]string{"message": "email atau password salah"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if !util.ValidateHash(request.Password, data.Password) {
		return util.FailedResponse(http.StatusUnauthorized, map[string]string{"message": "email atau password salah"})
	}

	var bagian string
	if data.Role == string(util.ADMIN) {
		if err := db.WithContext(ctx).Table("admin").Select("bagian").Where("id", data.ID).Scan(&bagian).Error; err != nil {
			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	var nama string
	if err := db.WithContext(ctx).Table(data.Role).Select("nama").Where("id", data.ID).Scan(&nama).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	token := util.GenerateJWT(data.ID, nama, data.Role, bagian)

	return util.SuccessResponse(c, http.StatusOK, response.Login{Token: token})
}

func ChangePasswordHandler(c echo.Context) error {
	request := &request.ChangePassword{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	claims := util.GetClaimsFromContext(c)
	id := int(claims["id"].(float64))

	if err := db.WithContext(ctx).First(new(model.Akun), "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "user tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := db.WithContext(ctx).Table("akun").Where("id", id).Update("password", util.HashPassword(request.PasswordBaru)).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func ResetPasswordHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Akun), "id", id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, map[string]string{"message": "user tidak ditemukan"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	password := util.GeneratePassword()

	if err := db.WithContext(ctx).Table("akun").Where("id", id).Update("password", util.HashPassword(password)).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, map[string]string{"password": password})
}
