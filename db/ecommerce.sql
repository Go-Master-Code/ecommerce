-- --------------------------------------------------------
-- Host:                         127.0.0.1
-- Server version:               5.7.44-log - MySQL Community Server (GPL)
-- Server OS:                    Win64
-- HeidiSQL Version:             12.10.0.7000
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

-- Dumping data for table ecommerce.barang: ~6 rows (approximately)
INSERT INTO `barang` (`id`, `nama_barang`, `harga`, `stok`, `id_kategori`, `deskripsi`, `image_path`, `icon_path`, `created_at`, `updated_at`, `deleted_at`) VALUES
	(1, 'Grapes', 3.2, 10, 1, 'Red Globe from USA', 'static/product_image/Grapes.jpg', 'static/product_icon/GrapesIcon.jpg', '2025-06-04 09:41:05.000', NULL, NULL),
	(2, 'Raspberries', 4.1, 36, 1, 'Raspberries from Rusia', 'static/product_image/Raspberries.jpg', 'static/product_icon/RaspberriesIcon.jpg', '2025-06-04 09:41:50.000', NULL, NULL),
	(3, 'Apricots', 3.18, 5, 1, 'Apricots from Italia', 'static/product_image/Apricots.jpg', 'static/product_icon/ApricotsIcon.jpg', '2025-06-04 09:42:06.000', NULL, NULL),
	(4, 'Oranges', 2.35, 14, 1, 'Oranges from Australia', 'static/product_image/Oranges.jpg', 'static/product_icon/OrangesIcon.jpg', '2025-06-04 09:43:52.000', NULL, NULL),
	(5, 'Apple', 1.95, 27, 1, 'Apple from China', 'static/product_image/Apple.jpg', 'static/product_icon/AppleIcon.jpg', '2025-06-04 09:43:53.000', NULL, NULL),
	(6, 'Banana', 2.2, 11, 1, 'Banana from Phipines', 'static/product_image/Banana.jpg', 'static/product_icon/BananaIcon.jpg', '2025-06-04 09:44:43.000', NULL, NULL);

-- Dumping data for table ecommerce.cart: ~2 rows (approximately)
INSERT INTO `cart` (`id_cart`, `user_id`, `created_at`, `updated_at`) VALUES
	(1, 'budi', '2025-06-05 11:52:18', '2025-06-05 18:52:18'),
	(2, 'rini', '2025-06-08 02:27:30', '2025-06-08 09:27:31');

-- Dumping data for table ecommerce.cart_items: ~5 rows (approximately)
INSERT INTO `cart_items` (`id_cart`, `id_barang`, `jumlah`) VALUES
	(1, 2, 3),
	(1, 3, 5),
	(2, 1, 1),
	(2, 4, 3),
	(2, 6, 4);

-- Dumping data for table ecommerce.kategori_barang: ~2 rows (approximately)
INSERT INTO `kategori_barang` (`id`, `nama_kategori`) VALUES
	(1, 'Fruits'),
	(2, 'Vegetables');

-- Dumping data for table ecommerce.order: ~3 rows (approximately)
INSERT INTO `order` (`id_order`, `user_id`, `created_at`, `payment`, `status`) VALUES
	(7, 'budi', '2025-06-08 01:41:25', 'transfer', 'Pending'),
	(8, 'budi', '2025-06-08 02:25:23', 'COD', 'Pending'),
	(9, 'budi', '2025-06-08 02:26:55', 'COD', 'Pending');

-- Dumping data for table ecommerce.order_items: ~6 rows (approximately)
INSERT INTO `order_items` (`id_order`, `id_barang`, `jumlah`, `harga_jual`) VALUES
	(7, 2, 4, 4.1),
	(7, 3, 2, 3.18),
	(8, 2, 9, 4.1),
	(8, 3, 1, 3.18),
	(9, 2, 3, 4.1),
	(9, 3, 5, 3.18);

-- Dumping data for table ecommerce.user: ~2 rows (approximately)
INSERT INTO `user` (`id`, `email`, `password`, `id_cart`, `first_name`, `last_name`, `telepon`, `negara`, `alamat`, `kota`, `kode_pos`) VALUES
	('budi', 'budijatmiko@gmail.com', 'budi123', 1, 'Budianto', 'Jatmiko', '081392837384', 'Indonesia', 'Jl. Cemara No. 13', 'Bandung', '40216'),
	('rini', 'riniastuti@gmail.com', 'rini123', 2, 'Rini', 'Astuti', '087827162736', 'Indonesia', 'Jl. Taman Sari No. 11', 'Bandung', '40227');

/*!40103 SET TIME_ZONE=IFNULL(@OLD_TIME_ZONE, 'system') */;
/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IFNULL(@OLD_FOREIGN_KEY_CHECKS, 1) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40111 SET SQL_NOTES=IFNULL(@OLD_SQL_NOTES, 1) */;
