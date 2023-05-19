package request

import "BE-6/src/model"

type Fakultas struct {
	Nama string `json:"nama" validate:"required"`
}

func (r *Fakultas) MapRequest() *model.Fakultas {
	return &model.Fakultas{
		Nama: r.Nama,
	}
}
