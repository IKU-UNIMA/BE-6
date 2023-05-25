package response

type (
	Dashboard struct {
		Target     string                       `json:"target"`
		Total      int                          `json:"total"`
		TotalProdi int                          `json:"total_prodi"`
		Pencapaian string                       `json:"pencapaian"`
		Detail     []DashboardDetailPerFakultas `json:"detail"`
	}

	DashboardDetailPerFakultas struct {
		ID          int    `json:"id"`
		Fakultas    string `json:"fakultas"`
		JumlahProdi int    `json:"jumlah_prodi"`
		Jumlah      int    `json:"jumlah"`
		Persentase  string `json:"persentase"`
	}

	DashboardPerProdi struct {
		Fakultas   string                    `json:"fakultas"`
		Total      int                       `json:"total"`
		TotalDosen int                       `json:"total_dosen"`
		Pencapaian string                    `json:"pencapaian"`
		Detail     []DashboardDetailPerProdi `json:"detail"`
	}

	DashboardDetailPerProdi struct {
		Prodi   string `json:"prodi"`
		Jumlah  int    `json:"jumlah"`
		Capaian int    `json:"capaian"`
	}
)
