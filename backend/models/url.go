package models

import (
	"time"
	"gorm.io/gorm"
)

type URL struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	OriginalURL string         `json:"original_url" gorm:"not null;index"`
	ShortCode   string         `json:"short_code" gorm:"uniqueIndex;not null"`
	Title       string         `json:"title"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
	Clicks      []Click        `json:"clicks,omitempty" gorm:"foreignKey:URLId"`
	ClickCount  int64          `json:"click_count" gorm:"-"`
}

type Click struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	URLId     uint           `json:"url_id" gorm:"not null;index"`
	IPAddress string         `json:"ip_address"`
	UserAgent string         `json:"user_agent"`
	Referer   string         `json:"referer"`
	Country   string         `json:"country"`
	City      string         `json:"city"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type CreateURLRequest struct {
	URL       string `json:"url" binding:"required,url"`
	CustomCode string `json:"custom_code,omitempty"`
}

type URLResponse struct {
	ID          uint      `json:"id"`
	OriginalURL string    `json:"original_url"`
	ShortCode   string    `json:"short_code"`
	ShortURL    string    `json:"short_url"`
	Title       string    `json:"title"`
	ClickCount  int64     `json:"click_count"`
	CreatedAt   time.Time `json:"created_at"`
}

type ClickStatsResponse struct {
	TotalClicks    int64                    `json:"total_clicks"`
	DailyClicks    []DailyClickStat         `json:"daily_clicks"`
	CountryClicks  []CountryClickStat       `json:"country_clicks"`
	RefererClicks  []RefererClickStat       `json:"referer_clicks"`
}

type DailyClickStat struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}

type CountryClickStat struct {
	Country string `json:"country"`
	Count   int64  `json:"count"`
}

type RefererClickStat struct {
	Referer string `json:"referer"`
	Count   int64  `json:"count"`
} 