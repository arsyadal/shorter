package handlers

import (
	"net/http"

	"shorter-backend/config"
	"shorter-backend/models"

	"github.com/gin-gonic/gin"
	"github.com/skip2/go-qrcode"
)

// GenerateQRCode generates QR code for a short URL
func GenerateQRCode(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	// Check if URL exists
	var url models.URL
	if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Get base URL for short URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + c.Request.Host
	shortURL := baseURL + "/" + shortCode

	// Generate QR Code for SHORT URL (not original URL)
	qrBytes, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	// Set response headers
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "inline; filename=qr-"+shortCode+".png")
	c.Header("Cache-Control", "public, max-age=3600")
	
	// Return image
	c.Data(http.StatusOK, "image/png", qrBytes)
}

// GetQRCodeHTML returns HTML page with QR code
func GetQRCodeHTML(c *gin.Context) {
	shortCode := c.Param("code")
	if shortCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Short code is required"})
		return
	}

	// Check if URL exists
	var url models.URL
	if err := config.DB.Where("short_code = ?", shortCode).First(&url).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
		return
	}

	// Get base URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	baseURL := scheme + "://" + c.Request.Host
	
	qrURL := baseURL + "/api/qr/" + shortCode + "/image"
	shortURL := baseURL + "/" + shortCode

	html := `<!DOCTYPE html>
<html>
<head>
    <title>QR Code - ` + shortCode + `</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, sans-serif; margin: 0; padding: 40px; background: #f8f9fa; }
        .container { max-width: 600px; margin: 0 auto; background: white; padding: 40px; border-radius: 12px; box-shadow: 0 2px 20px rgba(0,0,0,0.1); text-align: center; }
        .qr-code { margin: 20px 0; }
        .qr-code img { max-width: 256px; border: 1px solid #ddd; border-radius: 8px; }
        .url-info { background: #f8f9fa; padding: 20px; border-radius: 8px; margin: 20px 0; }
        .short-url { font-family: monospace; font-size: 18px; color: #0066cc; word-break: break-all; }
        .original-url { color: #666; font-size: 14px; word-break: break-all; margin-top: 10px; }
        .download-btn { display: inline-block; background: #0066cc; color: white; padding: 10px 20px; text-decoration: none; border-radius: 6px; margin: 10px; }
        .download-btn:hover { background: #0052a3; }
    </style>
</head>
<body>
    <div class="container">
        <h1>QR Code</h1>
        <div class="qr-code">
            <img src="` + qrURL + `" alt="QR Code for ` + shortCode + `">
        </div>
        <div class="url-info">
            <div class="short-url">` + shortURL + `</div>
            <div class="original-url">→ ` + url.OriginalURL + `</div>
        </div>
        <a href="` + qrURL + `" class="download-btn" download="qr-` + shortCode + `.png">Download QR Code</a>
        <a href="` + shortURL + `" class="download-btn">Visit URL</a>
    </div>
</body>
</html>`

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
} 