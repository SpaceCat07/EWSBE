package model

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/repository"

	"gorm.io/gorm"
)

type newsModel struct {
	db *gorm.DB
}

func NewNewsRepo(db *gorm.DB) repository.NewsRepository {
	return &newsModel{db: db}
}

func (r *newsModel) CreateNews(news *entity.News) error {
	return r.db.Create(news).Error
}

func (r *newsModel) GetAllNews() ([]entity.News, error) {
	var news []entity.News
	if err := r.db.Preload("Author").Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *newsModel) GetNewsByID(id uint) (*entity.News, error) {
	var news entity.News
	if err := r.db.Preload("Author").First(&news, id).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsModel) GetNewsBySlug(slug string) (*entity.News, error) {
	var news entity.News
	if err := r.db.Preload("Author").Where("slug = ?", slug).First(&news).Error; err != nil {
		return nil, err
	}
	return &news, nil
}

func (r *newsModel) UpdateNews(news *entity.News) error {
	return r.db.Save(news).Error
}

func (r *newsModel) DeleteNews(id uint) error {
	return r.db.Delete(&entity.News{}, id).Error
}

func (r *newsModel) GetNewsByAuthorID(authorID uint) ([]entity.News, error) {
	var news []entity.News
	if err := r.db.Preload("Author").Where("author_id = ?", authorID).Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}
