package entity

import "time"

type News struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Slug        string    `json:"slug" gorm:"unique;not null"`
	BannerPhoto *string   `json:"banner_photo,omitempty"` // optional
	Content     string    `json:"content" gorm:"type:text;not null"`
	AuthorID    uint      `json:"author_id" gorm:"not null"`
	Author      User      `json:"author" gorm:"foreignKey:AuthorID"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
