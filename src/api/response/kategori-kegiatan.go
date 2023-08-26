package response

type KategoriKegiatan struct {
	ID                      int                   `json:"id"`
	IdKerjasama             int                   `json:"-"`
	IdJenisKategoriKegiatan int                   `json:"-"`
	NilaiKontrak            float64               `json:"nilai_kontrak"`
	Volume                  string                `json:"volume"`
	At                      string                `json:"at"`
	Keterangan              string                `json:"keterangan"`
	Sasaran                 string                `json:"sasaran"`
	IndikatorKinerja        string                `json:"indikator_kinerja"`
	JenisKategoriKegiatan   JenisKategoriKegiatan `gorm:"foreignKey:IdJenisKategoriKegiatan" json:"jenis_kategori_kegiatan"`
}
