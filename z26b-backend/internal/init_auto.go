package internal

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AutoInitDB 自动检测并初始化数据库
// 如果数据库未初始化（admin表不存在或为空），则自动初始化
func AutoInitDB() error {
	db := GetDB()
	if db == nil {
		return fmt.Errorf("database connection not established")
	}

	// 检查是否已初始化
	if IsDatabaseInitialized(db) {
		log.Println("✅ Database already initialized, skipping...")
		return nil
	}

	log.Println("📦 Database not initialized, starting initialization...")

	// 创建表
	if err := CreateAllTables(db); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	// 插入初始数据
	if err := InsertInitialData(db); err != nil {
		return fmt.Errorf("failed to insert initial data: %w", err)
	}

	log.Println("🎉 Database initialization completed!")
	return nil
}

// IsDatabaseInitialized 检查数据库是否已初始化
func IsDatabaseInitialized(db *gorm.DB) bool {
	// 检查 admin 表是否存在
	var exists bool
	dbType := os.Getenv("DB_TYPE")
	if dbType == "" {
		dbType = "postgres"
	}

	var checkSQL string
	if dbType == "postgres" {
		checkSQL = `SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = 'admin'
		)`
	} else {
		// SQLite
		checkSQL = `SELECT COUNT(*) > 0 FROM sqlite_master WHERE type='table' AND name='admin'`
	}

	if err := db.Raw(checkSQL).Scan(&exists).Error; err != nil {
		log.Printf("Error checking if database is initialized: %v", err)
		return false
	}

	if !exists {
		return false
	}

	// 检查 admin 表是否有数据
	var count int64
	if err := db.Table("admin").Count(&count).Error; err != nil {
		log.Printf("Error counting admin records: %v", err)
		return false
	}

	return count > 0
}

// CreateAllTables 创建所有表
func CreateAllTables(db *gorm.DB) error {
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
			last_login_at TIMESTAMP,
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

		// CRM Event (CRM事件表)
		`CREATE TABLE IF NOT EXISTS crm_event (
			id TEXT PRIMARY KEY,
			user_id TEXT,
			event_type TEXT,
			spu_id TEXT,
			sku_id TEXT,
			order_id TEXT,
			amount DECIMAL(10,2),
			extra JSONB,
			ip_address TEXT,
			user_agent TEXT,
			created_at BIGINT
		)`,
		`CREATE INDEX IF NOT EXISTS idx_crm_event_user_id ON crm_event(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_crm_event_event_type ON crm_event(event_type)`,
		`CREATE INDEX IF NOT EXISTS idx_crm_event_spu_id ON crm_event(spu_id)`,
		`CREATE INDEX IF NOT EXISTS idx_crm_event_created_at ON crm_event(created_at)`,

		// Customer Stats (客户统计表)
		`CREATE TABLE IF NOT EXISTS customer_stats (
			id TEXT PRIMARY KEY,
			user_id TEXT UNIQUE,
			total_orders INTEGER DEFAULT 0,
			total_spent DECIMAL(10,2) DEFAULT 0,
			avg_order_value DECIMAL(10,2) DEFAULT 0,
			total_refunds INTEGER DEFAULT 0,
			refund_amount DECIMAL(10,2) DEFAULT 0,
			total_views INTEGER DEFAULT 0,
			total_carts INTEGER DEFAULT 0,
			total_comments INTEGER DEFAULT 0,
			total_shares INTEGER DEFAULT 0,
			last_order_at BIGINT,
			last_active_at BIGINT,
			customer_level TEXT DEFAULT 'normal',
			tags JSONB,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_customer_stats_user_id ON customer_stats(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_customer_stats_customer_level ON customer_stats(customer_level)`,
		`CREATE INDEX IF NOT EXISTS idx_customer_stats_total_spent ON customer_stats(total_spent)`,

		// Product Stats (商品统计表)
		`CREATE TABLE IF NOT EXISTS product_stats (
			id TEXT PRIMARY KEY,
			spu_id TEXT UNIQUE,
			total_views INTEGER DEFAULT 0,
			total_carts INTEGER DEFAULT 0,
			total_sales INTEGER DEFAULT 0,
			total_revenue DECIMAL(10,2) DEFAULT 0,
			total_refunds INTEGER DEFAULT 0,
			refund_amount DECIMAL(10,2) DEFAULT 0,
			total_comments INTEGER DEFAULT 0,
			avg_score DECIMAL(3,2) DEFAULT 0,
			total_shares INTEGER DEFAULT 0,
			conversion_rate DECIMAL(5,4) DEFAULT 0,
			created_at TIMESTAMP,
			updated_at TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_product_stats_spu_id ON product_stats(spu_id)`,
		`CREATE INDEX IF NOT EXISTS idx_product_stats_total_sales ON product_stats(total_sales)`,
		`CREATE INDEX IF NOT EXISTS idx_product_stats_total_revenue ON product_stats(total_revenue)`,
	}

	for _, sql := range sqlStatements {
		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to execute SQL: %w\nSQL: %s", err, sql)
		}
	}

	log.Println("✅ Tables created successfully!")
	return nil
}

// InsertInitialData 插入初始数据
func InsertInitialData(db *gorm.DB) error {
	now := time.Now()
	log.Println("📦 Inserting initial data...")

	// 1. 创建管理员账户
	if err := createDefaultAdmin(db, now); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	} else {
		log.Println("   ✓ Admin account created (email: admin@z26b.com, password: admin123)")
	}

	// 2. 创建测试用户
	if err := createDefaultTestUser(db, now); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	} else {
		log.Println("   ✓ Test user created (openid: oTest_dev_openid_001)")
	}

	// 3. 创建分类
	categoryIDs, err := createDefaultCategories(db, now)
	if err != nil {
		log.Printf("⚠️  Warning: %v", err)
	} else {
		log.Println("   ✓ Categories created")
	}

	// 4. 创建标签
	tagIDs, err := createDefaultTags(db, now)
	if err != nil {
		log.Printf("⚠️  Warning: %v", err)
	} else {
		log.Println("   ✓ Tags created")
	}

	// 5. 创建示例商品
	if len(categoryIDs) > 0 && len(tagIDs) > 0 {
		if err := createSampleProduct(db, now, categoryIDs[0], tagIDs[0]); err != nil {
			log.Printf("⚠️  Warning: %v", err)
		} else {
			log.Println("   ✓ Sample product created")
		}
	}

	// 6. 创建轮播图
	if err := createDefaultSwiper(db, now); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	} else {
		log.Println("   ✓ Swiper banner created")
	}

	// 7. 创建首页内容
	if err := createDefaultHomeContent(db, now); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	} else {
		log.Println("   ✓ Home content created")
	}

	// 8. 创建优惠券
	if err := createDefaultCoupon(db, now); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	} else {
		log.Println("   ✓ Welcome coupon created (code: WELCOME10)")
	}

	// 9. 创建促销活动
	if err := createDefaultPromotion(db, now); err != nil {
		log.Printf("⚠️  Warning: %v", err)
	} else {
		log.Println("   ✓ Promotion created")
	}

	printInitSummary()
	return nil
}

// hashPasswordInternal 加密密码
func hashPasswordInternal(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return password
	}
	return string(hashedPassword)
}

// createDefaultAdmin 创建默认管理员
func createDefaultAdmin(db *gorm.DB, now time.Time) error {
	adminID := GenerateUUID()
	adminSQL := fmt.Sprintf(`
		INSERT INTO admin (id, email, password, username, role, status, created_at, updated_at)
		VALUES ('%s', 'admin@z26b.com', '%s', 'Administrator', 'admin', 'active', '%s', '%s')
		ON CONFLICT (email) DO NOTHING
	`, adminID, hashPasswordInternal("admin123"), now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(adminSQL).Error
}

// createDefaultTestUser 创建默认测试用户
func createDefaultTestUser(db *gorm.DB, now time.Time) error {
	testUserID := GenerateUUID()
	testUserSQL := fmt.Sprintf(`
		INSERT INTO "user" (id, open_id, nick_name, avatar, created_at, updated_at)
		VALUES ('%s', 'oTest_dev_openid_001', '测试用户', '', '%s', '%s')
		ON CONFLICT (open_id) DO NOTHING
	`, testUserID, now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(testUserSQL).Error
}

// createDefaultCategories 创建默认分类
func createDefaultCategories(db *gorm.DB, now time.Time) ([]string, error) {
	categories := []struct {
		name string
		icon string
		sort int
	}{
		{"电子产品", "📱", 1},
		{"服装鞋包", "👔", 2},
		{"食品饮料", "🍔", 3},
		{"家居生活", "🏠", 4},
		{"运动户外", "⚽", 5},
	}

	var ids []string
	for _, cat := range categories {
		id := GenerateUUID()
		ids = append(ids, id)
		catSQL := fmt.Sprintf(`
			INSERT INTO category (id, name, icon, image, parent_id, sort, created_at, updated_at)
			VALUES ('%s', '%s', '%s', '', '', %d, '%s', '%s')
		`, id, cat.name, cat.icon, cat.sort, now.Format(time.RFC3339), now.Format(time.RFC3339))

		if err := db.Exec(catSQL).Error; err != nil {
			return ids, err
		}
	}

	return ids, nil
}

// createDefaultTags 创建默认标签
func createDefaultTags(db *gorm.DB, now time.Time) ([]string, error) {
	tags := []struct {
		name        string
		description string
		color       string
		sortOrder   int
	}{
		{"热销", "热销商品", "#FF6B6B", 1},
		{"新品", "新品上市", "#4ECDC4", 2},
		{"特价", "特价促销", "#FFD93D", 3},
		{"推荐", "精选推荐", "#6C5CE7", 4},
		{"限时", "限时优惠", "#FF6348", 5},
	}

	var ids []string
	for _, tag := range tags {
		id := GenerateUUID()
		ids = append(ids, id)
		tagSQL := fmt.Sprintf(`
			INSERT INTO tag (id, name, description, color, sort_order, status, created_at, updated_at)
			VALUES ('%s', '%s', '%s', '%s', %d, 'active', '%s', '%s')
			ON CONFLICT (name) DO NOTHING
		`, id, tag.name, tag.description, tag.color, tag.sortOrder, now.Format(time.RFC3339), now.Format(time.RFC3339))

		if err := db.Exec(tagSQL).Error; err != nil {
			return ids, err
		}
	}

	return ids, nil
}

// createSampleProduct 创建示例商品
func createSampleProduct(db *gorm.DB, now time.Time, categoryID, tagID string) error {
	spuID := GenerateUUID()
	spuSQL := fmt.Sprintf(`
		INSERT INTO spu (id, name, detail, cover_image, swipe_images, category_id, status, priority, owner, created_at, updated_at, created_by, updated_by, "_openid")
		VALUES ('%s', '示例商品', '这是一个示例商品，用于测试系统功能。', '', '[]', '%s', 'ENABLED', 100, 'system', %d, %d, 'system', 'system', '')
	`, spuID, categoryID, now.Unix(), now.Unix())

	if err := db.Exec(spuSQL).Error; err != nil {
		return err
	}

	// 创建SKU
	skuID := GenerateUUID()
	skuSQL := fmt.Sprintf(`
		INSERT INTO sku (id, "SPUID", description, image, price, count, owner, created_at, updated_at, created_by, updated_by, "_openid")
		VALUES ('%s', '%s', '默认规格', '', 99.99, 100, 'system', %d, %d, 'system', 'system', '')
	`, skuID, spuID, now.Unix(), now.Unix())

	if err := db.Exec(skuSQL).Error; err != nil {
		return err
	}

	// 关联标签
	spuTagID := GenerateUUID()
	spuTagSQL := fmt.Sprintf(`
		INSERT INTO spu_tag (id, spu_id, tag_id, created_at)
		VALUES ('%s', '%s', '%s', '%s')
	`, spuTagID, spuID, tagID, now.Format(time.RFC3339))

	return db.Exec(spuTagSQL).Error
}

// createDefaultSwiper 创建默认轮播图
func createDefaultSwiper(db *gorm.DB, now time.Time) error {
	swiperID := GenerateUUID()
	swiperSQL := fmt.Sprintf(`
		INSERT INTO swiper (id, images, title, link, priority, created_at, updated_at)
		VALUES ('%s', '[]', '欢迎使用Z26B商城系统', '', 1, '%s', '%s')
	`, swiperID, now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(swiperSQL).Error
}

// createDefaultHomeContent 创建默认首页内容
func createDefaultHomeContent(db *gorm.DB, now time.Time) error {
	homeContents := []struct {
		key     string
		title   string
		content string
	}{
		{"main", "欢迎光临", "<h1>欢迎使用Z26B商城系统</h1><p>这是一个功能完整的电商平台。</p>"},
		{"notice", "系统公告", "<p>系统正在正常运行中，如有问题请联系管理员。</p>"},
	}

	for _, hc := range homeContents {
		homeContentID := GenerateUUID()
		homeContentSQL := fmt.Sprintf(`
			INSERT INTO home_content (id, key, title, content, enabled, priority, created_at, updated_at)
			VALUES ('%s', '%s', '%s', '%s', true, 0, '%s', '%s')
			ON CONFLICT (key) DO NOTHING
		`, homeContentID, hc.key, hc.title, hc.content, now.Format(time.RFC3339), now.Format(time.RFC3339))

		if err := db.Exec(homeContentSQL).Error; err != nil {
			return err
		}
	}

	return nil
}

// createDefaultCoupon 创建默认优惠券
func createDefaultCoupon(db *gorm.DB, now time.Time) error {
	couponID := GenerateUUID()
	validFrom := now.Unix()
	validUntil := now.AddDate(0, 1, 0).Unix() // 1个月后
	couponSQL := fmt.Sprintf(`
		INSERT INTO coupon (id, code, discount_type, discount_value, min_amount, max_amount, usage_limit, usage_count, status, valid_from, valid_until, created_at, updated_at)
		VALUES ('%s', 'WELCOME10', 'percentage', 10.00, 50.00, 1000.00, 100, 0, 'active', %d, %d, '%s', '%s')
	`, couponID, validFrom, validUntil, now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(couponSQL).Error
}

// createDefaultPromotion 创建默认促销活动
func createDefaultPromotion(db *gorm.DB, now time.Time) error {
	promotionID := GenerateUUID()
	validFrom := now.Unix()
	validUntil := now.AddDate(0, 1, 0).Unix() // 1个月后
	promotionSQL := fmt.Sprintf(`
		INSERT INTO promotion (id, title, promotion_code, promotion_sub_code, tag, description, tag_text, promotion_status, min_amount, valid_from, valid_until, created_at, updated_at)
		VALUES ('%s', '新用户专享', 'NEW_USER', 'DISCOUNT', 'new', '新用户首单立减', '{"text":"新人专享","color":"#FF6B6B"}', 1, 0, %d, %d, '%s', '%s')
	`, promotionID, validFrom, validUntil, now.Format(time.RFC3339), now.Format(time.RFC3339))

	return db.Exec(promotionSQL).Error
}

// printInitSummary 打印初始化摘要
func printInitSummary() {
	log.Println("\n📋 Initialization Summary:")
	log.Println("   • Admin: admin@z26b.com / admin123")
	log.Println("   • Test User: oTest_dev_openid_001")
	log.Println("   • 5 Categories")
	log.Println("   • 5 Tags")
	log.Println("   • 1 Sample Product with SKU")
	log.Println("   • 1 Swiper Banner")
	log.Println("   • 2 Home Contents")
	log.Println("   • 1 Coupon (WELCOME10)")
	log.Println("   • 1 Promotion")
}
