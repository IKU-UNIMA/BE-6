package model

type Akun struct {
	ID       int    `gorm:"primaryKey"`
	Email    string `gorm:"type:varchar(255);unique"`
	Password string `gorm:"type:varchar(255)"`
	Role     string `gorm:"type:enum('admin','rektor')"`
	Admin    Admin  `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
	Rektor   Rektor `gorm:"foreignKey:ID;constraint:OnDelete:CASCADE"`
}
