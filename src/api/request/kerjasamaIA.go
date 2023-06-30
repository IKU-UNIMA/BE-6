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
	Mitra                 []MitraKerjasama
	Kegiatan              string `form:"kegiatan" validate:"required"`
	Status                string `form:"status" validate:"required"`
	TanggalAwal           string `form:"tanggal_awal" validate:"required"`
	TanggalBerakhir       string `form:"tanggal_berakhir" validate:"required"`
	KategoriKegiatan      []int  `form:"kategori_kegiatan" validate:"required"`
}

func (r *KerjasamaIA) MapRequest(dokumen string) (*model.Kerjasama, error) {
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
		Kegiatan:         r.Kegiatan,
		Status:           r.Status,
		TanggalAwal:      tanggalAwal,
		TanggalBerakhir:  tanggalBerakhir,
		Dokumen:          dokumen,
		KategoriKegiatan: MapToKategoriKegiatan(r.KategoriKegiatan),
	}, nil
}

func MapToKategoriKegiatan(ids []int) []model.KategoriKegiatan {
	result := make([]model.KategoriKegiatan, len(ids))
	if len(ids) == 1 {
		result = append(result, model.KategoriKegiatan{ID: ids[0]})
	}

	for i := 0; i < len(ids)/2; i++ {
		result[i] = model.KategoriKegiatan{ID: ids[i]}
		result[len(ids)-i-1] = model.KategoriKegiatan{ID: ids[len(ids)-i-1]}
	}

	return result
}
