package handlers

import (
	"net/http"
	"runtime"
	"time"

	"shorter-backend/config"
	"shorter-backend/models"

	"github.com/gin-gonic/gin"
)

// SystemStats holds system statistics
type SystemStats struct {
	// System info
	GoVersion     string `json:"go_version"`
	NumGoroutines int    `json:"num_goroutines"`
	NumCPU        int    `json:"num_cpu"`
	
	// Memory stats
	MemoryAlloc      uint64 `json:"memory_alloc_mb"`
	MemoryTotalAlloc uint64 `json:"memory_total_alloc_mb"`
	MemorySystem     uint64 `json:"memory_system_mb"`
	NumGC            uint32 `json:"num_gc"`
	
	// Database stats
	TotalURLs        int64     `json:"total_urls"`
	TotalClicks      int64     `json:"total_clicks"`
	URLsToday        int64     `json:"urls_today"`
	ClicksToday      int64     `json:"clicks_today"`
	TopDomains       []DomainStat `json:"top_domains"`
	
	// System status
	Uptime           string    `json:"uptime"`
	DatabaseStatus   string    `json:"database_status"`
	RedisStatus      string    `json:"redis_status"`
	LastHealthCheck  time.Time `json:"last_health_check"`
}

type DomainStat struct {
	Domain string `json:"domain"`
	Count  int64  `json:"count"`
}

var startTime = time.Now()

// GetSystemStats returns comprehensive system statistics
func GetSystemStats(c *gin.Context) {
	var stats SystemStats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	// System info
	stats.GoVersion = runtime.Version()
	stats.NumGoroutines = runtime.NumGoroutine()
	stats.NumCPU = runtime.NumCPU()
	
	// Memory stats (convert bytes to MB)
	stats.MemoryAlloc = bytesToMB(m.Alloc)
	stats.MemoryTotalAlloc = bytesToMB(m.TotalAlloc)
	stats.MemorySystem = bytesToMB(m.Sys)
	stats.NumGC = m.NumGC
	
	// Database stats
	config.DB.Model(&models.URL{}).Count(&stats.TotalURLs)
	config.DB.Model(&models.Click{}).Count(&stats.TotalClicks)
	
	// Today's stats
	today := time.Now().Format("2006-01-02")
	config.DB.Model(&models.URL{}).Where("DATE(created_at) = ?", today).Count(&stats.URLsToday)
	config.DB.Model(&models.Click{}).Where("DATE(created_at) = ?", today).Count(&stats.ClicksToday)
	
	// Top domains (last 30 days)
	config.DB.Raw(`
		SELECT 
			REGEXP_REPLACE(original_url, '^https?://(?:www\.)?([^/]+).*', '\1') as domain,
			COUNT(*) as count
		FROM urls 
		WHERE created_at >= NOW() - INTERVAL '30 days'
		GROUP BY domain 
		ORDER BY count DESC 
		LIMIT 10
	`).Scan(&stats.TopDomains)
	
	// System status
	stats.Uptime = time.Since(startTime).String()
	stats.DatabaseStatus = checkDatabaseStatus()
	stats.RedisStatus = checkRedisStatus()
	stats.LastHealthCheck = time.Now()
	
	c.JSON(http.StatusOK, stats)
}

// GetDetailedHealth returns detailed health check
func GetDetailedHealth(c *gin.Context) {
	health := gin.H{
		"status":    "ok",
		"timestamp": time.Now(),
		"services":  gin.H{},
	}
	
	// Check database
	if err := config.DB.Exec("SELECT 1").Error; err != nil {
		health["status"] = "degraded"
		health["services"].(gin.H)["database"] = gin.H{
			"status": "down",
			"error":  err.Error(),
		}
	} else {
		health["services"].(gin.H)["database"] = gin.H{
			"status": "up",
		}
	}
	
	// Check Redis
	if config.RDB != nil {
		if err := config.RDB.Ping(config.Ctx).Err(); err != nil {
			health["status"] = "degraded"
			health["services"].(gin.H)["redis"] = gin.H{
				"status": "down",
				"error":  err.Error(),
			}
		} else {
			health["services"].(gin.H)["redis"] = gin.H{
				"status": "up",
			}
		}
	} else {
		health["services"].(gin.H)["redis"] = gin.H{
			"status": "disabled",
		}
	}
	
	c.JSON(http.StatusOK, health)
}

// GetRecentActivity returns recent system activity
func GetRecentActivity(c *gin.Context) {
	limit := 50
	
	var recentURLs []models.URL
	var recentClicks []models.Click
	
	// Get recent URLs
	config.DB.Order("created_at desc").Limit(limit).Find(&recentURLs)
	
	// Get recent clicks with URL info
	config.DB.Preload("URL").Order("created_at desc").Limit(limit).Find(&recentClicks)
	
	activity := gin.H{
		"recent_urls":   recentURLs,
		"recent_clicks": recentClicks,
		"timestamp":     time.Now(),
	}
	
	c.JSON(http.StatusOK, activity)
}

// Utility functions
func bytesToMB(b uint64) uint64 {
	return b / 1024 / 1024
}

func checkDatabaseStatus() string {
	if err := config.DB.Exec("SELECT 1").Error; err != nil {
		return "down"
	}
	return "up"
}

func checkRedisStatus() string {
	if config.RDB == nil {
		return "disabled"
	}
	if err := config.RDB.Ping(config.Ctx).Err(); err != nil {
		return "down"
	}
	return "up"
} 