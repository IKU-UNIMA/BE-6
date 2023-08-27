package response

import (
	"BE-6/src/util"
	"strings"
	"time"

	"gorm.io/gorm"
)

type KerjasamaMOA struct {
	ID                    int                `json:"id"`
	IdFakultas            int                `json:"-"`
	IdDasarDokumen        int                `json:"-"`
	JenisDokumen          string             `json:"jenis_dokumen"`
	NomorDokumen          string             `json:"nomor_dokumen"`
	JenisKerjasama        string             `json:"jenis_kerjasama"`
	DasarDokumenKerjasama DasarKerjasama     `gorm:"-" json:"dasar_dokumen_kerjasama"`
	Judul                 string             `json:"judul"`
	Keterangan            string             `json:"keterangan"`
	Anggaran              string             `json:"anggaran"`
	SumberPendanaan       string             `json:"sumber_pendanaan"`
	Mitra                 []MitraKerjasama   `gorm:"foreignKey:IdKerjasama" json:"mitra"`
	Status                string             `json:"status"`
	TanggalAwal           string             `json:"tanggal_awal"`
	TanggalBerakhir       string             `json:"tanggal_berakhir"`
	Dokumen               string             `json:"dokumen"`
	Fakultas              FakultasReference  `gorm:"foreignKey:IdFakultas" json:"fakultas"`
	KategoriKegiatan      []KategoriKegiatan `gorm:"foreignKey:IdKerjasama" json:"kategori_kegiatan"`
}

func (p *KerjasamaMOA) AfterFind(tx *gorm.DB) (err error) {
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

func (KerjasamaMOA) TableName() string {
	return "kerjasama"
}
