package models

import (
	"gorm.io/gorm"
)

type User struct {
	ID        string `gorm:"primary_key;column:id;autoIncrement"`
	Password  string `gorm:"column:password"`
	IdCart    int    `gorm:"column:id_cart"`
	Email     string `gorm:"column:email"`
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Telepon   string `gorm:"column:telepon"`
	Negara    string `gorm:"column:negara"`
	Alamat    string `gorm:"column:alamat"`
	Kota      string `gorm:"column:kota"`
	KodePos   string `gorm:"column:kode_pos"`
}

func (u *User) TableName() string {
	return "user" //nama table pada db nya adalah user_logs
}

func TampilkanUser(db *gorm.DB, idUser string) []User { //return value slice [] of Barang
	var user []User

	err := db.Model(&User{}).Where("id=?", idUser).Find(&user).Error
	if err != nil {
		panic(err)
	}

	return user
}

func ValidasiUser(db *gorm.DB, idUser string, password string) (string, string) { //returnkan username dan password berdasarkan login
	var user []User

	err := db.Model(&User{}).Where("id = ? and password =?", idUser, password).Find(&user).Error
	if err != nil {
		panic(err)
	}

	if len(user) > 0 {
		//jika data user ketemu
		idUser := user[0].ID
		pwd := user[0].Password

		return idUser, pwd
	} else {
		//jika data user tidak ketemu
		return "", ""
	}
}
