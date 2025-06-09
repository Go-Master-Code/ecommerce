package models

import (
	"log"
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        string    `gorm:"primary_key;column:id_cart;autoIncrement"`
	UserID    string    `gorm:"column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"` //gorm tag untuk autocreatetime
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime;autoUpdateTime"`
	//Items     []CartItem `gorm:"foreignKey:id_cart;references:id_cart"` //relasi 1 to many: 1 user menangani beberapa transaksi. Buat foreign key nya, dan buat var User dari struct User serperti ini
	//relasi many to many
	//format: tabel_detil;foreignKey:PK_tabel_ini;joinForeignKey:nama_field_PK_di_tabel_detil;references:PK_tabel_master_lainnya;joinReferences:nama_field_PK_di_tabel_detil
	//Items []Barang `gorm:"many2many:cart_items;foreignKey:id_cart;joinForeignKey:id_barang;references:id_cart;joinReferences:id_cart"`
}

func (c *Cart) TableName() string {
	return "cart" //nama table pada db
}

func TampilkanCart(db *gorm.DB, idUser string) []Cart { //return value slice [] of Barang
	var cart []Cart

	err := db.Where("user_id=?", idUser).Find(&cart).Error
	if err != nil {
		panic(err)
	}

	return cart
}

// deklarasi struct dan query
type CartItems struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
	IdCart     int     `gorm:"column:id_cart"`
	IdBarang   int     `gorm:"column:id_barang"`
	IconPath   string  `gorm:"column:icon_path"`
	NamaBarang string  `gorm:"column:nama_barang"`
	Jumlah     float64 `gorm:"column:jumlah"`
	Harga      float64 `gorm:"column:harga"`
	Total      float64 `gorm:"column:total"`
}

type CartItemsView struct { //struct detilTrans ini harus didefinisikan setiap field datanya berdasarkan query di bawah, semua field yang dihasilkan harus punya representasi field pada struct, ditambah dengan tag gorm
	IdCart     int     `gorm:"column:id_cart"`
	IdBarang   int     `gorm:"column:id_barang"`
	IconPath   string  `gorm:"column:icon_path"`
	NamaBarang string  `gorm:"column:nama_barang"`
	Jumlah     float64 `gorm:"column:jumlah"`
	Harga      float64 `gorm:"column:harga"`
	Total      float64 `gorm:"column:total"`
	Index      int
}

// buat struct untuk menampung data update
type UpdateItemCart struct {
	IdBarang int
	Jumlah   float64
}

// buat struct untuk menampung data update
type AddItemCart struct {
	IdCart   int
	IdBarang int
	Jumlah   float64
}

func TampilkanCartItems(db *gorm.DB, idUser string) []CartItemsView {
	var ci []CartItems
	err := db.Table("cart_items").Select("cart_items.id_cart, id_barang, icon_path, nama_barang, jumlah, harga, (jumlah*harga) as total").Joins("join barang on cart_items.id_barang=barang.id join cart on cart.id_cart=cart_items.id_cart join user on user.id=cart.user_id").Where("user.id=?", idUser).Find(&ci).Error
	if err != nil {
		panic(err)
	}
	//return ci

	// Tambahkan index secara manual
	var civ []CartItemsView
	for i, item := range ci {
		civ = append(civ, CartItemsView{
			IdCart:     item.IdCart,
			IdBarang:   item.IdBarang,
			NamaBarang: item.NamaBarang,
			Harga:      item.Harga,
			Jumlah:     item.Jumlah,
			Total:      item.Total,
			IconPath:   item.IconPath,
			Index:      i, // Mulai dari nol
		})
	}
	return civ

}

func DeleteItem(db *gorm.DB, idCart string, idBarang string) {
	var ci []CartItems
	//Field deleted at akan terisi, record barang masih ada, tidak dihapus
	err := db.Table("cart_items").Where("id_cart = ? and id_barang =?", idCart, idBarang).Delete(&ci).Error
	if err != nil {
		panic(err)
	}
}

func AddItemToCart(db *gorm.DB, idCart int, idBarang int, jumlah float64) {
	addItemCart := AddItemCart{ //masukkan data (single) pada struct
		IdCart:   idCart,
		IdBarang: idBarang,
		Jumlah:   jumlah,
	}
	//db.Last()
	log.Println("Masuk ke add cart item models")
	err := db.Table("cart_items").Create(&addItemCart).Error
	if err != nil {
		panic(err)
	}
}

func UpdateCartItems(db *gorm.DB, idCart int, updateItemCart []UpdateItemCart) {
	var ci []CartItems
	err := db.Table("cart_items").Select("id_cart, id_barang, nama_barang, jumlah").Joins("join barang on cart_items.id_barang=barang.id").Where("id_cart=?", idCart).Find(&ci).Error
	if err != nil {
		panic(err)
	}

	/*
		log.Println("Masuk ke Method UpdateCartItems")
		log.Println("Data asli di database")
		log.Println(ci) //struct lama (data di database)

		log.Println("Data update dari form")
		log.Println(updateItemCart) //struct baru
	*/

	for _, p := range updateItemCart { //lakukan iterasi terhadap slice updateItemCart dengan var p sebagai recordnya
		db.Model(&CartItems{}).Where("id_cart = ? and id_barang =?", idCart, p.IdBarang).Updates(CartItems{Jumlah: p.Jumlah})
	}

}
