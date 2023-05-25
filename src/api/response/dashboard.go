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
