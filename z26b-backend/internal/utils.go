package internal

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
)

// ============================================
// 结构化日志工具
// ============================================

// LogLevel 日志级别
type LogLevel int

const (
	INFO LogLevel = iota
	WARN
	ERROR
)

// LogEntry 日志条目
type LogEntry struct {
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Timestamp string                 `json:"timestamp"`
	Service   string                 `json:"service,omitempty"`
	UserID    string                 `json:"userId,omitempty"`
	IP        string                 `json:"ip,omitempty"`
	Method    string                 `json:"method,omitempty"`
	Path      string                 `json:"path,omitempty"`
	Duration  string                 `json:"duration,omitempty"`
	Status    int                    `json:"status,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

// Logger 结构化日志器
type Logger struct {
	service string
}

// NewLogger 创建新的日志器
func NewLogger(service string) *Logger {
	return &Logger{service: service}
}

// Info 记录信息日志
func (l *Logger) Info(message string, extra ...map[string]interface{}) {
	var extraMap map[string]interface{}
	if len(extra) > 0 {
		extraMap = extra[0]
	}
	l.log(INFO, message, []map[string]interface{}{extraMap})
}

// Warn 记录警告日志
func (l *Logger) Warn(message string, extra ...map[string]interface{}) {
	var extraMap map[string]interface{}
	if len(extra) > 0 {
		extraMap = extra[0]
	}
	l.log(WARN, message, []map[string]interface{}{extraMap})
}

// Error 记录错误日志
func (l *Logger) Error(message string, err error, extra ...map[string]interface{}) {
	entry := LogEntry{
		Level:     "ERROR",
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   l.service,
	}
	if err != nil {
		entry.Error = err.Error()
	}
	if len(extra) > 0 {
		entry.Extra = extra[0]
	}
	l.writeLog(entry)
}

func (l *Logger) log(level LogLevel, message string, extra []map[string]interface{}) {
	var levelStr string
	switch level {
	case INFO:
		levelStr = "INFO"
	case WARN:
		levelStr = "WARN"
	case ERROR:
		levelStr = "ERROR"
	}

	entry := LogEntry{
		Level:     levelStr,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
		Service:   l.service,
	}
	if len(extra) > 0 {
		entry.Extra = extra[0]
	}
	l.writeLog(entry)
}

func (l *Logger) writeLog(entry LogEntry) {
	// 在开发环境中输出到控制台，在生产环境中可以输出到文件或外部服务
	if os.Getenv("GIN_MODE") == "release" {
		// 生产环境：JSON格式
		if jsonData, err := json.Marshal(entry); err == nil {
			log.Println(string(jsonData))
		}
	} else {
		// 开发环境：可读格式
		log.Printf("[%s] %s - %s", entry.Level, entry.Service, entry.Message)
		if entry.Error != "" {
			log.Printf("Error: %s", entry.Error)
		}
	}
}

// GlobalLogger 全局日志器
var GlobalLogger = NewLogger("z26b-backend")

// ============================================
// ID & Hash 工具
// ============================================

func GenerateID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("%d_%x", time.Now().UnixNano(), b)
}

func GenerateUUID() string {
	return uuid.New().String()
}

func HashString(s string) string {
	hashedBytes, _ := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	return string(hashedBytes)
}

// VerifyPassword 验证密码
func VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ToJSON(v interface{}) datatypes.JSON {
	data, _ := json.Marshal(v)
	return datatypes.JSON(data)
}

// ============================================
// JWT 认证中间件
// ============================================

var jwtSecret = []byte("temp-secret-key") // 临时密钥，实际使用时会重新设置

// init 初始化函数
func init() {
	// 延迟初始化JWT密钥，避免在测试时出错
}

// GetJWTSecret 获取JWT密钥

// GetJWTSecret 获取JWT密钥
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		if os.Getenv("GIN_MODE") == "release" {
			panic("JWT_SECRET environment variable must be set in production")
		}
		secret = "dev-secret-key-change-in-production"
	}
	// 更新全局密钥
	jwtSecret = []byte(secret)
	return secret
}

type AdminClaims struct {
	AdminID  string `json:"adminId"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAdminToken 生成管理员JWT
func GenerateAdminToken(adminID, email, username, role string) (string, error) {
	claims := AdminClaims{
		AdminID:  adminID,
		Email:    email,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "z26b-admin",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// AdminAuthMiddleware 管理员认证中间件
func AdminAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证已过期，请重新登录"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*AdminClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证信息"})
			c.Abort()
			return
		}

		c.Set("adminID", claims.AdminID)
		c.Set("adminEmail", claims.Email)
		c.Set("adminUsername", claims.Username)
		c.Set("adminRole", claims.Role)

		c.Next()
	}
}

// SuperAdminMiddleware 需要超级管理员权限
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("adminRole")
		if !exists || role != "super_admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要超级管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// ============================================
// 输入验证工具
// ============================================

// ValidateEmail 验证邮箱格式
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidatePassword 验证密码强度
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("密码长度至少8位")
	}
	if len(password) > 128 {
		return fmt.Errorf("密码长度不能超过128位")
	}
	// 检查是否包含至少一个字母和一个数字
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasLetter || !hasNumber {
		return fmt.Errorf("密码必须包含字母和数字")
	}
	return nil
}

// ValidateUsername 验证用户名
func ValidateUsername(username string) error {
	if len(username) < 3 || len(username) > 50 {
		return fmt.Errorf("用户名长度必须在3-50位之间")
	}
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(username) {
		return fmt.Errorf("用户名只能包含字母、数字、下划线和连字符")
	}
	return nil
}

// SanitizeString 清理字符串输入
func SanitizeString(input string) string {
	// 移除前后空格
	input = strings.TrimSpace(input)
	// 限制长度
	if utf8.RuneCountInString(input) > 1000 {
		runes := []rune(input)
		input = string(runes[:1000])
	}
	// 移除潜在的XSS字符
	input = strings.ReplaceAll(input, "<", "&lt;")
	input = strings.ReplaceAll(input, ">", "&gt;")
	input = strings.ReplaceAll(input, "\"", "&quot;")
	input = strings.ReplaceAll(input, "'", "&#x27;")
	return input
}

// ============================================
// SPU 价格同步
// ============================================

// UpdateSPUPriceRange 根据 SKU 更新 SPU 的价格范围
func UpdateSPUPriceRange(db interface {
	Exec(sql string, values ...interface{}) interface{ Error() error }
}, spuID string) error {
	sql := `
		UPDATE spu SET 
			min_price = COALESCE((SELECT MIN(price) FROM sku WHERE "SPUID" = $1), 0),
			max_price = COALESCE((SELECT MAX(price) FROM sku WHERE "SPUID" = $1), 0)
		WHERE id = $1
	`
	if err := db.Exec(sql, spuID).Error(); err != nil {
		return fmt.Errorf("failed to update SPU price range: %w", err)
	}
	return nil
}
