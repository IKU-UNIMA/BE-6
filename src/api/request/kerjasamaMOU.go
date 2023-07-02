package request

import (
	"BE-6/src/model"
	"BE-6/src/util"
	"errors"
)

type KerjasamaMOU struct {
	NomorDokumen     string `form:"nomor_dokumen" validate:"required"`
	JenisKerjasama   string `form:"jenis_kerjasama" validate:"required"`
	Judul            string `form:"judul" validate:"required"`
	Keterangan       string `form:"keterangan"`
	Mitra            []MitraKerjasama
	KategoriKegiatan []int  `form:"kategori_kegiatan"`
	TanggalAwal      string `form:"tanggal_awal" validate:"required"`
	TanggalBerakhir  string `form:"tanggal_berakhir" validate:"required"`
}

func (r *KerjasamaMOU) MapRequest() (*model.Kerjasama, error) {
	tanggalAwal, err := util.ConvertStringToDate(r.TanggalAwal)
	if err != nil {
		return nil, errors.New("format tanggal salah")
	}

	tanggalBerakhir, err := util.ConvertStringToDate(r.TanggalBerakhir)
	if err != nil {
		return nil, errors.New("format tanggal salah")
	}

	if tanggalBerakhir.Before(tanggalAwal) {
		return nil, errors.New("tanggal awal tidak boleh melebihi tanggal berakhir")
	}

	return &model.Kerjasama{
		JenisDokumen:     "Memorandum of Understanding (MoU)",
		NomorDokumen:     r.NomorDokumen,
		JenisKerjasama:   r.JenisKerjasama,
		Judul:            r.Judul,
		Keterangan:       r.Keterangan,
		TanggalAwal:      tanggalAwal,
		TanggalBerakhir:  tanggalBerakhir,
		KategoriKegiatan: MapBatchIDKategoriKegiatan(r.KategoriKegiatan),
	}, nil
}
