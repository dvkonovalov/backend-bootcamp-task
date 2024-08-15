package api

type Flat struct {
	Id       int    `json:"id"`
	House_id int    `json:"house_id"`
	Price    int    `json:"price"`
	Rooms    int    `json:"rooms"`
	Status   string `json:"status"`
}
