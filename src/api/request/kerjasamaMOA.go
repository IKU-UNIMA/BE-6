package request

import (
	"BE-6/src/model"
	"BE-6/src/util"
	"errors"
)

type KerjasamaMOA struct {
	IdFakultas            int                `json:"id_fakultas" validate:"required"`
	NomorDokumen          string             `json:"nomor_dokumen" validate:"required"`
	JenisKerjasama        string             `json:"jenis_kerjasama" validate:"required"`
	DasarDokumenKerjasama int                `json:"dasar_dokumen_kerjasama" validate:"required"`
	Judul                 string             `json:"judul" validate:"required"`
	Keterangan            string             `json:"keterangan"`
	Anggaran              float64            `json:"anggaran"`
	SumberPendanaan       string             `json:"sumber_pendanaan"`
	Mitra                 []MitraKerjasama   `json:"mitra" validate:"required"`
	KategoriKegiatan      []KategoriKegiatan `json:"kategori_kegiatan" validate:"required"`
	TanggalAwal           string             `json:"tanggal_awal" validate:"required"`
	TanggalBerakhir       string             `json:"tanggal_berakhir" validate:"required"`
}

func (r *KerjasamaMOA) MapRequest() (*model.Kerjasama, error) {
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
		IdFakultas:      r.IdFakultas,
		JenisDokumen:    "Memorandum of Aggreement (MoA)",
		NomorDokumen:    r.NomorDokumen,
		IdDasarDokumen:  r.DasarDokumenKerjasama,
		JenisKerjasama:  r.JenisKerjasama,
		Judul:           r.Judul,
		Anggaran:        r.Anggaran,
		SumberPendanaan: r.SumberPendanaan,
		Keterangan:      r.Keterangan,
		TanggalAwal:     tanggalAwal,
		TanggalBerakhir: tanggalBerakhir,
	}, nil
}
