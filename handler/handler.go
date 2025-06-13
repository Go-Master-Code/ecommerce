package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/Go-Master-Code/ecommerce/config"
	"github.com/Go-Master-Code/ecommerce/models"

	"github.com/gorilla/sessions"
)

// inisiasi database
var db = config.OpenConnectionMaster()

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

	barang := models.TampilkanBarang(db)
	kategori := models.ShowKategori(db)
	cart := models.TampilkanCart(db, username)
	log.Println("Eksekusi tampilkan shop sukses!")

	data := struct { //buat struct yang menampung data barang dan kategori barang
		Barangs    []models.Barang
		Kategories []models.Kategori
		Cart       []models.Cart
		Username   string
		NowUnix    int64 //Tambahkan timestamp ke data template
	}{
		Barangs:    barang,
		Kategories: kategori,
		Cart:       cart,
		Username:   username,
		NowUnix:    time.Now().Unix(), //Query &t=unix_timestamp bikin URL selalu beda tiap load halaman → browser tidak cache request.
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

	cart := models.TampilkanCart(db, username)
	cartItems := models.TampilkanCartItems(db, username)
	//cartItem := models.ShowKategori(db)

	data := struct { //buat struct yang menampung data barang dan kategori barang
		Cart      []models.Cart
		CartItems []models.CartItemsView
		Username  string
		NowUnix   int64 //Tambahkan timestamp ke data template
		//Kategories []models.Kategori
	}{
		Cart:      cart,
		CartItems: cartItems,
		Username:  username,
		//Kategories: kategori,
		NowUnix: time.Now().Unix(), //Query &t=unix_timestamp bikin URL selalu beda tiap load halaman → browser tidak cache request.
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())                                           //error spesifik untuk developer
		http.Error(w, "Terjadi kesalahan", http.StatusInternalServerError) //Error generik untuk user, pakai bahasa manusia
		return
	}
}

// OLD METHHOD -> Add Item to cart masih reload halaman
func AddCartItems(w http.ResponseWriter, r *http.Request) {
	//sudah pakai middleware untuk disable cache

	//log.Println("Masuk handler /add")

	//log.Println("URL Query:", r.URL.RawQuery)
	idBarang := r.URL.Query().Get("id_brg_shop")

	//log.Println("ID barang: " + idBarang)

	session, _ := store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)

	cart := models.TampilkanCart(db, username)
	idCart := cart[0].ID
	//log.Println("Cart ID: " + idCart)

	idBarangInt, _ := strconv.Atoi(idBarang)
	idCartInt, _ := strconv.Atoi(idCart)

	models.AddItemToCart(db, idCartInt, idBarangInt, 1)
	//log.Println("Barang: " + idBarang + "telah ditambahkan!")
	http.Redirect(w, r, "/shop", http.StatusSeeOther)
}

// eksperimental -> Add item to cart no reload
func AddBarangToCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Ambil data dari form-encoded
	idBarang := r.FormValue("product_id")
	//jumlah := r.FormValue("quantity")

	//fmt.Printf("Barang ditambahkan: %s (jumlah: %s)\n", idBarang, jumlah)

	w.Header().Set("Content-Type", "text/plain")
	//fmt.Fprintf(w, "Berhasil menambahkan %s (jumlah: %s) ke keranjang.", idBarang, jumlah)

	session, _ := store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)

	cart := models.TampilkanCart(db, username)
	idCart := cart[0].ID

	idBarangInt, _ := strconv.Atoi(idBarang)
	idCartInt, _ := strconv.Atoi(idCart)

	models.AddItemToCart(db, idCartInt, idBarangInt, 1)
}

// OLD METHOD DELETE ITEM FROM CART (Reload page)
func DeleteCartItems(w http.ResponseWriter, r *http.Request) {
	//sudah pakai middleware untuk disable cache

	log.Println("Masuk ke delete item")
	idBarang := r.URL.Query().Get("id_brg_delete")
	idCart := r.URL.Query().Get("id_cart_delete")

	log.Println("ID barang delete barang dari URL: " + idBarang)
	log.Println("ID cart delete barang dari URL: " + idCart)
	// //ambil data parmeter id dari URL

	idBarangInt, _ := strconv.Atoi(idBarang)

	models.DeleteItem(db, idCart, idBarangInt)
	log.Println("Data :" + idBarang + " berhasil dihapus!")
	http.Redirect(w, r, "/cart", http.StatusSeeOther)
}

// PERLU DISELIDIKI LAGI AGAR DATA BERUBAH TAPI GA RELOAD
func DeleteBarangFromCart(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idBarang := r.FormValue("product_id")

	w.Header().Set("Content-Type", "text/plain")

	session, _ := store.Get(r, "session-name")
	username, _ := session.Values["username"].(string)

	cart := models.TampilkanCart(db, username)
	idCart := cart[0].ID

	idBarangInt, _ := strconv.Atoi(idBarang)
	//idCartInt, _ := strconv.Atoi(idCart)

	models.DeleteItem(db, idCart, idBarangInt)
	log.Println("Data :" + idBarang + " dari cart: " + idCart + " berhasil dihapus!")
}

func UpdateCartItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//cetak username yang didapat dari session setelah validasiAJAX
		session, _ := store.Get(r, "session-name")
		username, _ := session.Values["username"].(string)

		cart := models.TampilkanCart(db, username)
		idCart := cart[0].ID
		idCartInt, _ := strconv.Atoi(idCart)

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
		models.UpdateCartItems(db, idCartInt, updateItemCart)
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

	user := models.TampilkanUser(db, username)
	cart := models.TampilkanCart(db, username)
	cartItems := models.TampilkanCartItems(db, username)
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

func GetOrderItems(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		//ambil session username
		session, _ := store.Get(r, "session-name")
		username, _ := session.Values["username"].(string)

		//ambil data payment method dari form
		payment := r.FormValue("payment")
		//Method SaveDataOrder secara default mengambalikan nilai primary key dalam hal ini id_order
		idOrder := models.SaveDataOrder(db, username, payment)

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

		//Update data barang di db
		models.UpdateStokBarangDetilOrder(db, bdo)

		//ambil id cart berdasarkan user id
		cart := models.TampilkanCart(db, username)
		models.ClearCart(db, cart[0].ID) //clear cart yang sudah di checkout

		//setelah selesai update stok barang dan clear cart, pindah ke halaman shop
		http.Redirect(w, r, "/shop", http.StatusSeeOther)

	}
}
