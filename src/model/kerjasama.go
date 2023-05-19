package model

import "time"

type Kerjasama struct {
	ID              int       `gorm:"primaryKey"`
	IdProdi         int       `gorm:"default:NULL"`
	IdFakultas      int       `gorm:"default:NULL"`
	Prodi           Prodi     `gorm:"foreignKey:IdProdi;constraint:OnDelete:SET NULL"`
	Fakultas        Fakultas  `gorm:"foreignKey:IdFakultas;constraint:OnDelete:SET NULL"`
	JenisDokumen    string    `gorm:"type:varchar(255)"`
	NomorDokumen    string    `gorm:"type:varchar(255);unique"`
	JenisKerjasama  string    `gorm:"type:varchar(255)"`
	Judul           string    `gorm:"type:text"`
	Keterangan      string    `gorm:"type:text"`
	Mitra           string    `gorm:"type:varchar(255)"`
	Kegiatan        string    `gorm:"type:varchar(255)"`
	Status          string    `gorm:"type:text"`
	TanggalAwal     time.Time `gorm:"type:date"`
	TanggalBerakhir time.Time `gorm:"type:date"`
}
