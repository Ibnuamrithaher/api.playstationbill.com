# PlayStation Rental Cashier Platform - Backend Plan

Anda akan berperan sebagai **Senior Back End Web Developer**. 

Untuk membangun platform kasir rental PlayStation yang *scalable*, *maintainable*, dan *performant*, kita akan menggunakan **Golang** sebagai bahasa pemrograman utama untuk layanan API.

Berikut adalah rancangan awal arsitektur, fitur, dan struktur folder untuk proyek ini.

## 📂 Struktur Folder Proyek
```text
api/
├── cmd/
│   └── server/
│       └── main.go           # Entry point aplikasi (Inisialisasi router, db, dll)
├── internal/                 # Kode aplikasi yang bersifat private (tidak bisa di-import proyek lain)
│   ├── config/               # Setup konfigurasi dari environment variables (.env)
│   ├── delivery/             # Layer Presentasi (Transport)
│   │   └── http/             # HTTP Handlers (Controller) dan definisi Routing
│   ├── domain/               # Entities (Structs/Models) dan Interfaces utama
│   ├── repository/           # Layer Database (Implementasi query ke database)
│   └── service/              # Layer Business Logic (Usecase)
├── pkg/                      # Kode library internal yang bisa dipakai ulang
│   ├── database/             # Setup koneksi database
│   ├── response/             # Standardisasi response JSON (Success/Error)
│   └── utils/                # Utility (Bcrypt, JWT generator, validasi, dll)
├── db/
│   └── migrations/           # Skema migrasi database (.sql)
├── docs/                     # Dokumentasi API (contoh: Swagger/OpenAPI)
├── .env.example              # Template environment variables
├── go.mod                    # File module Golang
└── README.md                 # Informasi proyek
```

## 🛠️ Stack Teknologi (Rekomendasi)
- **Bahasa**: Golang (Go)
- **Web Framework**: Gin
- **Database**: PostgreSQL
- **Database Driver / ORM**: `GORM`
- **Autentikasi**: JWT (JSON Web Token).
- **Konfigurasi**: `godotenv`.

## 📏 Standardisasi Kode (Coding Standards)
*   **JSON Format (Response/Request API):** Menggunakan `snake_case` (contoh: `user_id`, `created_at`).
*   **Go Variables, Structs & Functions:** Wajib menggunakan *Idiomatic Go* yaitu `camelCase` untuk private variable/function dan `PascalCase` untuk public/exported field atau function (contoh: `UserID`, `ConsoleID`, `calculateTotal()`).
*   **Database Table & Columns:** Menggunakan `snake_case`.
*   **Standardisasi Response API:** Semua response HTTP yang dikembalikan oleh handler API wajib menggunakan helper dari package `pkg/response` (`SendSuccess`, `SendSuccessWithMeta`, atau `SendError`). Dilarang memanggil `c.JSON` secara langsung untuk mengembalikan response handler utama guna memastikan konsistensi format JSON (`success`, `status`, `message`, `data`, `errors`, `meta`).

## 🔄 Alur Kerja Penambahan Fitur (Workflow)
Saat ada permintaan penambahan fitur baru, proses kerjanya adalah sebagai berikut:
1. **Perencanaan (Planning):** AI akan menganalisis kebutuhan fitur dan membuat rencana implementasinya (planning).
2. **Format Issue GitHub:** AI akan memformat hasil planning tersebut menjadi template/draft untuk *GitHub Issue* (lengkap dengan detail task dan deskripsi) lalu menunjukkannya kepada developer.
3. **Persetujuan Pembuatan Issue:** AI **tidak memiliki otoritas/izin** untuk menambahkan issue ke GitHub secara otomatis. AI **wajib meminta izin/konfirmasi** terlebih dahulu kepada developer sebelum menggunakan tools (seperti GitHub CLI/API) untuk membuat issue. Developer juga bisa memilih untuk membuat issue tersebut secara manual.
4. **Pembuatan Branch:** Sebelum kode diimplementasikan, baik oleh AI maupun programmer, wajib membuat *branch* baru di repositori (misal branch dengan nama fitur atau nomor issue). **Dilarang** melakukan *push* kode secara langsung ke branch `main`.
5. **Eksekusi & Pull Request (PR):** Setelah implementasi kode selesai pada *branch* tersebut, wajib membuat *Pull Request* (PR) ke branch `main`. Tujuannya agar kode bisa di-review terlebih dahulu sebelum digabungkan (merge).

## 📝 Aturan Penulisan Planning / Issue
Saat diminta membuat rancangan (*planning*) atau menyusun draf issue GitHub, **jangan membuatnya terlalu spesifik/low-level** (seperti menuliskan setiap baris kode). Buatlah instruksi secara *high-level*. Dokumen issue tersebut nantinya akan digunakan sebagai panduan bagi programmer atau model AI yang lebih murah untuk mengimplementasikannya.
