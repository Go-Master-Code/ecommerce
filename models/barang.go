package models

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Barang struct {
	ID         string         `gorm:"primary_key;column:id;autoIncrement"`
	IdKategori int            `gorm:"column:id_kategori"`
	NamaBarang string         `gorm:"column:nama_barang"`
	Harga      float64        `gorm:"column:harga"`
	Stok       int            `gorm:"column:stok"`
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

func TampilkanBarangPerKategori(db *gorm.DB, idKategori int) {
	var barang []Barang
	/*
		2 query:
		SELECT * FROM `addresses` WHERE `addresses`.`user_id` IN ('1','20','50','10','11','12','13','14','2','21','3','4','5','6','7','8','9')
		SELECT `users`.`id`,`users`.`password`,`users`.`first_name`,`users`.`middle_name`,`users`.`last_name`,`users`.`created_at`,`users`.`updated_at`,`Wallet`.`id` AS `Wallet__id`,`Wallet`.`user_id` AS `Wallet__user_id`,`Wallet`.`balance` AS `Wallet__balance`,`Wallet`.`created_at` AS `Wallet__created_at`,`Wallet`.`updated_at` AS `Wallet__updated_at` FROM `users` LEFT JOIN `wallets` `Wallet` ON `users`.`id` = `Wallet`.`user_id`
	*/
	err := db.Model(&Barang{}).Preload("Kategori").Where("barang.id_kategori = ?", idKategori).Find(&barang).Error
	if err != nil {
		panic(err)
	}

	fmt.Println("==========================Data Stok barang==========================")
	fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", "ID", "Nama Barang", "Harga", "Stok", "Kategori")
	//fmt.Println("ID | Nama Barang | Harga | Stok | Kategori |")
	fmt.Println("====================================================================")

	for i := range barang {
		fmt.Printf("%-5s %-30s %-10s %-8s %-12s\n", barang[i].ID, barang[i].NamaBarang, strconv.FormatFloat(float64(barang[i].Harga), 'f', -1, 64), strconv.Itoa(barang[i].Stok), barang[i].Kategori.NamaKategori)
		//fmt.Println(barang[i].ID + " | " + barang[i].NamaBarang + " | " + models.FormatAngka(barang[i].Harga) + " | " + strconv.Itoa(barang[i].Stok) + " | " + barang[i].KategoriBarang.NamaKategori + " | ")
	}
}

/* UNCOMMENT UNTUK UPDATE BARANG DETIL TRANSAKSI
var brg []BarangDetilTransaksi //slice untuk memasukkan data
// Update stok barang yang ada pada detil transaksi
func UpdateStokBarangDetilTransaksi(db *gorm.DB) { //save semua row barang ke dalam tabel detil_transaksi
	fr := Barang{}
	//fmt.Println("Kurangi stok tiap barang")

	for i := range brg {
		barang.ID = strconv.Itoa(brg[i].Id)
		log.Println("ID barang pertama: " + barang.ID)

		_ = db.First(&barang, "id = ?", barang.ID) //ambil 1 row dengan ID tertentu

		log.Println("Stok barang : " + strconv.Itoa(barang.Stok))

		barang.Stok = barang.Stok - brg[i].Jumlah //update stok berdasarkan qty terjual

		log.Println("Stok dikurangi : " + strconv.Itoa(brg[i].Jumlah))
		log.Println("Stok barang setelah dipotong : " + strconv.Itoa(barang.Stok))

		_ = db.Save(&barang) //update data ke database
		//fmt.Println("Stok barang " + barang.ID + " telah diupdate menjadi: " + strconv.Itoa(barang.Stok))
	}
}
*/

func TampilkanBarangSedikit(db *gorm.DB) []Barang {
	var barang []Barang

	//tampilkan barang yang stoknya < 10
	result := db.Model(Barang{}).Preload("Kategori").Where("stok < ?", "10").Find(&barang).Error
	if result != nil {
		panic(result)
	}

	return barang
}
