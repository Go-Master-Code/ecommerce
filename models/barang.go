package models

import (
	"time"

	"gorm.io/gorm"
)

type Barang struct {
	ID         int            `gorm:"primary_key;column:id;autoIncrement"`
	IdKategori int            `gorm:"column:id_kategori"`
	NamaBarang string         `gorm:"column:nama_barang"`
	Harga      float64        `gorm:"column:harga"`
	Stok       float64        `gorm:"column:stok"`
	Deskripsi  string         `gorm:"column:deskripsi"`
	ImagePath  string         `gorm:"column:image_path"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"` //tipe datanya bukan time.Time tapi gorm.DeletedAt -> penanda soft delete
	Kategori   Kategori       `gorm:"foreignKey:id_kategori;references:id"`
	//ItemsInCart []Cart         `gorm:"many2many:cart_items;foreignKey:id;joinForeignKey:id_barang;references:id_cart;joinReferences:id_cart"`

	//JualBarang []Transaksi    `gorm:"many2many:detil_transaksi;foreignKey:id;joinForeignKey:id_barang;references:id_transaksi;joinReferences:id_transaksi"`
	//format: tabel_many_to_many;foreignKey:PK_tabel_ini;joinForeignKey:nama_field_PK_di_tabel_detil;references:PK_tabel_master_lainnya;joinReferences:nama_field_PK_di_tabel_detil
}

func (b *Barang) TableName() string {
	return "barang" //nama table pada db nya adalah user_logs
}

func TampilkanBarang(db *gorm.DB) []Barang { //return value slice [] of Barang
	var barang []Barang

	err := db.Model(&Barang{}).Preload("Kategori").Find(&barang).Error
	if err != nil {
		panic(err)
	}

	return barang
}

func TampilkanBarangOrderByNama(db *gorm.DB) []Barang { //return value slice [] of Barang
	var barang []Barang

	err := db.Model(&Barang{}).Preload("Kategori").Order("nama_barang asc").Find(&barang).Error
	if err != nil {
		panic(err)
	}

	return barang
}

func TampilkanBarangById(db *gorm.DB, idBarang string) ([]Barang, int) { //return value slice [] of Barang dan selected idKategori untuk dropdown
	var barang []Barang

	err := db.Model(&Barang{}).Preload("Kategori").Where("id = ?", idBarang).Find(&barang).Error
	if err != nil {
		panic(err)
	}

	idKategori := barang[0].IdKategori //mengambil idKategori dari produk terpilih untuk dijadikan default value pada dropdown

	return barang, idKategori
}

// Update stok barang yang ada pada detil transaksi
func UpdateStokBarangDetilOrder(db *gorm.DB, bdo []BarangDetilOrder) {
	barang := Barang{}

	//var idBarang int

	for i := range bdo {
		barang.ID = bdo[i].IdBarang
		log.Println("ID barang ke : ", i+1, " :", barang.ID)

		_ = db.Find(&barang, "id = ?", barang.ID) //ambil 1 row dengan ID tertentu

		log.Println("Stok barang : ", barang.Stok)

		barang.Stok = barang.Stok - bdo[i].Jumlah //update stok berdasarkan qty terjual

		log.Println("Stok dikurangi : ", bdo[i].Jumlah)
		log.Println("Stok barang setelah dipotong : ", barang.Stok)

		_ = db.Save(&barang) //update data ke database
		log.Println("Stok barang ", barang.ID, " telah diupdate menjadi: ", barang.Stok)
	}
}

func TampilkanBarangSedikit(db *gorm.DB) []Barang {
	var barang []Barang

	//tampilkan barang yang stoknya < 10
	result := db.Model(Barang{}).Preload("Kategori").Where("stok < ?", "10").Find(&barang).Error
	if result != nil {
		panic(result)
	}

	return barang
}
