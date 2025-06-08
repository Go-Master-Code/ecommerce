package models

import (
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

func SaveDetilOrder(db *gorm.DB, idCart int, barangDetilOrder []BarangDetilOrder) {

	/*
		log.Println("Masuk ke Method UpdateCartItems")
		log.Println("Data asli di database")
		log.Println(ci) //struct lama (data di database)

		log.Println("Data update dari form")
		log.Println(updateItemCart) //struct baru
	*/

}

//var brg []BarangDetilOrder //slice untuk memasukkan data

/*
func TambahDetilTransaksi(idTransaksi string, idBarang int, jumlah int) {
	//tambah barang pada slice BarangDetilTransaksi untuk diinput sekaligus
	brg = append(brg, BarangDetilTransaksi{idTransaksi, idBarang, jumlah})
	//pengecekan barang masuk ke dalam slice detil_order
	fmt.Println(brg)
}

func SaveDetilTransaksi(db *gorm.DB) { //save semua row barang ke dalam tabel detil_transaksi
	//insert row(s) ke tabel detil transaksi
	result := db.Table("detil_transaksi").Create(brg)
	if result.Error != nil {
		panic(result.Error)
	}

	//update stok tiap row barang
	barang := Barang{}
	fmt.Println("Kurangi stok tiap barang")

	for i := range brg {
		barang.ID = strconv.Itoa(brg[i].Id)
		log.Println("ID barang : " + barang.ID)

		_ = db.First(&barang, "id = ?", barang.ID) //ambil 1 row dengan ID tertentu

		log.Println("Stok barang : " + strconv.Itoa(barang.Stok))

		barang.Stok = barang.Stok - brg[i].Jumlah //update stok berdasarkan qty terjual

		log.Println("Stok dikurangi : " + strconv.Itoa(brg[i].Jumlah))
		log.Println("Stok barang setelah dipotong : " + strconv.Itoa(barang.Stok))

		_ = db.Save(&barang) //update data ke database
	}

	//kosongkan kembali slice data detil barang
	brg = []BarangDetilTransaksi{}

	fmt.Println("Data master-detil transaksi telah ditambahkan!")
}
*/
