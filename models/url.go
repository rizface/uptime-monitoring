package models

import (
	"time"
	"gorm.io/gorm"
)

type URL struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	URL          string    `json:"url" gorm:"not null;unique"`
	Status       string    `json:"status" gorm:"default:unknown"`
	CheckInterval int      `json:"check_interval" gorm:"default:300"`
	LastChecked  *time.Time `json:"last_checked"`
	ResponseTime int       `json:"response_time" gorm:"default:0"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type URLRepository struct {
	DB *gorm.DB
}

func NewURLRepository(db *gorm.DB) *URLRepository {
	return &URLRepository{DB: db}
}

func (r *URLRepository) GetAll() ([]URL, error) {
	var urls []URL
	err := r.DB.Order("created_at DESC").Find(&urls).Error
	return urls, err
}

func (r *URLRepository) GetByID(id uint) (*URL, error) {
	var url URL
	err := r.DB.First(&url, id).Error
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *URLRepository) Create(url *URL) error {
	url.Status = "unknown"
	return r.DB.Create(url).Error
}

func (r *URLRepository) Update(url *URL) error {
	return r.DB.Save(url).Error
}

func (r *URLRepository) Delete(id uint) error {
	return r.DB.Delete(&URL{}, id).Error
}

func (r *URLRepository) UpdateStatus(id uint, status string, responseTime int) error {
	now := time.Now()
	return r.DB.Model(&URL{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":        status,
		"last_checked":  &now,
		"response_time": responseTime,
		"updated_at":    now,
	}).Error
}