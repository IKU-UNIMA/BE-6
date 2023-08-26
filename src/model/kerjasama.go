package model

import "time"

type Kerjasama struct {
	ID               int              `gorm:"primaryKey"`
	IdProdi          int              `gorm:"default:NULL"`
	IdFakultas       int              `gorm:"default:NULL"`
	IdDasarDokumen   int              `gorm:"default:NULL"`
	Prodi            Prodi            `gorm:"foreignKey:IdProdi;constraint:OnDelete:SET NULL"`
	Fakultas         Fakultas         `gorm:"foreignKey:IdFakultas;constraint:OnDelete:SET NULL"`
	JenisDokumen     string           `gorm:"type:varchar(255)"`
	NomorDokumen     string           `gorm:"type:varchar(255);unique"`
	JenisKerjasama   string           `gorm:"type:varchar(255)"`
	Judul            string           `gorm:"type:text"`
	Keterangan       string           `gorm:"type:text"`
	Mitra            []MitraKerjasama `gorm:"foreignKey:IdKerjasama;constraint:OnDelete:CASCADE"`
	Kegiatan         string           `gorm:"type:varchar(255)"`
	Dokumen          string
	TanggalAwal      time.Time          `gorm:"type:date"`
	TanggalBerakhir  time.Time          `gorm:"type:date"`
	KategoriKegiatan []KategoriKegiatan `gorm:"foreignKey:IdKerjasama;constraint:OnDelete:CASCADE"`
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
