package request

import "BE-6/src/model"

type JenisKategoriKegiatan struct {
	Nama string `json:"nama"`
}

func (r *JenisKategoriKegiatan) MapRequest() *model.JenisKategoriKegiatan {
	return &model.JenisKategoriKegiatan{Nama: r.Nama}
}
