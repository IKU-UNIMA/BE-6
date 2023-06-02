package handler

import (
	"BE-6/src/api/request"
	"BE-6/src/api/response"
	"BE-6/src/config/database"
	"BE-6/src/util"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type dashboardPathParam struct {
	Tahun    int `param:"tahun"`
	Fakultas int `param:"fakultas"`
}

func GetDashboardHandler(c echo.Context) error {
	params := &dashboardPathParam{}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, params); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.Dashboard{}

	conds := fmt.Sprintf("YEAR(tanggal_awal) <= %d AND YEAR(tanggal_berakhir) >= %d", params.Tahun, params.Tahun)
	dashboardQuery := fmt.Sprintf(`
	SELECT
		SUM(IF(jenis_dokumen='Memorandum of Understanding (MoU)', 1, 0))) AS mou,
		SUM(IF(jenis_dokumen='Memorandum of Agreement (MoA)', 1, 0))) AS moa,
		SUM(IF(jenis_dokumen='Implementation Arrangement (IA)', 1, 0))) AS ia,
		SUM(IF(jenis_kerjasama='Kerjasama Luar Negeri', 1, 0))) AS luar_negeri,
		SUM(IF(jenis_kerjasama='Kerjasama Dalam Negeri', 1, 0))) AS dalam_negeri
	FROM kerjasama
	WHERE %s
	`, conds)

	if err := db.WithContext(ctx).Raw(dashboardQuery).Find(data).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	prodi := []struct {
		ID       int
		Fakultas string
		Jumlah   int
	}{}
	prodiQuery := `
	SELECT fakultas.id, fakultas.nama AS fakultas, COUNT(prodi.id) AS jumlah FROM fakultas
	left JOIN prodi ON prodi.id_fakultas = fakultas.id
	GROUP BY fakultas.id ORDER BY fakultas.id
	`
	if err := db.WithContext(ctx).Raw(prodiQuery).Find(&prodi).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	var target float64
	targetQuery := fmt.Sprintf(`
	SELECT target FROM target
	WHERE bagian = 'IKU 6' AND tahun = %d
	`, params.Tahun)
	if err := db.WithContext(ctx).Raw(targetQuery).Find(&target).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.Target = fmt.Sprintf("%.1f", util.RoundFloat(target))

	query := fmt.Sprintf(`
	SELECT COUNT(DISTINCT id_prodi) AS jumlah_pencapaian FROM fakultas
	LEFT JOIN prodi ON prodi.id_fakultas = fakultas.id
	LEFT JOIN kerjasama ON kerjasama.id_prodi = prodi.id AND kerjasama.jenis_dokumen='Implementation Arrangement (IA)'
		AND %s
	GROUP BY fakultas.id ORDER BY fakultas.id
	`, conds)

	if err := db.WithContext(ctx).Raw(query).Find(&data.Detail).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	for i := 0; i < len(prodi); i++ {
		data.Total += data.Detail[i].JumlahPencapaian
		jumlahProdi := prodi[i].Jumlah
		data.TotalProdi += jumlahProdi

		var persentase float64
		if data.Detail[i].JumlahPencapaian != 0 {
			persentase = util.RoundFloat((float64(data.Detail[i].JumlahPencapaian) / float64(prodi[i].Jumlah)) * 100)
		}

		data.Detail[i].ID = prodi[i].ID
		data.Detail[i].Fakultas = prodi[i].Fakultas
		data.Detail[i].JumlahProdi = jumlahProdi
		data.Detail[i].Persentase = fmt.Sprintf("%.2f", persentase) + "%"
	}

	pencapaian := util.RoundFloat((float64(data.Total) / float64(data.TotalProdi)) * 100)

	data.Pencapaian = fmt.Sprintf("%.2f", pencapaian) + "%"

	return util.SuccessResponse(c, http.StatusOK, data)
}

func GetDashboardByFakultasHandler(c echo.Context) error {
	params := &dashboardPathParam{}
	if err := (&echo.DefaultBinder{}).BindPathParams(c, params); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	fakultasConds := fmt.Sprintf("WHERE prodi.id_fakultas=%d", params.Fakultas)

	db := database.InitMySQL()
	ctx := c.Request().Context()
	data := &response.DashboardPerProdi{}

	fakultas := ""
	if err := db.WithContext(ctx).Raw("SELECT nama FROM fakultas WHERE id = ?", params.Fakultas).First(&fakultas).Error; err != nil {
		if err.Error() == util.NOT_FOUND_ERROR {
			return util.FailedResponse(http.StatusNotFound, nil)
		}
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.Fakultas = fakultas

	conds := fmt.Sprintf("AND YEAR(tanggal_awal) <= %d AND YEAR(tanggal_berakhir) >= %d", params.Tahun, params.Tahun)
	query := fmt.Sprintf(`
	SELECT
		CONCAT(kode_prodi, " - ", prodi.nama, " (", prodi.jenjang, ")") AS prodi,
		COUNT(kerjasama.id) AS jumlah_kerjasama,
		COUNT(DISTINCT id_prodi) AS capaian
	FROM prodi
	LEFT JOIN kerjasama ON kerjasama.id_prodi = prodi.id
		AND kerjasama.jenis_dokumen='Implementation Arrangement (IA)'
		%s
	%s
	GROUP BY prodi.id ORDER BY prodi.id
	`, conds, fakultasConds)

	if err := db.WithContext(ctx).Raw(query).Find(&data.Detail).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	data.TotalProdi = len(data.Detail)
	for i := 0; i < len(data.Detail); i++ {
		data.Total += data.Detail[i].Capaian
	}

	pencapaian := util.RoundFloat((float64(data.Total) / float64(data.TotalProdi)) * 100)
	data.Pencapaian = fmt.Sprintf("%.2f", pencapaian) + "%"

	return util.SuccessResponse(c, http.StatusOK, data)
}

func InsertTargetHandler(c echo.Context) error {
	req := &request.Target{}
	if err := c.Bind(req); err != nil {
		return util.FailedResponse(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	db := database.InitMySQL()
	ctx := c.Request().Context()
	conds := fmt.Sprintf("bagian='%s' AND tahun=%d", util.IKU6, req.Tahun)

	if err := db.WithContext(ctx).Where(conds).Save(req.MapRequest()).Error; err != nil {
		return util.FailedResponse(http.StatusInternalServerError, nil)
	}

	return util.SuccessResponse(c, http.StatusOK, nil)
}
