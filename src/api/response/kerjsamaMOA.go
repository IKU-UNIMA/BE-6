package response

import (
	"strings"

	"gorm.io/gorm"
)

type KerjasamaMOA struct {
	ID                    int               `json:"id"`
	IdFakultas            int               `json:"-"`
	IdDasarDokumen        int               `json:"-"`
	JenisDokumen          string            `json:"jenis_dokumen"`
	NomorDokumen          string            `json:"nomor_dokumen"`
	JenisKerjasama        string            `json:"jenis_kerjasama"`
	DasarDokumenKerjasama DasarKerjasama    `gorm:"-" json:"dasar_dokumen_kerjasama"`
	Judul                 string            `json:"judul"`
	Keterangan            string            `json:"keterangan"`
	Mitra                 []MitraKerjasama  `gorm:"foreignKey:IdKerjasama" json:"mitra"`
	Kegiatan              string            `json:"kegiatan"`
	Status                string            `json:"status"`
	TanggalAwal           string            `json:"tanggal_awal"`
	TanggalBerakhir       string            `json:"tanggal_akhir"`
	Dokumen               string            `json:"dokumen"`
	Fakultas              FakultasReference `gorm:"foreignKey:IdFakultas" json:"fakultas"`
}

func (p *KerjasamaMOA) AfterFind(tx *gorm.DB) (err error) {
	p.TanggalAwal = strings.Split(p.TanggalAwal, "T")[0]

	p.TanggalBerakhir = strings.Split(p.TanggalBerakhir, "T")[0]

	return
}

func (KerjasamaMOA) TableName() string {
	return "kerjasama"
}
