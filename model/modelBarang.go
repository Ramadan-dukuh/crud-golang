package model

import "time"

type Barang struct {
	Id          int       `json:"id,omitempty"`
	Nama_barang string    `json:"nama_barang,omitempty"`
	Deskripsi   string    `json:"deskripsi,omitempty"`
	Kategori    string    `json:"kategori,omitempty"`
	Harga       int       `json:"harga,omitempty"`
	Stok        int       `json:"stok,omitempty"`
	Gambar      string    `json:"gambar,omitempty"`
	Created_at  time.Time `json:"created_at,omitempty"`
	Updated_at  time.Time `json:"updated_at,omitempty"`
}
