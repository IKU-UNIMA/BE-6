package request

import (
	"BE-6/src/model"
	"BE-6/src/util"
	"errors"
)

type KerjasamaMOU struct {
	JenisDokumen    string `json:"jenis_dokumen" validate:"required"`
	NomorDokumen    string `json:"nomor_dokumen" validate:"required"`
	JenisKerjasama  string `json:"jenis_kerjasama" validate:"required"`
	Judul           string `json:"judul" validate:"required"`
	Keterangan      string `json:"keterangan"`
	Mitra           string `json:"mitra"`
	Kegiatan        string `json:"kegiatan" validate:"required"`
	Status          string `json:"status" validate:"required"`
	TanggalAwal     string `json:"tanggal_awal" validate:"required"`
	TanggalBerakhir string `json:"tanggal_akhir" validate:"required"`
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
	return &model.Kerjasama{
		JenisDokumen:    r.JenisDokumen,
		NomorDokumen:    r.NomorDokumen,
		JenisKerjasama:  r.JenisKerjasama,
		Judul:           r.Judul,
		Keterangan:      r.Keterangan,
		Mitra:           r.Mitra,
		Kegiatan:        r.Kegiatan,
		Status:          r.Status,
		TanggalAwal:     tanggalAwal,
		TanggalBerakhir: tanggalBerakhir,
	}, nil
}
