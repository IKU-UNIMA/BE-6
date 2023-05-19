package response

import "time"

type KerjasamaMOU struct {
	ID              int       `json:"id"`
	JenisDokumen    string    `json:"jenis_dokumen"`
	NomorDokumen    string    `json:"nomor_dokumen"`
	JenisKerjasama  string    `json:"jenis_kerjasama"`
	Judul           string    `json:"judul"`
	Keterangan      string    `json:"keterangan"`
	Mitra           string    `json:"mitra"`
	Kegiatan        string    `json:"kegiatan"`
	Status          string    `json:"status"`
	TanggalAwal     time.Time `json:"tanggal_awal"`
	TanggalBerakhir time.Time `json:"tanggal_akhir"`
}

func (KerjasamaMOU) TableName() string {
	return "kerjasama"
}
