package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// generateUUID 生成UUID
func generateUUID() string {
	return uuid.New().String()
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return uuid.New().String()
	}
	return hex.EncodeToString(b)[:length]
}

// hashPassword 加密密码
func hashPassword(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}
	return string(hashedPassword)
}

// DropAllTables 删除所有表
func DropAllTables(db *gorm.DB) error {
	log.Println("🗑️  Dropping existing tables...")

	tables := []string{
		"recommended_product", "home_content",
		"spu_tag", "order_item", "cart_item", "comment",
		"sku", "spu", "tag", "category",
		"address", `"order"`, "coupon", "promotion", "swiper",
		`"user"`, "admin",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	log.Println("✅ Tables dropped successfully")
	return nil
}

// CreateTables 创建所有表
func CreateTables(db *gorm.DB) error {
	log.Println("📝 Creating tables...")

	sqlStatements := []string{
		// Admin
		`CREATE TABLE IF NOT EXISTS admin (
			id TEXT PRIMARY KEY,
			email TEXT UNIQUE,
			password TEXT,
			username TEXT,
			role TEXT DEFAULT 'admin',
			status TEXT DEFAULT 'active',
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// User
		`CREATE TABLE IF NOT EXISTS "user" (
			id TEXT PRIMARY KEY,
			open_id TEXT UNIQUE,
			nick_name TEXT,
			avatar TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// Category
		`CREATE TABLE IF NOT EXISTS category (
			id TEXT PRIMARY KEY,
			name TEXT,
			icon TEXT,
			image TEXT,
			parent_id TEXT,
			sort INTEGER,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// Tag
		`CREATE TABLE IF NOT EXISTS tag (
			id TEXT PRIMARY KEY,
			name TEXT UNIQUE,
			description TEXT,
			color TEXT,
			sort_order INTEGER,
			status TEXT DEFAULT 'active',
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// SPU (商品)
		`CREATE TABLE IF NOT EXISTS spu (
			id TEXT PRIMARY KEY,
			name TEXT,
			detail TEXT,
			cover_image TEXT,
			swipe_images JSONB,
			category_id TEXT REFERENCES category(id),
			min_price DECIMAL(10,2) DEFAULT 0,
			max_price DECIMAL(10,2) DEFAULT 0,
			status TEXT,
			priority INTEGER,
			owner TEXT,
			created_at BIGINT,
			updated_at BIGINT,
			created_by TEXT,
			updated_by TEXT,
			"_openid" TEXT
		)`,

		// SKU (商品规格)
		`CREATE TABLE IF NOT EXISTS sku (
			id TEXT PRIMARY KEY,
			"SPUID" TEXT REFERENCES spu(id),
			description TEXT,
			image TEXT,
			price DECIMAL(10,2),
			count INTEGER,
			owner TEXT,
			created_at BIGINT,
			updated_at BIGINT,
			created_by TEXT,
			updated_by TEXT,
			"_openid" TEXT
		)`,

		// SPU Tag (商品标签关联)
		`CREATE TABLE IF NOT EXISTS spu_tag (
			id TEXT PRIMARY KEY,
			spu_id TEXT REFERENCES spu(id),
			tag_id TEXT REFERENCES tag(id),
			created_at TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_spu_tag_spu_id ON spu_tag(spu_id)`,
		`CREATE INDEX IF NOT EXISTS idx_spu_tag_tag_id ON spu_tag(tag_id)`,

		// Address
		`CREATE TABLE IF NOT EXISTS address (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			name TEXT,
			phone TEXT,
			country_name TEXT,
			country_code TEXT,
			province_name TEXT,
			province_code TEXT,
			city_name TEXT,
			city_code TEXT,
			district_name TEXT,
			district_code TEXT,
			detail_address TEXT,
			is_default INTEGER DEFAULT 0,
			address_tag TEXT,
			latitude TEXT,
			longitude TEXT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// Order
		`CREATE TABLE IF NOT EXISTS "order" (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			status TEXT,
			delivery_info JSONB,
			total_price DECIMAL(10,2),
			discount_price DECIMAL(10,2),
			final_price DECIMAL(10,2),
			remarks TEXT,
			created_at BIGINT,
			updated_at BIGINT
		)`,

		// OrderItem
		`CREATE TABLE IF NOT EXISTS order_item (
			id TEXT PRIMARY KEY,
			order_id TEXT REFERENCES "order"(id),
			sku_id TEXT REFERENCES sku(id),
			quantity INTEGER,
			price DECIMAL(10,2),
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// CartItem
		`CREATE TABLE IF NOT EXISTS cart_item (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			sku_id TEXT REFERENCES sku(id),
			quantity INTEGER,
			is_selected BOOLEAN DEFAULT true,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// Comment
		`CREATE TABLE IF NOT EXISTS comment (
			id TEXT PRIMARY KEY,
			spu_id TEXT,
			sku_id TEXT,
			user_id TEXT,
			user_name TEXT,
			user_head_url TEXT,
			comment_content TEXT,
			comment_score INTEGER,
			comment_resources JSONB,
			is_anonymity BOOLEAN DEFAULT false,
			seller_reply TEXT,
			created_at BIGINT,
			updated_at TIMESTAMP
		)`,

		// Coupon
		`CREATE TABLE IF NOT EXISTS coupon (
			id TEXT PRIMARY KEY,
			code TEXT,
			discount_type TEXT,
			discount_value DECIMAL(10,2),
			min_amount DECIMAL(10,2),
			max_amount DECIMAL(10,2),
			usage_limit INTEGER,
			usage_count INTEGER DEFAULT 0,
			status TEXT,
			valid_from BIGINT,
			valid_until BIGINT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// Promotion
		`CREATE TABLE IF NOT EXISTS promotion (
			id TEXT PRIMARY KEY,
			title TEXT,
			promotion_code TEXT,
			promotion_sub_code TEXT,
			tag TEXT,
			description TEXT,
			tag_text JSONB,
			promotion_status INTEGER,
			min_amount DECIMAL(10,2),
			valid_from BIGINT,
			valid_until BIGINT,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// Swiper
		`CREATE TABLE IF NOT EXISTS swiper (
			id TEXT PRIMARY KEY,
			images JSONB,
			title TEXT,
			link TEXT,
			priority INTEGER,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,

		// RecommendedProduct
		`CREATE TABLE IF NOT EXISTS recommended_product (
			id TEXT PRIMARY KEY,
			spu_id TEXT REFERENCES spu(id),
			tags JSONB,
			priority INTEGER DEFAULT 0,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_recommended_product_spu_id ON recommended_product(spu_id)`,
		`CREATE INDEX IF NOT EXISTS idx_recommended_product_priority ON recommended_product(priority)`,

		// HomeContent
		`CREATE TABLE IF NOT EXISTS home_content (
			id TEXT PRIMARY KEY,
			key TEXT UNIQUE,
			title TEXT,
			content TEXT,
			enabled BOOLEAN DEFAULT true,
			priority INTEGER DEFAULT 0,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_home_content_key ON home_content(key)`,
		`CREATE INDEX IF NOT EXISTS idx_home_content_enabled ON home_content(enabled)`,
	}

	for _, sql := range sqlStatements {
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to execute SQL: %w\nSQL: %s", err, sql)
		}
	}

	log.Println("✅ Tables created successfully!")
	return nil
}

// InsertSampleData 插入示例数据
func InsertSampleData(db *gorm.DB) error {
	log.Println("📦 Inserting initial data...")
	now := time.Now()

	// 1. 创建管理员账户
	if err := insertAdmin(db, now); err != nil {
		log.Printf("⚠️  Warning: Failed to insert admin: %v", err)
	} else {
		log.Println("   ✓ Admin account created (email: admin@z26b.com, password: admin123)")
	}

	// 2. 创建测试用户
	if err := insertTestUser(db, now); err != nil {
		log.Printf("⚠️  Warning: Failed to insert test user: %v", err)
	} else {
		log.Println("   ✓ Test user created (openid: oTest_dev_openid_001)")
	}

	// 3. 创建分类
	categories, err := insertCategories(db, now)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to insert categories: %v", err)
	} else {
		log.Println("   ✓ Categories created")
	}

	// 4. 创建标签
	tags, err := insertTags(db, now)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to insert tags: %v", err)
	} else {
		log.Println("   ✓ Tags created")
	}

	// 5. 创建示例商品
	if len(categories) > 0 && len(tags) > 0 {
		if err := insertSampleProduct(db, now, categories[0].id, tags[0].id); err != nil {
			log.Printf("⚠️  Warning: Failed to insert sample product: %v", err)
		} else {
			log.Println("   ✓ Sample product created")
		}
	}

	// 6. 创建轮播图
	if err := insertSwiper(db, now); err != nil {
		log.Printf("⚠️  Warning: Failed to insert swiper: %v", err)
	} else {
		log.Println("   ✓ Swiper banner created")
	}

	// 7. 创建首页内容
	if err := insertHomeContent(db, now); err != nil {
		log.Printf("⚠️  Warning: Failed to insert home content: %v", err)
	} else {
		log.Println("   ✓ Home content created")
	}

	// 8. 创建优惠券
	if err := insertCoupon(db, now); err != nil {
		log.Printf("⚠️  Warning: Failed to insert coupon: %v", err)
	} else {
		log.Println("   ✓ Welcome coupon created (code: WELCOME10)")
	}

	// 9. 创建促销活动
	if err := insertPromotion(db, now); err != nil {
		log.Printf("⚠️  Warning: Failed to insert promotion: %v", err)
	} else {
		log.Println("   ✓ Promotion created")
	}

	return nil
}

// insertAdmin 插入管理员
func insertAdmin(db *gorm.DB, now time.Time) error {
	adminID := generateUUID()
	adminSQL := fmt.Sprintf(`
		INSERT INTO admin (id, email, password, username, role, status, created_at, updated_at)
		VALUES ('%s', 'admin@z26b.com', '%s', 'Administrator', 'admin', 'active', '%s', '%s')
	`, adminID, hashPassword("admin123"), now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(adminSQL).Error
}

// insertTestUser 插入测试用户
func insertTestUser(db *gorm.DB, now time.Time) error {
	testUserID := generateUUID()
	testUserSQL := fmt.Sprintf(`
		INSERT INTO "user" (id, open_id, nick_name, avatar, created_at, updated_at)
		VALUES ('%s', 'oTest_dev_openid_001', '测试用户', '', '%s', '%s')
	`, testUserID, now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(testUserSQL).Error
}

// Category 分类结构
type Category struct {
	id       string
	name     string
	icon     string
	parentID string
	sort     int
}

// insertCategories 插入分类
func insertCategories(db *gorm.DB, now time.Time) ([]Category, error) {
	categories := []Category{
		{generateUUID(), "电子产品", "📱", "", 1},
		{generateUUID(), "服装鞋包", "👔", "", 2},
		{generateUUID(), "食品饮料", "🍔", "", 3},
		{generateUUID(), "家居生活", "🏠", "", 4},
		{generateUUID(), "运动户外", "⚽", "", 5},
	}

	for _, cat := range categories {
		catSQL := fmt.Sprintf(`
			INSERT INTO category (id, name, icon, image, parent_id, sort, created_at, updated_at)
			VALUES ('%s', '%s', '%s', '', '%s', %d, '%s', '%s')
		`, cat.id, cat.name, cat.icon, cat.parentID, cat.sort, now.Format(time.RFC3339), now.Format(time.RFC3339))

		if err := db.Exec(catSQL).Error; err != nil {
			return categories, err
		}
	}

	return categories, nil
}

// Tag 标签结构
type Tag struct {
	id          string
	name        string
	description string
	color       string
	sortOrder   int
}

// insertTags 插入标签
func insertTags(db *gorm.DB, now time.Time) ([]Tag, error) {
	tags := []Tag{
		{generateUUID(), "热销", "热销商品", "#FF6B6B", 1},
		{generateUUID(), "新品", "新品上市", "#4ECDC4", 2},
		{generateUUID(), "特价", "特价促销", "#FFD93D", 3},
		{generateUUID(), "推荐", "精选推荐", "#6C5CE7", 4},
		{generateUUID(), "限时", "限时优惠", "#FF6348", 5},
	}

	for _, tag := range tags {
		tagSQL := fmt.Sprintf(`
			INSERT INTO tag (id, name, description, color, sort_order, status, created_at, updated_at)
			VALUES ('%s', '%s', '%s', '%s', %d, 'active', '%s', '%s')
		`, tag.id, tag.name, tag.description, tag.color, tag.sortOrder, now.Format(time.RFC3339), now.Format(time.RFC3339))

		if err := db.Exec(tagSQL).Error; err != nil {
			return tags, err
		}
	}

	return tags, nil
}

// insertSampleProduct 插入示例商品
func insertSampleProduct(db *gorm.DB, now time.Time, categoryID, tagID string) error {
	spuID := generateUUID()
	// 示例商品价格
	minPrice := 99.99
	maxPrice := 99.99

	spuSQL := fmt.Sprintf(`
		INSERT INTO spu (id, name, detail, cover_image, swipe_images, category_id, min_price, max_price, status, priority, owner, created_at, updated_at, created_by, updated_by, "_openid")
		VALUES ('%s', '示例商品', '这是一个示例商品，用于测试系统功能。', '', '[]', '%s', %.2f, %.2f, 'ENABLED', 100, 'system', %d, %d, 'system', 'system', '')
	`, spuID, categoryID, minPrice, maxPrice, now.Unix(), now.Unix())

	if err := db.Exec(spuSQL).Error; err != nil {
		return err
	}

	// 创建SKU
	skuID := generateUUID()
	skuSQL := fmt.Sprintf(`
		INSERT INTO sku (id, "SPUID", description, image, price, count, owner, created_at, updated_at, created_by, updated_by, "_openid")
		VALUES ('%s', '%s', '默认规格', '', %.2f, 100, 'system', %d, %d, 'system', 'system', '')
	`, skuID, spuID, minPrice, now.Unix(), now.Unix())

	if err := db.Exec(skuSQL).Error; err != nil {
		return err
	}

	// 关联标签
	spuTagID := generateUUID()
	spuTagSQL := fmt.Sprintf(`
		INSERT INTO spu_tag (id, spu_id, tag_id, created_at)
		VALUES ('%s', '%s', '%s', '%s')
	`, spuTagID, spuID, tagID, now.Format(time.RFC3339))

	return db.Exec(spuTagSQL).Error
}

// insertSwiper 插入轮播图
func insertSwiper(db *gorm.DB, now time.Time) error {
	swiperID := generateUUID()
	swiperSQL := fmt.Sprintf(`
		INSERT INTO swiper (id, images, title, link, priority, created_at, updated_at)
		VALUES ('%s', '[]', '欢迎使用Z26B商城系统', '', 1, '%s', '%s')
	`, swiperID, now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(swiperSQL).Error
}

// insertHomeContent 插入首页内容
func insertHomeContent(db *gorm.DB, now time.Time) error {
	homeContents := []struct {
		key     string
		title   string
		content string
	}{
		{"main", "欢迎光临", "<h1>欢迎使用Z26B商城系统</h1><p>这是一个功能完整的电商平台。</p>"},
		{"notice", "系统公告", "<p>系统正在正常运行中，如有问题请联系管理员。</p>"},
	}

	for _, hc := range homeContents {
		homeContentID := generateUUID()
		homeContentSQL := fmt.Sprintf(`
			INSERT INTO home_content (id, key, title, content, enabled, priority, created_at, updated_at)
			VALUES ('%s', '%s', '%s', '%s', true, 0, '%s', '%s')
		`, homeContentID, hc.key, hc.title, hc.content, now.Format(time.RFC3339), now.Format(time.RFC3339))

		if err := db.Exec(homeContentSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

// insertCoupon 插入优惠券
func insertCoupon(db *gorm.DB, now time.Time) error {
	couponID := generateUUID()
	validFrom := now.Unix()
	validUntil := now.AddDate(0, 1, 0).Unix() // 1个月后
	couponSQL := fmt.Sprintf(`
		INSERT INTO coupon (id, code, discount_type, discount_value, min_amount, max_amount, usage_limit, usage_count, status, valid_from, valid_until, created_at, updated_at)
		VALUES ('%s', 'WELCOME10', 'percentage', 10.00, 50.00, 1000.00, 100, 0, 'active', %d, %d, '%s', '%s')
	`, couponID, validFrom, validUntil, now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(couponSQL).Error
}

// insertPromotion 插入促销活动
func insertPromotion(db *gorm.DB, now time.Time) error {
	promotionID := generateUUID()
	validFrom := now.Unix()
	validUntil := now.AddDate(0, 1, 0).Unix() // 1个月后
	promotionSQL := fmt.Sprintf(`
		INSERT INTO promotion (id, title, promotion_code, promotion_sub_code, tag, description, tag_text, promotion_status, min_amount, valid_from, valid_until, created_at, updated_at)
		VALUES ('%s', '新用户专享', 'NEW_USER', 'DISCOUNT', 'new', '新用户首单立减', '{"text":"新人专享","color":"#FF6B6B"}', 1, 0, %d, %d, '%s', '%s')
	`, promotionID, validFrom, validUntil, now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(promotionSQL).Error
}

// PrintSummary 打印初始化摘要
func PrintSummary() {
	log.Println("\n🎉 Database initialization completed successfully!")
	log.Println("\n📋 Summary:")
	log.Println("   • Admin: admin@z26b.com / admin123")
	log.Println("   • Test User: oTest_dev_openid_001")
	log.Println("   • 5 Categories")
	log.Println("   • 5 Tags")
	log.Println("   • 1 Sample Product with SKU")
	log.Println("   • 1 Swiper Banner")
	log.Println("   • 2 Home Contents")
	log.Println("   • 1 Coupon (WELCOME10)")
	log.Println("   • 1 Promotion")
	log.Println("\n✨ You can now start the server with: go run main.go")
}
