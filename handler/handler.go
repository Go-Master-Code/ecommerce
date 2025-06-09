package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/Go-Master-Code/ecommerce/config"
	"github.com/Go-Master-Code/ecommerce/models"

	"github.com/gorilla/sessions"
)

// inisiasi database
var db = config.OpenConnectionMaster()

// global var sementara untuk idCart dan idUser
var idUser = "budi"
var idCart = 1

func Login(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "login.html")) //buka form untuk input barang
		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

// deklarasi var store untuk session
var store = sessions.NewCookieStore([]byte("secret-key"))

func ValidasiAJAX(w http.ResponseWriter, r *http.Request) {
	type LoginResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	if r.Method == http.MethodPost {
		// Mengambil data dari form login
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")

		//log username n password
		//log.Println(username)
		//log.Println(password)

		// validasi input username dan password
		id, pwd := models.ValidasiUser(db, username, password)

		// Validasi login
		if username == id && password == pwd {
			// Jika login berhasil
			// Login berhasil, buat session
			session, _ := store.Get(r, "session-name")
			session.Values["authenticated"] = true
			session.Save(r, w)

			sessionUser, _ := store.Get(r, "session-name")
			session.Values["username"] = username
			sessionUser.Save(r, w)
			//log.Println(session) bernilai true

			//log.Println("Login sukses")
			response := LoginResponse{
				Success: true,
				Message: "Login berhasil!",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		} else {
			// Jika login gagal
			//log.Println("Login gagal")
			response := LoginResponse{
				Success: false,
				Message: "Username atau password salah.",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(response)
		}
	} else {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

// Fungsi logout untuk menghapus sesi
func Logout(w http.ResponseWriter, r *http.Request) {
	// Menghapus session (logout)
	session, _ := store.Get(r, "session-name")
	session.Values["authenticated"] = false
	session.Save(r, w)

	//log.Println(session) bernilai false
	// Pindah ke halaman login
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//cetak username yang didapat dari session setelah validasiAJAX
	session, _ := store.Get(r, "session-name")

	// Ambil username dari session
	username, _ := session.Values["username"].(string)

	data := struct { //buat struct yang menampung data barang dan kategori barang
		Username string
	}{
		Username: username,
	}

	//handler Form ini hanya bisa menerima request berupa get
	if r.Method == "GET" {
		//jika methodnya get
		tmpl, err := template.ParseFiles(path.Join("views", "main.html"), path.Join("views", "template.html"))

		if err != nil {
			log.Println(err)
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		err = tmpl.Execute(w, data)
		if err != nil {
			log.Println(err.Error())                                           //error spesifik untuk developer
			http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}
		return
	}
	//selain request GET, akan error
	http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError)
}

func ShopHandler(w http.ResponseWriter, r *http.Request) {
	//cetak username yang didapat dari session setelah validasiAJAX
	session, _ := store.Get(r, "session-name")

	// Ambil username dari session
	username, _ := session.Values["username"].(string)

	tmpl, err := template.ParseFiles(path.Join("views", "shop.html"), path.Join("views", "template.html")) //buka file product.html
	if err != nil {
		log.Println(err)
		return
	}

	cart := models.TampilkanCart(db, idUser)
	log.Println("Eksekusi tampilkan cart sukses!")
	barang := models.TampilkanBarang(db)
	log.Println("Eksekusi tampilkan barang sukses!")
	kategori := models.ShowKategori(db)
	log.Println("Eksekusi show kategori sukses!")

	data := struct { //buat struct yang menampung data barang dan kategori barang
		Cart       []models.Cart
		Barangs    []models.Barang
		Kategories []models.Kategori
		Username   string
	}{
		Cart:       cart,
		Barangs:    barang,
		Kategories: kategori,
		Username:   username,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

func CartViewHandler(w http.ResponseWriter, r *http.Request) {
	//cetak username yang didapat dari session setelah validasiAJAX
	session, _ := store.Get(r, "session-name")

	// Ambil username dari session
	username, _ := session.Values["username"].(string)

	tmpl, err := template.ParseFiles(path.Join("views", "cart.html"), path.Join("views", "template.html")) //buka file product.html
	if err != nil {
		log.Println(err)
		return
	}

	//contoh var dibuat global

	cart := models.TampilkanCart(db, username)
	cartItems := models.TampilkanCartItems(db, username)
	//cartItem := models.ShowKategori(db)

	data := struct { //buat struct yang menampung data barang dan kategori barang
		Cart      []models.Cart
		CartItems []models.CartItemsView
		Username  string
		//Kategories []models.Kategori
	}{
		Cart:      cart,
		CartItems: cartItems,
		Username:  username,
		//Kategories: kategori,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

func DeleteCartItems(w http.ResponseWriter, r *http.Request) {
	log.Println("Masuk ke delete item")
	idBarang := r.URL.Query().Get("id")
	idCart := r.URL.Query().Get("cart")

	log.Println("ID barang delete barang dari URL: " + idBarang)
	log.Println("ID cart delete barang dari URL: " + idCart)
	// //ambil data parmeter id dari URL

	models.DeleteItem(db, idCart, idBarang)
	log.Println("Data :" + idBarang + " berhasil dihapus!")
	http.Redirect(w, r, "/cart", http.StatusMovedPermanently)
}

func AddCartItems(w http.ResponseWriter, r *http.Request) {
	log.Println("Masuk handler /add")

	session, _ := store.Get(r, "session-name")

	// Ambil username dari session
	username, _ := session.Values["username"].(string)

	var idBarang string = r.URL.Query().Get("id")

	log.Println("ID barang: " + idBarang)

	cart := models.TampilkanCart(db, username)
	idCart := cart[0].ID
	log.Println("Cart ID: " + cart[0].ID)

	idBarangInt, _ := strconv.Atoi(idBarang)
	idCartInt, _ := strconv.Atoi(idCart)

	models.AddItemToCart(db, idCartInt, idBarangInt, 1)
	http.Redirect(w, r, "/shop", http.StatusMovedPermanently)
}

func UpdateCartItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("Masuk method update data barang")
		//ambil data dari form, lakukan post sesuai tipe form : POST
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())                                        //error spesifik untuk developer
			http.Error(w, "Ini bukan POST", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		//log.Println("Parse Form berhasil!")
		//Dapatkan jumlah iterasi
		itemCount := r.Form.Get("itemCount")
		itemCountInt, _ := strconv.Atoi(itemCount)

		var updateItemCart []models.UpdateItemCart

		var iterasi string
		var idBarang string
		var jumlah string

		var idBarangInt int
		var jumlahFloat float64

		for i := 0; i < itemCountInt; i++ {
			iterasi = strconv.Itoa(i)
			idBarang = r.Form.Get("id" + iterasi)
			jumlah = r.Form.Get("jumlah" + iterasi)

			idBarangInt, _ = strconv.Atoi(idBarang)
			jumlahFloat, _ = strconv.ParseFloat(jumlah, 32)
			updateItemCart = append(updateItemCart, models.UpdateItemCart{
				IdBarang: idBarangInt,
				Jumlah:   jumlahFloat})
		}

		log.Println(updateItemCart)
		models.UpdateCartItems(db, idCart, updateItemCart)
		//setelah selesai update cart, pindah ke halaman checkout
		http.Redirect(w, r, "/checkout", http.StatusMovedPermanently)
	}
}

func CheckoutHandler(w http.ResponseWriter, r *http.Request) { //parameter handler func wajib seperti ini
	//update data cart
	//cetak username yang didapat dari session setelah validasiAJAX
	session, _ := store.Get(r, "session-name")

	// Ambil username dari session
	username, _ := session.Values["username"].(string)

	tmpl, err := template.ParseFiles(path.Join("views", "checkout.html"), path.Join("views", "template.html")) //buka file product.html
	if err != nil {
		log.Println(err)
		return
	}

	user := models.TampilkanUser(db, idUser)
	cart := models.TampilkanCart(db, idUser)
	cartItems := models.TampilkanCartItems(db, idUser)
	//cartItem := models.ShowKategori(db)

	//hitung total harga
	var totalCheckout float64
	for _, item := range cartItems {
		totalCheckout = totalCheckout + (item.Harga * item.Jumlah)
	}

	//formatted := fmt.Sprintf("%.2f", totalCheckout)
	//formattedFloat, _ := strconv.ParseFloat(formatted, 64)

	log.Println("Total checkout: ", totalCheckout)

	data := struct { //buat struct yang menampung data barang dan kategori barang
		User               []models.User
		Cart               []models.Cart
		CartItems          []models.CartItemsView
		TotalCheckOut      float64
		TotalCheckOutFinal float64
		Username           string
		//Kategories []models.Kategori
	}{
		User:               user,
		Cart:               cart,
		CartItems:          cartItems,
		TotalCheckOut:      totalCheckout,
		TotalCheckOutFinal: totalCheckout + 3,
		Username:           username,
		//Kategories: kategori,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

// func SaveOrder(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		log.Println("Masuk method save checkout data")
// 		//ambil data dari form, lakukan post sesuai tipe form : POST
// 		err := r.ParseForm()
// 		if err != nil {
// 			log.Println(err.Error())                                        //error spesifik untuk developer
// 			http.Error(w, "Ini bukan POST", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
// 			return
// 		}

// 		idOrder := models.SaveDataOrder(db, "budi")
// 		log.Println("ID Order: " + idOrder + " telah tersimpan!")

// 		/*
// 			var kodeRow string = ""
// 			var qtyRow string = ""
// 			var hargaRow string = ""

// 			log.Println("Jumlah baris :" + r.Form.Get("jmlBaris"))
// 			//Ambil value jumlah iterasi dari form web (text box dengan id=jml)
// 			iterasi, _ := strconv.Atoi(r.Form.Get("jmlBaris"))
// 		*/
// 	}
// }

func GetOrderItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		payment := r.FormValue("payment")
		idOrder := models.SaveDataOrder(db, "budi", payment)

		log.Println("ID Order tersimpan: " + idOrder + " telah tersimpan!")

		log.Println("Masuk method get order items")
		//ambil data dari form, lakukan post sesuai tipe form : POST
		err := r.ParseForm()
		if err != nil {
			log.Println(err.Error())                                        //error spesifik untuk developer
			http.Error(w, "Ini bukan POST", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
			return
		}

		//log.Println("Parse Form berhasil!")
		//Dapatkan jumlah iterasi
		itemCount := r.Form.Get("itemCount")
		itemCountInt, _ := strconv.Atoi(itemCount)
		log.Println("Item count: " + itemCount)

		var bdo []models.BarangDetilOrder

		var iterasi string
		var idBarang string
		var jumlah string
		var hargajual string

		var idBarangInt int
		var jumlahFloat float64
		var hargajualFloat float64

		for i := 0; i < itemCountInt; i++ {
			iterasi = strconv.Itoa(i)
			idBarang = r.Form.Get("id" + iterasi)
			jumlah = r.Form.Get("jumlah" + iterasi)
			hargajual = r.Form.Get("harga" + iterasi)

			idBarangInt, _ = strconv.Atoi(idBarang)
			jumlahFloat, _ = strconv.ParseFloat(jumlah, 64)
			hargajualFloat, _ = strconv.ParseFloat(hargajual, 64)

			bdo = append(bdo, models.BarangDetilOrder{
				IdOrder:   idOrder,
				IdBarang:  idBarangInt,
				Jumlah:    jumlahFloat,
				HargaJual: hargajualFloat,
			})
		}

		log.Println(bdo)

		//Save data detil order

		result := db.Table("order_items").Create(bdo)
		if result.Error != nil {
			panic(result.Error)
		}

		http.Redirect(w, r, "/shop", http.StatusMovedPermanently)
		/*
			models.UpdateCartItems(db, idCart, updateItemCart)
			//setelah selesai update cart, pindah ke halaman checkout
			http.Redirect(w, r, "/checkout", http.StatusMovedPermanently)
		*/
	}
}
