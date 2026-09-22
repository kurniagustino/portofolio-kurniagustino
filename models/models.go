package models

import "time"

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
	TechStack   string         `gorm:"size:200"`
	Image       string         `gorm:"size:255"`
	Link        string         `gorm:"size:255"`
	Order       int            `gorm:"default:0"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	Images      []ProjectImage `gorm:"foreignKey:ProjectID"`
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

type Admin struct {
	ID        uint      `gorm:"primaryKey"`
	Username  string    `gorm:"size:100;unique"`
	Password  string    `gorm:"size:255"` // Simpan hash password
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

func (Admin) TableName() string {
	return "admin_user"
}
