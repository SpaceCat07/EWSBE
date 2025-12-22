package usecase

import (
	"EWSBE/internal/entity"
	"EWSBE/internal/repository"
	"errors"
	"strings"
	"unicode"

	"gorm.io/gorm"
)

type NewsUsecase struct {
	newsRepo repository.NewsRepository
}

func NewNewsUsecase(newsRepo repository.NewsRepository) *NewsUsecase {
	return &NewsUsecase{newsRepo: newsRepo}
}

func (uc *NewsUsecase) CreateNews(title, content string, bannerPhoto *string, authorID uint) (*entity.News, error) {
	if strings.TrimSpace(title) == "" || strings.TrimSpace(content) == "" {
		return nil, errors.New("title and content cannot be empty")
	}

	slug := generateSlug(title)

	news := &entity.News{
		Title:       title,
		Slug:        slug,
		BannerPhoto: bannerPhoto,
		Content:     content,
		AuthorID:    authorID,
	}

	if err := uc.newsRepo.CreateNews(news); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, errors.New("slug already exists")
		}
		return nil, err
	}

	// Load author
	newsWithAuthor, err := uc.newsRepo.GetNewsBySlug(slug)
	if err != nil {
		return nil, err
	}

	return newsWithAuthor, nil
}

func (uc *NewsUsecase) GetAllNews() ([]entity.News, error) {
	return uc.newsRepo.GetAllNews()
}

func (uc *NewsUsecase) GetNewsBySlug(slug string) (*entity.News, error) {
	return uc.newsRepo.GetNewsBySlug(slug)
}

func (uc *NewsUsecase) UpdateNews(id uint, title, content string, bannerPhoto *string, authorID uint) (*entity.News, error) {
	news, err := uc.newsRepo.GetNewsByID(id)
	if err != nil {
		return nil, errors.New("news not found")
	}

	if news.AuthorID != authorID {
		return nil, errors.New("unauthorized")
	}

	if title != "" {
		news.Title = title
		news.Slug = generateSlug(title)
	}
	if content != "" {
		news.Content = content
	}
	if bannerPhoto != nil {
		news.BannerPhoto = bannerPhoto
	}

	if err := uc.newsRepo.UpdateNews(news); err != nil {
		return nil, err
	}

	return news, nil
}

func (uc *NewsUsecase) DeleteNews(id uint, authorID uint) error {
	news, err := uc.newsRepo.GetNewsByID(id)
	if err != nil {
		return errors.New("news not found")
	}

	if news.AuthorID != authorID {
		return errors.New("unauthorized")
	}

	return uc.newsRepo.DeleteNews(id)
}

func (uc *NewsUsecase) GetNewsByAuthorID(authorID uint) ([]entity.News, error) {
	return uc.newsRepo.GetNewsByAuthorID(authorID)
}

// Helper function to generate slug
func generateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) || r == ' ' {
			return r
		}
		return -1
	}, slug)
	slug = strings.ReplaceAll(slug, " ", "-")
	return strings.Trim(slug, "-")
}
