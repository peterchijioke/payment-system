package middlewares

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	suspiciousUserAgents = []string{
		"scrapy",
		"wget",
		"python-requests",
		"bot",
		"spider",
		"crawler",
	}

	suspiciousURLPatterns = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(union.*select|select.*from|insert.*into|delete.*from|drop.*table|update.*set)`),
		regexp.MustCompile(`(?i)(<script|javascript:|onerror=|onload=)`),
		regexp.MustCompile(`(?i)(\.\./|\.\.\\)`),
		regexp.MustCompile(`(?i)(etc/passwd|etc/shadow|win.ini|boot.ini)`),
		regexp.MustCompile(`(?i)(eval\(|base64_decode|shell_exec|system\()`),
	}
)

type MonitoringService struct {
	mu                 sync.RWMutex
	requestCounts      map[string][]time.Time
	endpointCounts     map[string]map[string]int
	suspiciousIPs      map[string]time.Time
	totalRequests      int64
	blockedRequests    int64
	suspiciousPatterns int64
	blockDuration      time.Duration
	windowDuration     time.Duration
	scrapingThreshold  int
	rateThreshold      int
}

func NewMonitoringService() *MonitoringService {
	return &MonitoringService{
		requestCounts:     make(map[string][]time.Time),
		endpointCounts:    make(map[string]map[string]int),
		suspiciousIPs:     make(map[string]time.Time),
		blockDuration:     getEnvDuration("MONITOR_BLOCK_DURATION", 5*time.Minute),
		windowDuration:    getEnvDuration("MONITOR_WINDOW_DURATION", time.Minute),
		scrapingThreshold: getEnvInt("MONITOR_SCRAPING_THRESHOLD", 100),
		rateThreshold:     getEnvInt("MONITOR_RATE_THRESHOLD", 300),
	}
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return time.Duration(intVal) * time.Minute
		}
	}
	return defaultValue
}

func Monitoring(monitor *MonitoringService) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")
		method := c.Request.Method
		path := c.Request.URL.Path

		if monitor.isBlocked(clientIP) {
			monitor.mu.Lock()
			monitor.blockedRequests++
			monitor.mu.Unlock()

			log.Printf("[BLOCKED] Suspicious IP accessing: %s %s from %s", method, path, clientIP)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Access denied. Suspicious activity detected.",
				"success": false,
			})
			return
		}

		if blocked, reason := monitor.checkSuspiciousPatterns(c, userAgent, path); blocked {
			monitor.mu.Lock()
			monitor.suspiciousPatterns++
			monitor.suspiciousIPs[clientIP] = time.Now().Add(monitor.blockDuration)
			monitor.mu.Unlock()

			log.Printf("[BLOCKED] Suspicious pattern detected from %s: %s | Path: %s", clientIP, reason, path)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Access denied. Suspicious pattern detected.",
				"success": false,
			})
			return
		}

		if monitor.detectScraping(clientIP, path) {
			monitor.mu.Lock()
			monitor.suspiciousPatterns++
			monitor.suspiciousIPs[clientIP] = time.Now().Add(monitor.blockDuration)
			monitor.mu.Unlock()

			log.Printf("[BLOCKED] Scraping detected from %s: %s %s", clientIP, method, path)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   "Access denied. Scraping behavior detected.",
				"success": false,
			})
			return
		}

		if monitor.detectDDoS(clientIP) {
			monitor.mu.Lock()
			monitor.blockedRequests++
			monitor.suspiciousIPs[clientIP] = time.Now().Add(monitor.blockDuration)
			monitor.mu.Unlock()

			log.Printf("[BLOCKED] DDoS attempt detected from %s: %s %s", clientIP, method, path)
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":   "Too many requests. You have been temporarily blocked.",
				"success": false,
			})
			return
		}

		monitor.recordRequest(clientIP, path)

		log.Printf("[REQUEST] %s %s from %s | User-Agent: %s", method, path, clientIP, truncateUserAgent(userAgent))

		monitor.mu.Lock()
		monitor.totalRequests++
		monitor.mu.Unlock()

		c.Next()

		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			log.Printf("[WARNING] %s %s returned status %d", method, path, statusCode)
		}
	}
}

func (m *MonitoringService) isBlocked(ip string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if blockTime, exists := m.suspiciousIPs[ip]; exists {
		if time.Now().Before(blockTime) {
			return true
		}
		delete(m.suspiciousIPs, ip)
	}
	return false
}

func (m *MonitoringService) checkSuspiciousPatterns(c *gin.Context, userAgent, path string) (bool, string) {
	uaLower := strings.ToLower(userAgent)
	for _, suspicious := range suspiciousUserAgents {
		if strings.Contains(uaLower, suspicious) {
			return true, "Suspicious User-Agent: " + userAgent
		}
	}

	for _, pattern := range suspiciousURLPatterns {
		if pattern.MatchString(path) {
			return true, "Malicious URL pattern: " + pattern.String()
		}
	}

	if len(path) > 500 {
		return true, "Unusually long URL"
	}

	return false, ""
}

func (m *MonitoringService) detectScraping(ip, path string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()


	if firstReq, exists := m.requestCounts[ip]; exists {
		if len(firstReq) > 0 && now.Sub(firstReq[0]) > m.windowDuration {

			delete(m.endpointCounts, ip)
		}
	}

	if _, exists := m.endpointCounts[ip]; !exists {
		m.endpointCounts[ip] = make(map[string]int)
	}

	m.endpointCounts[ip][path]++

	uniqueEndpoints := len(m.endpointCounts[ip])


	if uniqueEndpoints > m.scrapingThreshold {
		return true
	}

	return false
}

func (m *MonitoringService) detectDDoS(ip string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	requests := m.requestCounts[ip]

	var validRequests []time.Time
	for _, reqTime := range requests {
		if now.Sub(reqTime) < m.windowDuration {
			validRequests = append(validRequests, reqTime)
		}
	}

	if len(validRequests) > m.rateThreshold {
		return true
	}

	return false
}

func (m *MonitoringService) recordRequest(ip, path string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()

	m.requestCounts[ip] = append(m.requestCounts[ip], now)

	if _, exists := m.endpointCounts[ip]; !exists {
		m.endpointCounts[ip] = make(map[string]int)
	}

	m.endpointCounts[ip][path]++
}

func (m *MonitoringService) cleanupOldEntries(ip string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var validRequests []time.Time
	for _, reqTime := range m.requestCounts[ip] {
		if time.Since(reqTime) < m.windowDuration {
			validRequests = append(validRequests, reqTime)
		}
	}
	m.requestCounts[ip] = validRequests

	if len(m.requestCounts[ip]) == 0 {
		delete(m.endpointCounts, ip)
	}
}

func truncateUserAgent(ua string) string {
	if len(ua) > 50 {
		return ua[:50] + "..."
	}
	return ua
}

func (m *MonitoringService) GetMonitoringStats() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.mu.RLock()
		defer m.mu.RUnlock()

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"stats": gin.H{
				"total_requests":      m.totalRequests,
				"blocked_requests":    m.blockedRequests,
				"suspicious_patterns": m.suspiciousPatterns,
				"active_blocked_ips":  len(m.suspiciousIPs),
			},
		})
	}
}
