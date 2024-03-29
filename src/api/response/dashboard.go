package response

type (
	Dashboard struct {
		Mou         int                          `json:"mou"`
		Moa         int                          `json:"moa"`
		Ia          int                          `json:"ia"`
		LuarNegeri  int                          `json:"luar_negeri"`
		DalamNegeri int                          `json:"dalam_negeri"`
		Target      string                       `json:"target"`
		Total       int                          `json:"total"`
		TotalProdi  int                          `json:"total_prodi"`
		Pencapaian  string                       `json:"pencapaian"`
		Detail      []DashboardDetailPerFakultas `gorm:"-" json:"detail"`
	}

	DashboardDetailPerFakultas struct {
		ID               int    `json:"id"`
		Fakultas         string `json:"fakultas"`
		JumlahProdi      int    `json:"jumlah_prodi"`
		JumlahPencapaian int    `json:"jumlah_pencapaian"`
		Persentase       string `json:"persentase"`
	}

	DashboardPerProdi struct {
		Fakultas   string                    `json:"fakultas"`
		Total      int                       `json:"total"`
		TotalProdi int                       `json:"total_prodi"`
		Pencapaian string                    `json:"pencapaian"`
		Detail     []DashboardDetailPerProdi `json:"detail"`
	}

	DashboardDetailPerProdi struct {
		Prodi           string `json:"prodi"`
		JumlahKerjasama int    `json:"jumlah_kerjasama"`
		Capaian         int    `json:"capaian"`
	}
)
