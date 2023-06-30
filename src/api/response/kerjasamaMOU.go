package response

import (
	"BE-6/src/util"
	"strings"
	"time"

	"gorm.io/gorm"
)

type KerjasamaMOU struct {
	ID              int              `json:"id"`
	JenisDokumen    string           `json:"jenis_dokumen"`
	NomorDokumen    string           `json:"nomor_dokumen"`
	JenisKerjasama  string           `json:"jenis_kerjasama"`
	Judul           string           `json:"judul"`
	Keterangan      string           `json:"keterangan"`
	Mitra           []MitraKerjasama `gorm:"foreignKey:IdKerjasama" json:"mitra"`
	Kegiatan        string           `json:"kegiatan"`
	Status          string           `json:"status"`
	TanggalAwal     string           `json:"tanggal_awal"`
	TanggalBerakhir string           `json:"tanggal_akhir"`
	Dokumen         string           `json:"dokumen"`
}

func (p *KerjasamaMOU) AfterFind(tx *gorm.DB) (err error) {
	today := time.Now()
	tanggalBerakhir, _ := util.ConvertStringToDate(p.TanggalBerakhir)
	if today.Before(tanggalBerakhir) {
		p.Status = "Aktif"
	} else {
		p.Status = "Kadaluarsa"
	}

	p.TanggalAwal = strings.Split(p.TanggalAwal, "T")[0]

	p.TanggalBerakhir = strings.Split(p.TanggalBerakhir, "T")[0]

	return
}

func (KerjasamaMOU) TableName() string {
	return "kerjasama"
}
