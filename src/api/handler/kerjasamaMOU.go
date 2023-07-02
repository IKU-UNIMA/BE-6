package handler

import (
	"BE-6/src/api/request"
	"BE-6/src/api/response"
	"BE-6/src/config/database"
	"BE-6/src/config/env"
	"BE-6/src/config/storage"
	"BE-6/src/model"
	"BE-6/src/util"
	"BE-6/src/util/validation"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type mouQueryParam struct {
	NomorDokumen string `query:"nomor_dokumen"`
	Judul        string `query:"judul"`
	Mitra        string `query:"mitra"`
	Page         int    `query:"page"`
}

func GetAllKerjasamaMOUHandler(c echo.Context) error {
	queryParams := &mouQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	condition := "jenis_dokumen='Memorandum of Understanding (MoU)'"

	if queryParams.NomorDokumen != "" {

		condition += " AND UPPER(nomor_dokumen) LIKE '%" + strings.ToUpper(queryParams.NomorDokumen) + "%'"

	}

	if queryParams.Judul != "" {

		condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Judul) + "%'"

	}

	db := database.DB
	ctx := c.Request().Context()
	limit := 20
	data := []response.KerjasamaMOU{}

	if queryParams.Mitra != "" {

		condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Mitra) + "%'"

	}

	if err := db.WithContext(ctx).Debug().Table("kerjasama").
		Where(condition).
		Offset(util.CountOffset(queryParams.Page, limit)).
		Limit(limit).Where(condition).
		Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var totalResult int64
	if err := db.WithContext(ctx).Table("kerjasama").
		Where(condition).Count(&totalResult).Error; err != nil {
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

// func GetAllKerjasamaMOUHandler(c echo.Context) error {
// 	db := database.DB
// 	ctx := c.Request().Context()
// 	result := []response.KerjasamaMOU{}

// 	if err := db.WithContext(ctx).Order("id").Where("jenis_dokumen", "Memorandum of Aggreement (MOU)").Preload("Prodi").Find(&result).Error; err != nil {
// 		return util.FailedResponse(http.StatusInternalServerError, nil)
// 	}

// 	return util.SuccessResponse(c, http.StatusOK, result)
// }

func GetKerjasamaMOUByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	result := &response.KerjasamaMOU{}

	if err := db.WithContext(ctx).
		Where("jenis_dokumen", "Memorandum of Understanding (MoU)").
		Preload("Mitra").Preload("KategoriKegiatan").
		First(result, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, result)
}

// func GetKerjasamaMOUByUnitKerjaHandler(c echo.Context) error {
// 	id, err := util.GetId(c)
// 	if err != nil {
// 		return err
// 	}

// 	db := database.DB
// 	ctx := c.Request().Context()
// 	result := &response.KerjasamaMOU{}

// 	if err := db.WithContext(ctx).Where("prodi", result.IdProdi).Preload("Prodi").First(result, id).Error; err != nil {
// 		if err.Error() == util.NOT_FOUND_ERROR {
// 			return util.FailedResponse(http.StatusNotFound, nil)
// 		}

// 		return util.FailedResponse(http.StatusInternalServerError, nil)
// 	}

// 	return util.SuccessResponse(c, http.StatusOK, result)
// }

func InsertKerjasamaMOUHandler(c echo.Context) error {
	request := &request.KerjasamaMOU{}
	reqData := c.FormValue("mitra")

	if err := json.Unmarshal([]byte(reqData), &request.Mitra); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"mitra": err.Error()})
	}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	dokumen, _ := c.FormFile("file")
	if dokumen == nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "file tidak boleh kosong"})
	}

	if err := util.CheckFileIsPDF(dokumen); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if len(request.KategoriKegiatan) == 0 {
		form, _ := c.MultipartForm()
		kategoriKegiatan := form.Value["kategori_kegiatan[]"]
		for _, v := range kategoriKegiatan {
			id, err := strconv.Atoi(v)
			if err != nil {
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "id kategori kegiatan harus berupa angka"})
			}

			request.KategoriKegiatan = append(request.KategoriKegiatan, id)
		}
	}

	if len(request.KategoriKegiatan) == 0 {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"kategori_kegiatan": "field ini wajib diisi"})
	}

	// validate kategori kegiatan
	for _, v := range request.KategoriKegiatan {
		if err := db.WithContext(ctx).Select("id").First(new(model.KategoriKegiatan), "id", v).Error; err != nil {
			if err.Error() == util.NOT_FOUND_ERROR {
				return util.FailedResponse(
					http.StatusNotFound,
					map[string]string{"message": fmt.Sprintf("kategori kegiatan dengan id '%d' tidak ditemukan", v)},
				)
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	mitra := []model.MitraKerjasama{}
	for _, v := range request.Mitra {
		if err := validation.ValidateKerjasama(&v); err != nil {
			return err
		}

		mitra = append(mitra, *v.MapRequestToKerjasama())
	}

	data, err := request.MapRequest()
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	data.Mitra = mitra

	dDokumen, err := storage.CreateFile(dokumen, env.GetDokumenFolderId())
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.Dokumen = util.CreateFileUrl(dDokumen.Id)

	if err := db.WithContext(ctx).Create(data).Error; err != nil {
		storage.DeleteFile(dDokumen.Id)
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "nomor surat duplikasi"})
		}

		return nil
	}

	return util.SuccessResponse(c, http.StatusCreated, data.ID)
}

func EditKerjasamaMOUHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	request := &request.KerjasamaMOU{}
	if err := c.Bind(request); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(request); err != nil {
		return err
	}

	db := database.DB
	tx := db.Begin()
	ctx := c.Request().Context()

	if err := db.WithContext(ctx).First(new(model.Kerjasama), id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if len(request.KategoriKegiatan) == 0 {
		form, _ := c.MultipartForm()
		kategoriKegiatan := form.Value["kategori_kegiatan[]"]
		for _, v := range kategoriKegiatan {
			id, err := strconv.Atoi(v)
			if err != nil {
				return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "id kategori kegiatan harus berupa angka"})
			}

			request.KategoriKegiatan = append(request.KategoriKegiatan, id)
		}
	}

	if len(request.KategoriKegiatan) == 0 {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"kategori_kegiatan": "field ini tidak boleh kosong"})
	}

	// validate kategori kegiatan
	for _, v := range request.KategoriKegiatan {
		if err := db.WithContext(ctx).Select("id").First(new(model.KategoriKegiatan), "id", v).Error; err != nil {
			if err.Error() == util.NOT_FOUND_ERROR {
				return util.FailedResponse(
					http.StatusNotFound,
					map[string]string{"message": fmt.Sprintf("kategori kegiatan dengan id '%d' tidak ditemukan", v)},
				)
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}
	}

	data, errMapping := request.MapRequest()
	if errMapping != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": errMapping.Error()})
	}

	reqData := c.FormValue("mitra")

	if reqData != "" {
		if err := json.Unmarshal([]byte(reqData), &request.Mitra); err != nil {
			tx.Rollback()
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"mitra": err.Error()})
		}
	}

	mitra := []model.MitraKerjasama{}
	for _, v := range request.Mitra {
		if err := validation.ValidateKerjasama(&v); err != nil {
			return err
		}

		mitra = append(mitra, *v.MapRequestToKerjasama())
	}

	if err := tx.WithContext(ctx).Omit("jenis_dokumen", "dokumen", "KategoriKegiatan").
		Where("jenis_dokumen = 'Memorandum of Understanding (MoU)' AND id = ?", id).Updates(data).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "nomor surat tidak boleh sama"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Delete(new(model.MitraKerjasama), "id_kerjasama", id).Error; err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	mou := &model.Kerjasama{ID: id}
	if err := tx.WithContext(ctx).Model(mou).Association("Mitra").Replace(&mitra); err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Model(mou).Association("KategoriKegiatan").Replace(&data.KategoriKegiatan); err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteKerjasamaMOUHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
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
