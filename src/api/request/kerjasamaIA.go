package request

import (
	"BE-6/src/model"
	"BE-6/src/util"
	"errors"
)

type KerjasamaIA struct {
	IdProdi               int                `json:"id_prodi" validate:"required"`
	NomorDokumen          string             `json:"nomor_dokumen" validate:"required"`
	JenisKerjasama        string             `json:"jenis_kerjasama" validate:"required"`
	DasarDokumenKerjasama int                `json:"dasar_dokumen_kerjasama" validate:"required"`
	Judul                 string             `json:"judul" validate:"required"`
	Keterangan            string             `json:"keterangan"`
	Anggaran              string             `json:"anggaran"`
	SumberPendanaan       string             `json:"sumber_pendanaan"`
	KategoriKegiatan      []KategoriKegiatan `json:"kategori_kegiatan" validate:"required"`
	Mitra                 []MitraKerjasama   `json:"mitra" validate:"required"`
	TanggalAwal           string             `json:"tanggal_awal" validate:"required"`
	TanggalBerakhir       string             `json:"tanggal_berakhir" validate:"required"`
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
		IdProdi:         r.IdProdi,
		JenisDokumen:    "Implementation Arrangement (IA)",
		NomorDokumen:    r.NomorDokumen,
		IdDasarDokumen:  r.DasarDokumenKerjasama,
		JenisKerjasama:  r.JenisKerjasama,
		Judul:           r.Judul,
		Keterangan:      r.Keterangan,
		TanggalAwal:     tanggalAwal,
		TanggalBerakhir: tanggalBerakhir,
	}, nil
}
