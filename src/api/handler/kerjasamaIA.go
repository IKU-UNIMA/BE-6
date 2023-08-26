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
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

type iaQueryParam struct {
	Prodi        int    `query:"prodi"`
	NomorDokumen string `query:"nomor_dokumen"`
	Judul        string `query:"judul"`
	Mitra        string `query:"mitra"`
	Page         int    `query:"page"`
}

func GetDasarKerjasamaIAHandler(c echo.Context) error {
	queryParams := &iaQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	condition := "jenis_dokumen IN ('Memorandum of Aggreement (MoA)', 'Memorandum of Understanding (MoU)')"

	if queryParams.Judul != "" {
		condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Judul) + "%'"
	}

	db := database.DB
	ctx := c.Request().Context()
	data := []response.DasarKerjasama{}

	if err := db.WithContext(ctx).Where(condition).Find(&data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetAllKerjasamaIAHandler(c echo.Context) error {
	queryParams := &iaQueryParam{}
	if err := (&echo.DefaultBinder{}).BindQueryParams(c, queryParams); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	condition := "jenis_dokumen='Implementation Arrangement (IA)'"
	if queryParams.Prodi != 0 {
		condition += fmt.Sprintf(` AND id_prodi = %d`, queryParams.Prodi)
	}

	if queryParams.NomorDokumen != "" {

		condition += " AND UPPER(nomor_dokumen) LIKE '%" + strings.ToUpper(queryParams.NomorDokumen) + "%'"

	}

	if queryParams.Judul != "" {

		condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Judul) + "%'"

	}

	db := database.DB
	ctx := c.Request().Context()
	limit := 20
	data := []response.KerjasamaIA{}

	if queryParams.Mitra != "" {
		condition += " AND UPPER(judul) LIKE '%" + strings.ToUpper(queryParams.Mitra) + "%'"

	}

	if err := db.WithContext(ctx).Debug().Table("kerjasama").
		Where(condition).Preload("Prodi").Preload("Mitra").
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

func GetKerjasamaIAByIdHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	db := database.DB
	ctx := c.Request().Context()
	result := &response.KerjasamaIA{}

	if err := db.WithContext(ctx).
		Where("jenis_dokumen", "Implementation Arrangement (IA)").
		Preload("Prodi").Preload("Mitra").Preload("KategoriKegiatan.JenisKategoriKegiatan").
		First(result, id).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	newResponse := response.DasarKerjasama{}

	if err := db.WithContext(ctx).
		Where("jenis_dokumen IN ('Memorandum of Aggreement (MoA)', 'Memorandum of Understanding (MoU)') AND id=?", result.IdDasarDokumen).
		Find(&newResponse).Error; err != nil {

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}
	result.DasarDokumenKerjasama = newResponse

	return util.SuccessResponse(c, http.StatusOK, result)
}

// func GetKerjasamaIAByUnitKerjaHandler(c echo.Context) error {
// 	id, err := util.GetId(c)
// 	if err != nil {
// 		return err
// 	}

// 	db := database.DB
// 	ctx := c.Request().Context()
// 	result := &response.KerjasamaIA{}

// 	if err := db.WithContext(ctx).Where("prodi", result.IdProdi).Preload("Prodi").First(result, id).Error; err != nil {
// 		if err.Error() == util.NOT_FOUND_ERROR {
// 			return util.FailedResponse(http.StatusNotFound, nil)
// 		}

// 		return util.FailedResponse(http.StatusInternalServerError, nil)
// 	}

// 	return util.SuccessResponse(c, http.StatusOK, result)
// }

func ImportKerjasamaIAHandler(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.DB
	ctx := c.Request().Context()
	data := []model.Kerjasama{}

	defer func() {
		os.Remove(file.Filename)
	}()

	if err := util.WriteFile(file); err != nil {
		return err
	}

	excel, err := excelize.OpenFile(file.Filename)
	if err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}
	defer excel.Close()

	rows, err := excel.GetRows(excel.GetSheetName(0))
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if len(rows[0]) != 8 {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "jumlah kolom tidak sesuai format"})
	}

	for i := 1; i < len(rows); i++ {

		data = append(data, model.Kerjasama{
			JenisDokumen: rows[i][1],
			NomorDokumen: rows[i][2],
			Judul:        rows[i][3],
			Keterangan:   rows[i][4],
			Kegiatan:     rows[i][6],
		})
	}

	if err := db.WithContext(ctx).Create(&data).Error; err != nil {
		if strings.Contains(err.Error(), "nomor_dokumen") {
			return util.FailedResponse(
				http.StatusBadRequest,
				map[string]string{"message": "terdapat duplikasi untuk nomor dokumen " + strings.Split(err.Error(), "'")[1]},
			)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func InsertKerjasamaIAHandler(c echo.Context) error {
	req := &request.KerjasamaIA{}
	reqData := c.FormValue("data")
	if err := json.Unmarshal([]byte(reqData), req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
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

	if err := db.WithContext(ctx).
		Where("jenis_dokumen IN ('Memorandum of Aggreement (MoA)', 'Memorandum of Understanding (MoU)') AND id=?", req.DasarDokumenKerjasama).
		First(new(model.Kerjasama)).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// validate kategori kegiatan
	kategoriKegiatan := []model.KategoriKegiatan{}
	for _, v := range req.KategoriKegiatan {
		if err := db.WithContext(ctx).Select("id").
			First(new(model.JenisKategoriKegiatan), "id", v.IdJenisKategoriKegiatan).Error; err != nil {
			if err.Error() == util.NOT_FOUND_ERROR {
				return util.FailedResponse(
					http.StatusNotFound,
					map[string]string{"message": fmt.Sprintf("kategori kegiatan dengan id '%d' tidak ditemukan", v.IdJenisKategoriKegiatan)},
				)
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		kategoriKegiatan = append(kategoriKegiatan, *v.MapRequest())
	}

	mitra := []model.MitraKerjasama{}
	for _, v := range req.Mitra {
		if err := validation.ValidateKerjasama(&v); err != nil {
			return err
		}

		mitra = append(mitra, *v.MapRequestToKerjasama())
	}

	data, err := req.MapRequest()
	if err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	data.KategoriKegiatan = kategoriKegiatan
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

func EditKerjasamaIAHandler(c echo.Context) error {
	id, err := util.GetId(c)
	if err != nil {
		return err
	}

	req := &request.KerjasamaIA{}
	reqData := c.FormValue("data")
	if err := json.Unmarshal([]byte(reqData), req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"kategori_kegiatan": err.Error()})
	}

	if err := c.Validate(req); err != nil {
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

	if err := db.WithContext(ctx).
		Where("jenis_dokumen IN ('Memorandum of Aggreement (MoA)', 'Memorandum of Understanding (MoU)') AND id=?", req.DasarDokumenKerjasama).
		First(new(model.Kerjasama)).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	// validate kategori kegiatan
	kategoriKegiatan := []model.KategoriKegiatan{}
	for _, v := range req.KategoriKegiatan {
		if err := db.WithContext(ctx).Select("id").First(new(model.KategoriKegiatan), "id", v.IdJenisKategoriKegiatan).Error; err != nil {
			if err.Error() == util.NOT_FOUND_ERROR {
				return util.FailedResponse(
					http.StatusNotFound,
					map[string]string{"message": fmt.Sprintf("kategori kegiatan dengan id '%d' tidak ditemukan", v.IdJenisKategoriKegiatan)},
				)
			}

			return util.FailedResponse(http.StatusInternalServerError, nil)
		}

		kategoriKegiatan = append(kategoriKegiatan, *v.MapRequest())
	}

	mitra := []model.MitraKerjasama{}
	for _, v := range req.Mitra {
		if err := validation.ValidateKerjasama(&v); err != nil {
			return err
		}

		mitra = append(mitra, *v.MapRequestToKerjasama())
	}

	data, errMapping := req.MapRequest()
	if errMapping != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": errMapping.Error()})
	}

	if err := tx.WithContext(ctx).Omit("jenis_dokumen", "dokumen", "KategoriKegiatan").
		Where("jenis_dokumen = 'Implementation Arrangement (IA)' AND id = ?", id).Updates(data).Error; err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), util.UNIQUE_ERROR) {
			return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": "nomor surat tidak boleh sama"})
		}

		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Delete(new(model.KategoriKegiatan), "id_kerjasama", id).Error; err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Delete(new(model.MitraKerjasama), "id_kerjasama", id).Error; err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	ia := &model.Kerjasama{ID: id}
	if err := tx.WithContext(ctx).Model(ia).Association("Mitra").Replace(&mitra); err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.WithContext(ctx).Model(ia).Association("KategoriKegiatan").Replace(&kategoriKegiatan); err != nil {
		tx.Rollback()
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	if err := tx.Commit().Error; err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}

func DeleteKerjasamaIAHandler(c echo.Context) error {
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
