package response

import (
	"BE-6/src/util"
	"strings"
	"time"

	"gorm.io/gorm"
)

type KerjasamaMOU struct {
	ID               int                `json:"id"`
	JenisDokumen     string             `json:"jenis_dokumen"`
	NomorDokumen     string             `json:"nomor_dokumen"`
	JenisKerjasama   string             `json:"jenis_kerjasama"`
	Judul            string             `json:"judul"`
	Keterangan       string             `json:"keterangan"`
	Mitra            []MitraKerjasama   `gorm:"foreignKey:IdKerjasama" json:"mitra"`
	Status           string             `json:"status"`
	TanggalAwal      string             `json:"tanggal_awal"`
	TanggalBerakhir  string             `json:"tanggal_berakhir"`
	Dokumen          string             `json:"dokumen"`
	KategoriKegiatan []KategoriKegiatan `gorm:"foreignKey:Idkerjasama" json:"kategori_kegiatan"`
}

func (p *KerjasamaMOU) AfterFind(tx *gorm.DB) (err error) {
	p.TanggalAwal = strings.Split(p.TanggalAwal, "T")[0]

	p.TanggalBerakhir = strings.Split(p.TanggalBerakhir, "T")[0]

	today := time.Now()
	tanggalBerakhir, err := util.ConvertStringToDate(p.TanggalBerakhir)
	if today.Before(tanggalBerakhir) {
		p.Status = "Aktif"
	} else {
		p.Status = "Kadaluarsa"
	}

	return
}

func (KerjasamaMOU) TableName() string {
	return "kerjasama"
}
