package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BlogPost struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"size:200"`
	Slug      string `gorm:"size:255;unique"`
	Category  string `gorm:"size:50;default:'General'"`
	Content   string
	Image     string    `gorm:"size:255"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Project struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:200"`
	Slug        string `gorm:"size:255;unique"`
	Description string
	TechStack   string `gorm:"size:200"`
	Image       string `gorm:"size:255"`
	Link        string `gorm:"size:255"`
	// Case study
	Year      string         `gorm:"size:20"`  // contoh: "2025"
	Category  string         `gorm:"size:100"` // contoh: "Healthcare"
	Role      string         `gorm:"size:150"` // contoh: "Backend Developer"
	Problem   string         `gorm:"type:text"`
	Solution  string         `gorm:"type:text"`
	Result    string         `gorm:"type:text"`
	Order     int            `gorm:"default:0"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	Images    []ProjectImage `gorm:"foreignKey:ProjectID"`
}

type ProjectImage struct {
	ID        uint `gorm:"primaryKey"`
	ProjectID uint
	Image     string `gorm:"size:255"`
	Caption   string `gorm:"size:200"`
	Order     int    `gorm:"default:0"`
}

func (BlogPost) TableName() string {
	return "landing_blogpost"
}

func (Project) TableName() string {
	return "landing_project"
}

func (ProjectImage) TableName() string {
	return "landing_projectimage"
}

// BeforeCreate membuat slug otomatis bila kosong.
func (p *Project) BeforeCreate(tx *gorm.DB) error {
	if p.Slug == "" {
		p.Slug = uniqueSlugInTx(tx, Slugify(p.Title), "landing_project")
	}
	return nil
}

// BeforeCreate membuat slug otomatis bila kosong.
func (b *BlogPost) BeforeCreate(tx *gorm.DB) error {
	if b.Slug == "" {
		b.Slug = uniqueSlugInTx(tx, Slugify(b.Title), "landing_blogpost")
	}
	return nil
}

func uniqueSlugInTx(tx *gorm.DB, base, table string) string {
	slug := base
	for i := 2; ; i++ {
		var count int64
		if err := tx.Table(table).Where("slug = ?", slug).Count(&count).Error; err != nil {
			return slug
		}
		if count == 0 {
			return slug
		}
		slug = fmt.Sprintf("%s-%d", base, i)
	}
}

type Admin struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:100;unique"`
	Password  string    `gorm:"size:255"` // Simpan hash password
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Admin) TableName() string {
	return "admin_user"
}

type SiteSetting struct {
	ID           uint   `gorm:"primaryKey"`
	SiteTitle    string `gorm:"default:'Portfolio'"`
	Favicon      string
	ProfileImage string // URL foto profil (kosongkan untuk memakai placeholder inisial)
	CVLink       string // URL CV/PDF, contoh: /public/cv.pdf
	HeroTitle    string `gorm:"default:'Hi, I am a Developer'"`
	HeroSubtitle string `gorm:"default:'Welcome to my portfolio'"`
	AboutText    string `gorm:"type:text"`
	FooterText   string `gorm:"default:'© 2026 Portfolio'"`
	GithubLink   string
	LinkedInLink string
	EmailLink    string
	WhatsAppLink string // nomor format https://wa.me/628xxxx
	Theme        string `gorm:"size:20;default:'dark'"` // dark | light
}

func (SiteSetting) TableName() string {
	return "site_setting"
}

type MenuLink struct {
	ID    uint   `gorm:"primaryKey"`
	Label string `gorm:"size:50"`
	URL   string `gorm:"size:255"`
	Order int    `gorm:"default:0"`
}

func (MenuLink) TableName() string {
	return "site_menulink"
}

// ContactMessage adalah pesan yang masuk lewat form kontak di homepage.
type ContactMessage struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"size:100"`
	Email     string    `gorm:"size:200"`
	Message   string    `gorm:"type:text"`
	IsRead    bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (ContactMessage) TableName() string {
	return "contact_message"
}

// Experience adalah entri riwayat kerja (timeline homepage + resume).
type Experience struct {
	ID      uint   `gorm:"primaryKey"`
	Company string `gorm:"size:200"`
	Role    string `gorm:"size:200"`
	Period  string `gorm:"size:100"`  // contoh: "Juni 2025 - Sekarang"
	Bullets string `gorm:"type:text"` // satu poin per baris
	Order   int    `gorm:"default:0"`
}

func (Experience) TableName() string {
	return "experience_item"
}

// Credential adalah pendidikan atau sertifikasi.
type Credential struct {
	ID          uint   `gorm:"primaryKey"`
	Type        string `gorm:"size:30;default:'education'"` // education | certification
	Title       string `gorm:"size:200"`
	Institution string `gorm:"size:200"`
	Year        string `gorm:"size:50"`
	Order       int    `gorm:"default:0"`
}

func (Credential) TableName() string {
	return "credential"
}

// Testimonial adalah kutipan klien / rekan kerja.
type Testimonial struct {
	ID    uint   `gorm:"primaryKey"`
	Quote string `gorm:"type:text"`
	Name  string `gorm:"size:100"`
	Role  string `gorm:"size:150"`
	Order int    `gorm:"default:0"`
}

func (Testimonial) TableName() string {
	return "testimonial"
}
