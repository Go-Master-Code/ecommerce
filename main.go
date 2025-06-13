package main

import (
	"log"
	"net/http"

	"github.com/Go-Master-Code/ecommerce/handler"
	"github.com/Go-Master-Code/ecommerce/middleware"
)

func main() {
	mux := http.NewServeMux() //buat server mux
	mux.HandleFunc("/login", handler.Login)
	mux.HandleFunc("/logout", handler.Logout)
	mux.HandleFunc("/ajax", handler.ValidasiAJAX)
	//mux.HandleFunc("/home", handler.HomeHandler)
	mux.HandleFunc("/shop", handler.ShopHandler)
	mux.HandleFunc("/cart", handler.CartViewHandler)
	mux.HandleFunc("/delete", middleware.NoCache(handler.DeleteCartItems)) //handler untuk delete barang
	//mux.HandleFunc("/add", middleware.NoCache(handler.AddCartItems))               ////OLD METHHOD -> Add Item to cart masih reload halaman
	mux.HandleFunc("/add", middleware.NoCache(handler.AddBarangToCart))           //handler untuk add barang ke cart
	mux.HandleFunc("/hapus", middleware.NoCache(handler.DeleteBarangFromCart))    //handler untuk hapus barang dari cart
	mux.HandleFunc("/update", handler.UpdateCartItems)                            //handler untuk update items pada cart
	mux.HandleFunc("/checkout", handler.CheckoutHandler)                          //handler untuk checkout (data cart final)
	mux.HandleFunc("/get_order_items", middleware.NoCache(handler.GetOrderItems)) //handler untuk save data master order

	fileServer := http.FileServer(http.Dir("assets")) //load directory
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	//buat log untuk tahu server sudah berjalan di port 3000
	log.Println("Server running on port 8080")

	//menjalankan server
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err) //jika terjadi error
}
