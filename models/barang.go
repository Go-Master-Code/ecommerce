package models

import (
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        string    `gorm:"primary_key;column:id_cart;autoIncrement"`
	UserID    string    `gorm:"column:user_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"` //gorm tag untuk autocreatetime
	Payment   string    `gorm:"column:payment"`
	Status    string    `gorm:"column:status"`
	//Items     []CartItem `gorm:"foreignKey:id_cart;references:id_cart"` //relasi 1 to many: 1 user menangani beberapa transaksi. Buat foreign key nya, dan buat var User dari struct User serperti ini
	//relasi many to many
	//format: tabel_detil;foreignKey:PK_tabel_ini;joinForeignKey:nama_field_PK_di_tabel_detil;references:PK_tabel_master_lainnya;joinReferences:nama_field_PK_di_tabel_detil
	//Items []Barang `gorm:"many2many:cart_items;foreignKey:id_cart;joinForeignKey:id_barang;references:id_cart;joinReferences:id_cart"`
}

func (o *Order) TableName() string {
	return "order" //nama table pada db
}

// Save Data Master Transaksi
func SaveDataOrder(db *gorm.DB, userID string, payment string) string {
	order := Order{
		UserID:  userID,
		Payment: payment,
		Status:  "Pending",
	}

	//query: INSERT INTO `products` (`id`,`name`,`price`,`created_at`,`updated_at`) VALUES ('P001','Laptop ASUS',10250000,'2024-12-06 15:15:51.069','2024-12-06 15:15:51.069')
	err := db.Create(&order).Error
	if err != nil {
		panic(err)
	}

	idOrder := order.ID //simpan ID transaksi untuk add more detil_jual
	return idOrder      //return value untuk detil transaksi
}

type BarangDetilOrder struct {
	IdOrder   string  `gorm:"column:id_order"`
	IdBarang  int     `gorm:"column:id_barang"`
	Jumlah    float64 `gorm:"column:jumlah"`
	HargaJual float64 `gorm:"column:harga_jual"`
}
