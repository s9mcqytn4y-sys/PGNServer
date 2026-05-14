package utilitas

// ResponsAPI adalah format baku untuk seluruh respons JSON
type ResponsAPI struct {
	Status bool        `json:"status"`
	Pesan  string      `json:"pesan"`
	Data   interface{} `json:"data"`
}

// FormatRespons membungkus data ke dalam format ResponsAPI
func FormatRespons(status bool, pesan string, data interface{}) ResponsAPI {
	return ResponsAPI{
		Status: status,
		Pesan:  pesan,
		Data:   data,
	}
}
