package repository

import "EWSBE/internal/entity"

type NewsRepository interface {
	CreateNews(news *entity.News) error
	GetAllNews() ([]entity.News, error)
	GetNewsByID(id uint) (*entity.News, error)
	GetNewsBySlug(slug string) (*entity.News, error)
	UpdateNews(news *entity.News) error
	DeleteNews(id uint) error
	GetNewsByAuthorID(authorID uint) ([]entity.News, error)
}
