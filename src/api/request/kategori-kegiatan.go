package request

import "BE-6/src/model"

type KategoriKegiatan struct {
	IdJenisKategoriKegiatan int     `json:"id_jenis_kategori_kegiatan"`
	NilaiKontrak            float64 `json:"nilai_kontrak"`
	Volume                  string  `json:"volume"`
	At                      string  `json:"at"`
	Keterangan              string  `json:"keterangan"`
	Sasaran                 string  `json:"sasaran"`
	IndikatorKinerja        string  `json:"indikator_kinerja"`
}

func (r *KategoriKegiatan) MapRequest() *model.KategoriKegiatan {
	return &model.KategoriKegiatan{
		IdJenisKategoriKegiatan: r.IdJenisKategoriKegiatan,
		NilaiKontrak:            r.NilaiKontrak,
		Volume:                  r.Volume,
		At:                      r.At,
		Keterangan:              r.Keterangan,
		Sasaran:                 r.Sasaran,
		IndikatorKinerja:        r.IndikatorKinerja,
	}
}
