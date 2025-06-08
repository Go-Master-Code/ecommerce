package models

import (
	"gorm.io/gorm"
)

type Kategori struct {
	ID           int      `gorm:"primary_key;column:id;autoIncrement"`
	NamaKategori string   `gorm:"primary_key;column:nama_kategori"`
	Barang       []Barang `gorm:"foreignKey:id;references:id"`
	//relasi many to one terhadap barang
	//ID kategori barang yang sama dapat dimiliki beberapa barang
}

func (k *Kategori) TableName() string {
	return "kategori_barang" //nama table pada db
}

// kategori barang
func TambahKategori(db *gorm.DB, kategoriBarang string) { //untuk memasukkan 1 row (data) ke db
	kategori := Kategori{ //masukkan data (single) pada struct
		NamaKategori: kategoriBarang, //ID tidak didefinisikan karena autoincrement
	}

	err := db.Create(&kategori).Error
	if err != nil {
		panic(err)
	}
}

func ShowKategori(db *gorm.DB) []Kategori { //return value: slice of KategoriBarang
	var kategori []Kategori

	err := db.Find(&kategori).Error
	if err != nil {
		panic(err)
	}

	return kategori
}
