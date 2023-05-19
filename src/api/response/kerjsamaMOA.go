package response

import "time"

type KerjasamaMOA struct {
	ID              int               `json:"id"`
	IdFakultas      int               `json:"-"`
	JenisDokumen    string            `json:"jenis_dokumen"`
	NomorDokumen    string            `json:"nomor_dokumen"`
	JenisKerjasama  string            `json:"jenis_kerjasama"`
	Judul           string            `json:"judul"`
	Keterangan      string            `json:"keterangan"`
	Mitra           string            `json:"mitra"`
	Kegiatan        string            `json:"kegiatan"`
	Status          string            `json:"status"`
	TanggalAwal     time.Time         `json:"tanggal_awal"`
	TanggalBerakhir time.Time         `json:"tanggal_akhir"`
	Fakultas        FakultasReference `gorm:"foreignKey:IdFakultas" json:"fakultas"`
}

func (KerjasamaMOA) TableName() string {
	return "kerjasama"
}
