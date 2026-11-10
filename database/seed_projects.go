package database

import (
	"fmt"
	"portfolio_kurnia/models"
	"time"
)

type seedProject struct {
	title       string
	description string
	techStack   string
	imageDir    string
}

// SeedProjectsIfEmpty memindahkan 4 proyek yang sebelumnya hardcoded
// di homepage ke dalam database, agar tampil lewat admin dan
// halaman homepage tetap menampilkan proyek seperti sediakala.
func SeedProjectsIfEmpty() {
	var count int64
	if err := DB.Model(&models.Project{}).Count(&count).Error; err != nil {
		fmt.Println("Cek proyek gagal:", err)
		return
	}
	if count > 0 {
		return
	}

	seeds := []seedProject{
		{
			title:       "Sistem Informasi Pegawai",
			description: "Aplikasi web untuk manajemen data pegawai perusahaan yang komprehensif, mencakup modul absensi, payroll, dan evaluasi kinerja.",
			techStack:   "Laravel, MySQL",
			imageDir:    "/public/proyek1",
		},
		{
			title:       "Indikator Mutu RS",
			description: "Sistem analitik untuk penilaian dan monitoring indikator mutu di RS Bhayangkara Jambi dengan dashboard real-time.",
			techStack:   "Laravel, Filament",
			imageDir:    "/public/proyek2",
		},
		{
			title:       "Survey Kepuasan Pasien",
			description: "Aplikasi interaktif untuk pengumpulan data survei kepuasan pasien secara efisien dengan pelaporan otomatis.",
			techStack:   "Laravel, MySQL",
			imageDir:    "/public/proyek3",
		},
		{
			title:       "Aplikasi Storage Data",
			description: "Sistem manajemen penyimpanan dokumen digital terpusat yang aman dan terpusat.",
			techStack:   "Laravel, MySQL",
			imageDir:    "/public/proyek4",
		},
	}

	for i, s := range seeds {
		project := models.Project{
			Title:       s.title,
			Slug:        models.Slugify(s.title),
			Description: s.description,
			TechStack:   s.techStack,
			Order:       i + 1,
			CreatedAt:   time.Now(),
		}
		if err := DB.Create(&project).Error; err != nil {
			fmt.Println("Seed proyek gagal:", err)
			continue
		}
		for j := 1; j <= 3; j++ {
			img := models.ProjectImage{
				ProjectID: project.ID,
				Image:     fmt.Sprintf("%s/project-gambar%d.png", s.imageDir, j),
				Order:     j,
			}
			if err := DB.Create(&img).Error; err != nil {
				fmt.Println("Seed gambar gagal:", err)
			}
		}
	}
	fmt.Printf("Seed proyek selesai: %d proyek\n", len(seeds))
}
