package response

import "time"

type KerjasamaIA struct {
	ID              int            `json:"id"`
	IdProdi         int            `json:"-"`
	JenisDokumen    string         `json:"jenis_dokumen"`
	NomorDokumen    string         `json:"nomor_dokumen"`
	JenisKerjasama  string         `json:"jenis_kerjasama"`
	Judul           string         `json:"judul"`
	Keterangan      string         `json:"keterangan"`
	Mitra           string         `json:"mitra"`
	Kegiatan        string         `json:"kegiatan"`
	Status          string         `json:"status"`
	TanggalAwal     time.Time      `json:"tanggal_awal"`
	TanggalBerakhir time.Time      `json:"tanggal_akhir"`
	Prodi           ProdiReference `gorm:"foreignKey:IdProdi" json:"prodi"`
}

func (KerjasamaIA) TableName() string {
	return "kerjasama"
}
