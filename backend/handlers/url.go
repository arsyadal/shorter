package handlers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"shorter-backend/config"
	"shorter-backend/models"
	"shorter-backend/utils"

	"github.com/gin-gonic/gin"
)

// ShortenURL creates a new short URL
func ShortenURL(c *gin.Context) {
	var req models.CreateURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Normalize and validate URL
	normalizedURL := utils.NormalizeURL(req.URL)
	if !utils.IsValidURL(normalizedURL) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid URL format"})
		return
	}

	// Check if URL already exists
	var existingURL models.URL
	if err := config.DB.Where("original_url = ?", normalizedURL).First(&existingURL).Error; err == nil {
		// URL already exists, return existing short code
		baseURL := getBaseURL(c)
		response := models.URLResponse{
			ID:          existingURL.ID,
			OriginalURL: existingURL.OriginalURL,
			ShortCode:   existingURL.ShortCode,
			ShortURL:    fmt.Sprintf("%s/%s", baseURL, existingURL.ShortCode),
			Title:       existingURL.Title,
			CreatedAt:   existingURL.CreatedAt,
		}

		// Get click count
		var clickCount int64
		config.DB.Model(&models.Click{}).Where("url_id = ?", existingURL.ID).Count(&clickCount)
		response.ClickCount = clickCount

		c.JSON(http.StatusOK, response)
		return
	}

	// Generate short code
	var shortCode string
	if req.CustomCode != "" {
		// Validate custom code
		if !utils.IsValidCustomCode(req.CustomCode) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid custom code format"})
			return
		}

		// Check if custom code already exists
		var existingCustom models.URL
		if err := config.DB.Where("short_code = ?", req.CustomCode).First(&existingCustom).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Custom code already exists"})
			return
		}
		shortCode = req.CustomCode
	} else {
		// Generate random short code
		for {
			shortCode = utils.GenerateShortCode()
			var existing models.URL
			if err := config.DB.Where("short_code = ?", shortCode).First(&existing).Error; err != nil {
				break // Short code is unique
			}
		}
	}

	// Get title from URL (optional)
	title := utils.GetTitleFromURL(normalizedURL)

	// Create new URL entry
	newURL := models.URL{
		OriginalURL: normalizedURL,
		ShortCode:   shortCode,
		Title:       title,
	}

	if err := config.DB.Create(&newURL).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create short URL"})
		return
	}

	// Cache the short URL
	config.CacheSet(shortCode, normalizedURL, 24*time.Hour)

	// Prepare response
	baseURL := getBaseURL(c)
	response := models.URLResponse{
		ID:          newURL.ID,
		OriginalURL: newURL.OriginalURL,
		ShortCode:   newURL.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", baseURL, newURL.ShortCode),
		Title:       newURL.Title,
		ClickCount:  0,
		CreatedAt:   newURL.CreatedAt,
	}

	c.JSON(http.StatusCreated, response)
}

// RedirectURL handles the redirect from short URL to original URL
func RedirectURL(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	var originalURL string
	var urlID uint

	// Try to get from cache first
	cachedURL, err := config.CacheGet(shortCode)
	if err == nil && cachedURL != "" {
		originalURL = cachedURL
		// Get URL ID from database for click tracking
		var url models.URL
		if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err == nil {
			urlID = url.ID
		}
	} else {
		// Get from database
		var url models.URL
		if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		originalURL = url.OriginalURL
		urlID = url.ID

		// Cache for future requests
		config.CacheSet(shortCode, originalURL, 24*time.Hour)
	}

	// Track click asynchronously
	go trackClick(urlID, c.Request)

	// Redirect to original URL
	c.Redirect(http.StatusMovedPermanently, originalURL)
}

// GetURLStats returns statistics for a short URL
func GetURLStats(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	// Get URL from database
	var url models.URL
	if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Get total clicks
	var totalClicks int64
	config.DB.Model(&models.Click{}).Where("url_id = ?", url.ID).Count(&totalClicks)

	// Get daily clicks (last 30 days)
	var dailyClicks []models.DailyClickStat
	config.DB.Raw(`
		SELECT DATE(created_at) as date, COUNT(*) as count 
		FROM clicks 
		WHERE url_id = ? AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE(created_at) 
		ORDER BY date DESC
	`, url.ID).Scan(&dailyClicks)

	// Get country clicks (top 10)
	var countryClicks []models.CountryClickStat
	config.DB.Raw(`
		SELECT country, COUNT(*) as count 
		FROM clicks 
		WHERE url_id = ? AND country != '' 
		GROUP BY country 
		ORDER BY count DESC 
		LIMIT 10
	`, url.ID).Scan(&countryClicks)

	// Get referer clicks (top 10)
	var refererClicks []models.RefererClickStat
	config.DB.Raw(`
		SELECT referer, COUNT(*) as count 
		FROM clicks 
		WHERE url_id = ? AND referer != '' 
		GROUP BY referer 
		ORDER BY count DESC 
		LIMIT 10
	`, url.ID).Scan(&refererClicks)

	response := models.ClickStatsResponse{
		TotalClicks:    totalClicks,
		DailyClicks:    dailyClicks,
		CountryClicks:  countryClicks,
		RefererClicks:  refererClicks,
	}

	c.JSON(http.StatusOK, response)
}

// GetAllURLs returns all URLs created (with pagination)
func GetAllURLs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	var urls []models.URL
	var total int64

	// Get total count
	config.DB.Model(&models.URL{}).Count(&total)

	// Get URLs with pagination
	config.DB.Limit(limit).Offset(offset).Order("created_at desc").Find(&urls)

	// Prepare response
	baseURL := getBaseURL(c)
	var responses []models.URLResponse
	
	for _, url := range urls {
		// Get click count for each URL
		var clickCount int64
		config.DB.Model(&models.Click{}).Where("url_id = ?", url.ID).Count(&clickCount)
		
		response := models.URLResponse{
			ID:          url.ID,
			OriginalURL: url.OriginalURL,
			ShortCode:   url.ShortCode,
			ShortURL:    fmt.Sprintf("%s/%s", baseURL, url.ShortCode),
			Title:       url.Title,
			ClickCount:  clickCount,
			CreatedAt:   url.CreatedAt,
		}
		responses = append(responses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"urls":       responses,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// trackClick records a click for analytics
func trackClick(urlID uint, r *http.Request) {
	if urlID == 0 {
		return
	}

	click := models.Click{
		URLId:     urlID,
		IPAddress: utils.GetClientIP(r),
		UserAgent: r.UserAgent(),
		Referer:   r.Referer(),
		// You can enhance this with IP geolocation service
		Country:   "", // Get from IP geolocation service
		City:      "", // Get from IP geolocation service
	}

	config.DB.Create(&click)
}

// getBaseURL returns the base URL for short links
func getBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	
	host := c.Request.Host
	if customDomain := os.Getenv("CUSTOM_DOMAIN"); customDomain != "" {
		host = customDomain
	}
	
	return fmt.Sprintf("%s://%s", scheme, host)
} 