package internal

import (
	"fmt"
	"log"
	"os"
	"time"

	sqlite "github.com/glebarez/sqlite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() *gorm.DB {
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "postgres" // 默认使用 PostgreSQL
	}

	var db *gorm.DB
	var err error

	logLevel := logger.Info
	if os.Getenv("GIN_MODE") == "release" {
		logLevel = logger.Warn
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	switch dbType {
	case "postgres":
		db, err = initPostgres(gormConfig)
	default:
		db, err = initSQLite(gormConfig)
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 跳过自动迁移 - 使用 cmd/initdb 手动初始化数据库
	// 如果需要自动迁移，设置环境变量 AUTO_MIGRATE=true
	if os.Getenv("AUTO_MIGRATE") == "true" {
		err = db.AutoMigrate(
			&Admin{},
			&User{},
			&Category{},
			&Tag{},
			&SPU{},
			&SKU{},
			&SPUTag{},
			&Address{},
			&Order{},
			&OrderItem{},
			&CartItem{},
			&Comment{},
			&Coupon{},
			&Promotion{},
			&Swiper{},
			&RecommendedProduct{},
			&HomeContent{},
		)

		if err != nil {
			log.Fatalf("Failed to migrate database: %v", err)
		}
		log.Println("Database migrated successfully")
	}

	DB = db
	log.Println("Database connected successfully")

	// 自动检测并初始化数据库
	if err := AutoInitDB(); err != nil {
		log.Printf("Warning: Auto initialization failed: %v", err)
	}

	// 创建测试用户（如果不存在）
	EnsureTestUser(db)

	return db
}

// EnsureTestUser 确保测试用户存在
func EnsureTestUser(db *gorm.DB) {
	testOpenID := "oTest_dev_openid_001"
	var user User
	if err := db.First(&user, "open_id = ?", testOpenID).Error; err != nil {
		// 用户不存在，创建测试用户
		user = User{
			ID:        GenerateUUID(),
			OpenID:    testOpenID,
			NickName:  "测试用户",
			Avatar:    "",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := db.Create(&user).Error; err != nil {
			log.Printf("Failed to create test user: %v", err)
		} else {
			log.Printf("Test user created: openid=%s, id=%s", testOpenID, user.ID)
		}
	} else {
		log.Printf("Test user already exists: openid=%s, id=%s", testOpenID, user.ID)
	}
}

func initSQLite(config *gorm.Config) (*gorm.DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./z26b.db"
	}
	return gorm.Open(sqlite.Open(dbPath), config)
}

func initPostgres(config *gorm.Config) (*gorm.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" {
		host = "localhost"
	}
	if port == "" {
		port = "5432"
	}
	if user == "" {
		user = "postgres"
	}
	if dbName == "" {
		dbName = "z26b"
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(dsn), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}

// GetDB 返回数据库实例
func GetDB() *gorm.DB {
	return DB
}
