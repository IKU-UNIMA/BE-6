package model

type Prodi struct {
	ID         int `gorm:"primaryKey"`
	IdFakultas int
	KodeProdi  int      `gorm:"unique"`
	Nama       string   `gorm:"type:varchar(255)"`
	Jenjang    string   `gorm:"type:varchar(60)"`
	Fakultas   Fakultas `gorm:"foreignKey:IdFakultas;constraint:OnDelete:CASCADE"`
}
