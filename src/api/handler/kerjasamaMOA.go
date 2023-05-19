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

type moaQueryParam struct {
	Fakultas     int    `query:"fakultas"`
	NomorDokumen string `query:"nomor_dokumen"`
	Judul        string `query:"judul"`
	Mitra        string `query:"mitra"`
	Page         int    `query:"page"`
}

func GetAllKerjasamaMOAHandler(c echo.Context) error {
	queryParams := &moaQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	condition := "jenis_dokumen='Memorandum of Aggreement (MoA)'"
	if queryParams.Fakultas != 0 {
		condition += fmt.Sprintf(` AND id_fakultas = %d`, queryParams.Fakultas)
	}

	if queryParams.NomorDokumen != "" {

		condition += " AND UPPER(nomor_dokumen) LIKE '%" + strings.ToUpper(queryParams.NomorDokumen) + "%'"

	}

	if queryParams.Judul != "" {

		condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Judul) + "%'"

	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	limit := 20
	data := []response.KerjasamaMOA{}

	if queryParams.Mitra != "" {

		condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Mitra) + "%'"

	}

	var totalResult int64

	// if err := db.WithContext(ctx).Order("id").Where("jenis_dokumen", "Implementation Arrangement (IA)").Preload("Prodi").Find(&result).Error; err != nil {
	// 	return util.FailedResponse(http.StatusInternalServerError, nil)
	// }
	if err := db.WithContext(ctx).Table("kerjasama").Where(condition).Preload("Fakultas").Find(&data).Count(&totalResult).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, util.Pagination{
		Limit:       limit,
		Page:        queryParams.Page,
		TotalPage:   util.CountTotalPage(int(totalResult), limit),
		TotalResult: int(totalResult),
		Data:        data,
	})
}

// func GetAllKerjasamaMOAHandler(c echo.Context) error {
// 	db := database.InitMySQL()
// 	ctx := c.Request().Context()
// 	result := []response.KerjasamaMOA{}

// 	if err := db.WithContext(ctx).Order("id").Where("jenis_dokumen", "Memorandum of Aggreement (MoA)").Preload("Prodi").Find(&result).Error; err != nil {
// 		return util.FailedResponse(http.StatusInternalServerError, nil)
// 	}

// 	return util.SuccessResponse(c, http.StatusOK, result)
// }

func GetKerjasamaMOAByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	result := &response.KerjasamaMOA{}

	if err := db.WithContext(ctx).Where("jenis_dokumen", "Memorandum of Aggreement (MoA)").Preload("Fakultas").First(result, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

// func GetKerjasamaMOAByUnitKerjaHandler(c echo.Context) error {
// 	id, err := util.GetId(c)
// 	if err != nil {
// 		return err
// 	}

// 	db := database.InitMySQL()
// 	ctx := c.Request().Context()
// 	result := &response.KerjasamaMOA{}

// 	if err := db.WithContext(ctx).Where("prodi", result.IdProdi).Preload("Prodi").First(result, id).Error; err != nil {
// 		if err.Error() == util.NOT_FOUND_ERROR {
// 			return util.FailedResponse(http.StatusNotFound, nil)
// 		}

// 		return util.FailedResponse(http.StatusInternalServerError, nil)
// 	}

// 	return util.SuccessResponse(c, http.StatusOK, result)
// }

func InsertKerjasamaMOAHandler(c echo.Context) error {
	request := &request.KerjasamaMOA{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	data, errMapping := request.MapRequest()
	if errMapping != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": errMapping.Error()})
	}
	if err := db.WithContext(ctx).Create(data).Error; err != nil {
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "fakultas sudah ada"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusCreated, data.ID)
}

func EditKerjasamaMOAHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	request := &request.KerjasamaMOA{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Kerjasama), id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}
	data, errMapping := request.MapRequest()
	if errMapping != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": errMapping.Error()})
	}
	if err := db.WithContext(ctx).Where("id", id).Updates(data).Error; err != nil {
		if err != nil {
			if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "nomor surat tidak boleh sama"})
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteKerjasamaMOAHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()

	query := db.WithContext(ctx).Delete(new(model.Kerjasama), id)
	if query.Error == nil && query.RowsAffected < 1 {
		return util.FailedResponse(http.StatusNotFound, nil)
	}

	if query.Error != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
