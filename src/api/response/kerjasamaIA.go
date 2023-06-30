package response

import (
	"BE-6/src/util"
	"strings"
	"time"

	"gorm.io/gorm"
)

type KerjasamaIA struct {
	ID                    int                `json:"id"`
	IdProdi               int                `json:"-"`
	IdDasarDokumen        int                `json:"-"`
	JenisDokumen          string             `json:"jenis_dokumen"`
	NomorDokumen          string             `json:"nomor_dokumen"`
	JenisKerjasama        string             `json:"jenis_kerjasama"`
	DasarDokumenKerjasama DasarKerjasama     `gorm:"-" json:"dasar_dokumen_kerjasama"`
	Judul                 string             `json:"judul"`
	Keterangan            string             `json:"keterangan"`
	Mitra                 []MitraKerjasama   `gorm:"foreignKey:IdKerjasama" json:"mitra"`
	Kegiatan              string             `json:"kegiatan"`
	Status                string             `json:"status"`
	TanggalAwal           string             `json:"tanggal_awal"`
	TanggalBerakhir       string             `json:"tanggal_akhir"`
	Dokumen               string             `json:"dokumen"`
	Prodi                 ProdiReference     `gorm:"foreignKey:IdProdi" json:"prodi"`
	KategoriKegiatan      []KategoriKegiatan `gorm:"many2many:kategori_kegiatan_kerjasama;joinForeignKey:kerjasama_id" json:"kategori_kegiatan"`
}

type MitraKerjasama struct {
	ID                     int `gorm:"primaryKey"`
	IdKerjasama            int
	NamaInstansi           string `gorm:"type:text"`
	NegaraAsal             string `gorm:"type:varchar(255)"`
	Bidang                 string `gorm:"type:varchar(255)"`
	Penandatangan          string `gorm:"type:varchar(255)"`
	JabatanPenandatangan   string `gorm:"type:varchar(255)"`
	PenanggungJawab        string `gorm:"type:varchar(255)"`
	JabatanPenanggungJawab string `gorm:"type:varchar(255)"`
}

func (p *KerjasamaIA) AfterFind(tx *gorm.DB) (err error) {
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

func (KerjasamaIA) TableName() string {
	return "kerjasama"
}
