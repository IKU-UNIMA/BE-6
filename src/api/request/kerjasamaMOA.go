package request

import (
	"BE-6/src/model"
	"BE-6/src/util"
	"errors"
)

type KerjasamaMOA struct {
	IdFakultas            int    `form:"id_fakultas" validate:"required"`
	NomorDokumen          string `form:"nomor_dokumen" validate:"required"`
	JenisKerjasama        string `form:"jenis_kerjasama" validate:"required"`
	DasarDokumenKerjasama int    `form:"dasar_dokumen_kerjasama" validate:"required"`
	Judul                 string `form:"judul" validate:"required"`
	Keterangan            string `form:"keterangan"`
	Mitra                 []MitraKerjasama
	Kegiatan              string `form:"kegiatan" validate:"required"`
	Status                string `form:"status" validate:"required"`
	TanggalAwal           string `form:"tanggal_awal" validate:"required"`
	TanggalBerakhir       string `form:"tanggal_akhir" validate:"required"`
}

func (r *KerjasamaMOA) MapRequest(dokumen string) (*model.Kerjasama, error) {
	tanggalAwal, err := util.ConvertStringToDate(r.TanggalAwal)
	if err != nil {
		return nil, errors.New("format tanggal salah")
	}

	tanggalBerakhir, err := util.ConvertStringToDate(r.TanggalBerakhir)
	if err != nil {
		return nil, errors.New("format tanggal salah")
	}
	return &model.Kerjasama{
		IdFakultas:       r.IdFakultas,
		JenisDokumen:     "Memorandum of Aggreement (MoA)",
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
	}, nil
}
