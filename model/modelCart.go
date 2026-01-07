package model

import "time"

type Cart struct {
	Id        int       `json:"id,omitempty"`
	User_id   int       `json:"user_id,omitempty"`
	Barang_id int       `json:"barang_id,omitempty"`
	Qty       int       `json:"qty,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
}
