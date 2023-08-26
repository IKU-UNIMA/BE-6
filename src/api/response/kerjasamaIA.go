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
	TanggalBerakhir       string             `json:"tanggal_berakhir"`
	Dokumen               string             `json:"dokumen"`
	Prodi                 ProdiReference     `gorm:"foreignKey:IdProdi" json:"prodi"`
	KategoriKegiatan      []KategoriKegiatan `gorm:"foreignKey:IdKerjasama" json:"kategori_kegiatan"`
}

type MitraKerjasama struct {
	ID                     int    `json:"id"`
	IdKerjasama            int    `json:"-"`
	NamaInstansi           string `json:"nama_instansi"`
	NegaraAsal             string `json:"negara_asal"`
	Bidang                 string `json:"bidang"`
	Penandatangan          string `json:"penandatangan"`
	JabatanPenandatangan   string `json:"jabatan_penandatangan"`
	PenanggungJawab        string `json:"penanggung_jawab"`
	JabatanPenanggungJawab string `json:"jabatan_penanggung_jawab"`
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
