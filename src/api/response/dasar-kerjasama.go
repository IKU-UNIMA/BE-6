package response

type DasarKerjasama struct {
	ID           int    `json:"id"`
	NomorDokumen string `json:"nomor_dokumen"`
	Judul        string `json:"judul"`
	JenisDokumen string `json:"jenis_dokumen"`
}

func (DasarKerjasama) TableName() string {
	return "kerjasama"
}
