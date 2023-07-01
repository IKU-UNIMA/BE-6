package request

import (
	"BE-6/src/model"
	"BE-6/src/util"
	"errors"
)

type KerjasamaIA struct {
	IdProdi               int    `form:"id_prodi" validate:"required"`
	NomorDokumen          string `form:"nomor_dokumen" validate:"required"`
	JenisKerjasama        string `form:"jenis_kerjasama" validate:"required"`
	DasarDokumenKerjasama int    `form:"dasar_dokumen_kerjasama" validate:"required"`
	Judul                 string `form:"judul" validate:"required"`
	Keterangan            string `form:"keterangan"`
	KategoriKegiatan      []int  `form:"kategori_kegiatan"`
	Mitra                 []MitraKerjasama
	TanggalAwal           string `form:"tanggal_awal" validate:"required"`
	TanggalBerakhir       string `form:"tanggal_berakhir" validate:"required"`
}

func (r *KerjasamaIA) MapRequest() (*model.Kerjasama, error) {
	tanggalAwal, err := util.ConvertStringToDate(r.TanggalAwal)
	if err != nil {
		return nil, errors.New("format tanggal awal salah")
	}

	tanggalBerakhir, err := util.ConvertStringToDate(r.TanggalBerakhir)
	if err != nil {
		return nil, errors.New("format tanggal berakhir salah")
	}

	if tanggalBerakhir.Before(tanggalAwal) {
		return nil, errors.New("tanggal awal tidak boleh melebihi tanggal berakhir")
	}

	return &model.Kerjasama{
		IdProdi:          r.IdProdi,
		JenisDokumen:     "Implementation Arrangement (IA)",
		NomorDokumen:     r.NomorDokumen,
		IdDasarDokumen:   r.DasarDokumenKerjasama,
		JenisKerjasama:   r.JenisKerjasama,
		Judul:            r.Judul,
		Keterangan:       r.Keterangan,
		TanggalAwal:      tanggalAwal,
		TanggalBerakhir:  tanggalBerakhir,
		KategoriKegiatan: MapToKategoriKegiatan(r.KategoriKegiatan),
	}, nil
}

func MapToKategoriKegiatan(ids []int) (result []model.KategoriKegiatan) {
	for i := 0; i < len(ids); i++ {
		result = append(result, model.KategoriKegiatan{ID: ids[i]})
	}
	return result
}
