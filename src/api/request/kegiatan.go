package request

import "BE-6/src/model"

type KategoriKegiatan struct {
	Nama string `json:"nama"`
}

func (r *KategoriKegiatan) MapRequestToKerjasama() *model.KategoriKegiatan {
	return &model.KategoriKegiatan{
		Nama: r.Nama,
	}
}
