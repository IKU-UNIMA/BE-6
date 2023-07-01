package request

import "BE-6/src/model"

type KategoriKegiatan struct {
	Nama string `json:"nama"`
}

func (r *KategoriKegiatan) MapRequest() *model.KategoriKegiatan {
	return &model.KategoriKegiatan{Nama: r.Nama}
}

func MapBatchIDKategoriKegiatan(ids []int) (result []model.KategoriKegiatan) {
	for i := 0; i < len(ids); i++ {
		result = append(result, model.KategoriKegiatan{ID: ids[i]})
	}
	return result
}
