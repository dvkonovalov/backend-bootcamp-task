package api

type House struct {
	Id         int    `json:"id"`
	Address    string `json:"address"`
	Developer  string `json:"developer"`
	Year       int    `json:"year"`
	DateCreate string `json:"created_at"`
	UpdateAt   string `json:"update_at"`
}
