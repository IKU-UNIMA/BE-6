package request

import "BE-6/src/model"

type MitraKerjasama struct {
	NamaInstansi           string `json:"nama_instansi"`
	NegaraAsal             string `json:"negara_asal"`
	Bidang                 string `json:"bidang"`
	Penandatangan          string `json:"penandatangan"`
	JabatanPenandatangan   string `json:"jabatan_penandatangan"`
	PenanggungJawab        string `json:"penanggung_jawab"`
	JabatanPenanggungJawab string `json:"jabatan_penanggung_jawab"`
}

func (r *MitraKerjasama) MapRequestToKerjasama() *model.MitraKerjasama {
	return &model.MitraKerjasama{
		NamaInstansi:           r.NamaInstansi,
		NegaraAsal:             r.NegaraAsal,
		Bidang:                 r.Bidang,
		Penandatangan:          r.Penandatangan,
		JabatanPenandatangan:   r.JabatanPenandatangan,
		PenanggungJawab:        r.PenanggungJawab,
		JabatanPenanggungJawab: r.JabatanPenandatangan,
	}
}
