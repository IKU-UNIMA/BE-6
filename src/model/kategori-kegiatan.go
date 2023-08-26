package model

type KategoriKegiatan struct {
	ID                      int `gorm:"primaryKey"`
	IdKerjasama             int
	IdJenisKategoriKegiatan int
	NilaiKontrak            float64
	Volume                  string `gorm:"type:varchar(255)"`
	At                      string `gorm:"type:varchar(255)"`
	Keterangan              string
	Sasaran                 string                `gorm:"type:varchar(255)"`
	IndikatorKinerja        string                `gorm:"type:varchar(255)"`
	JenisKategoriKegiatan   JenisKategoriKegiatan `gorm:"foreignKey:IdJenisKategoriKegiatan;constraint:OnDelete:SET NULL"`
}
