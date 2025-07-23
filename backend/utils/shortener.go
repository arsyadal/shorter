package utils

import (
	"crypto/rand"
	"math/big"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	shortCodeLength = 6
)

// GenerateShortCode generates a random Base62 short code
func GenerateShortCode() string {
	result := make([]byte, shortCodeLength)
	for i := range result {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		result[i] = base62Chars[num.Int64()]
	}
	return string(result)
}

// GenerateUUIDShortCode generates a short code from UUID (alternative method)
func GenerateUUIDShortCode() string {
	id := uuid.New()
	return EncodeBase62(id.ID())
}

// EncodeBase62 encodes a number to Base62
func EncodeBase62(num uint32) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	result := ""
	base := uint32(len(base62Chars))

	for num > 0 {
		result = string(base62Chars[num%base]) + result
		num = num / base
	}

	return result
}

// DecodeBase62 decodes a Base62 string to number
func DecodeBase62(str string) uint32 {
	result := uint32(0)
	base := uint32(len(base62Chars))

	for _, char := range str {
		result = result*base + uint32(strings.IndexRune(base62Chars, char))
	}

	return result
}

// IsValidURL validates if the given string is a valid URL
func IsValidURL(rawURL string) bool {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return false
	}

	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

// NormalizeURL normalizes the URL by adding http:// if no scheme is present
func NormalizeURL(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		return "http://" + rawURL
	}
	return rawURL
}

// IsValidCustomCode validates custom short code
func IsValidCustomCode(code string) bool {
	if len(code) < 3 || len(code) > 20 {
		return false
	}

	// Only allow alphanumeric characters and hyphens
	matched, _ := regexp.MatchString("^[a-zA-Z0-9-]+$", code)
	return matched
}

// GetClientIP extracts the real client IP from request
func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, get the first one
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	if strings.Contains(ip, ":") {
		ip = strings.Split(ip, ":")[0]
	}
	return ip
}

// GetTitleFromURL fetches the title of a webpage (simplified version)
func GetTitleFromURL(rawURL string) string {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(rawURL)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// Simple title extraction (you might want to use a proper HTML parser)
	// For now, return empty string - you can enhance this later
	return ""
}

// Contains checks if a slice contains a string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
} 