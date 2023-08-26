package request

import (
	"BE-6/src/model"
	"BE-6/src/util"
	"errors"
)

type KerjasamaMOU struct {
	NomorDokumen     string             `json:"nomor_dokumen" validate:"required"`
	JenisKerjasama   string             `json:"jenis_kerjasama" validate:"required"`
	Judul            string             `json:"judul" validate:"required"`
	Keterangan       string             `json:"keterangan"`
	Mitra            []MitraKerjasama   `json:"mitra" validate:"required"`
	KategoriKegiatan []KategoriKegiatan `json:"kategori_kegiatan" validate:"required"`
	TanggalAwal      string             `json:"tanggal_awal" validate:"required"`
	TanggalBerakhir  string             `json:"tanggal_berakhir" validate:"required"`
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
		JenisDokumen:    "Memorandum of Understanding (MoU)",
		NomorDokumen:    r.NomorDokumen,
		JenisKerjasama:  r.JenisKerjasama,
		Judul:           r.Judul,
		Keterangan:      r.Keterangan,
		TanggalAwal:     tanggalAwal,
		TanggalBerakhir: tanggalBerakhir,
	}, nil
}
