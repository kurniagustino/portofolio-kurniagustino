package database

import (
	"fmt"
	"portfolio_kurnia/models"
)

// SeedExperienceIfEmpty memindahkan timeline pengalaman kerja yang sebelumnya
// hardcoded di homepage ke dalam database (hanya bila tabel masih kosong).
func SeedExperienceIfEmpty() {
	var count int64
	if err := DB.Model(&models.Experience{}).Count(&count).Error; err != nil {
		fmt.Println("Cek pengalaman gagal:", err)
		return
	}
	if count > 0 {
		return
	}

	items := []models.Experience{
		{
			Company: "RS Bhayangkara Jambi",
			Role:    "IT Programmer",
			Period:  "Juni 2025 - Sekarang",
			Bullets: "• Membangun dan mengembangkan aplikasi internal rumah sakit menggunakan Laravel, Python, dan Flutter.\n• Maintenance SIMRS dan memastikan kelancaran operasional 24/7.\n• Menjaga konektivitas jaringan (MikroTik) untuk 50 workstation dan troubleshooting Hardware/Software.",
			Order:   1,
		},
		{
			Company: "PT. THRIVENI",
			Role:    "Subcontractor Administrator",
			Period:  "Agustus 2022 - Juni 2023",
			Bullets: "• Mengelola kontrak sewa 50+ unit alat berat, LMV, Bus, dan Tongkang.\n• Validasi dan audit timesheet dari site mining secara berkala untuk compliance.\n• Memproses invoice (Tally) dan monitor evaluasi efisiensi biaya operasional (15%).",
			Order:   2,
		},
		{
			Company: "Program Kotaku PUPR",
			Role:    "Asisten Manajemen Data",
			Period:  "Mei 2020 - Mei 2022",
			Bullets: "• Memimpin survey baseline data 5.000 bangunan hunian di lebih dari 25 kelurahan.\n• Implementasi Geographic Information System (GIS) dan pemutakhiran data 95%.\n• Analisis kelengkapan, akurasi, anomali data, serta pembuatan laporan berkala.",
			Order:   3,
		},
	}

	for _, item := range items {
		if err := DB.Create(&item).Error; err != nil {
			fmt.Println("Seed pengalaman gagal:", err)
		}
	}
	fmt.Printf("Seed pengalaman selesai: %d entri\n", len(items))
}
