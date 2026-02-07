package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"z26b-backend/handlers"
	"z26b-backend/handlers/admin"
	"z26b-backend/handlers/miniprogram"
	"z26b-backend/internal"
	"z26b-backend/middleware"
	admin_services "z26b-backend/services/admin_services"
	"z26b-backend/services/crm"
	miniprogram_services "z26b-backend/services/miniprogram"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default configuration")
	}
}

func main() {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database
	db := internal.InitDB()
	if db == nil {
		internal.GlobalLogger.Error("Failed to initialize database", fmt.Errorf("database initialization failed"))
		log.Fatal("Failed to initialize database")
	}

	// Initialize default admin account
	if err := admin.InitDefaultAdmin(db); err != nil {
		internal.GlobalLogger.Warn("Failed to create default admin", map[string]interface{}{"error": err.Error()})
	}

	// Initialize MinIO storage (optional)
	minioConfig := internal.GetDefaultMinIOConfig()
	if err := internal.InitMinIO(minioConfig); err != nil {
		internal.GlobalLogger.Warn("MinIO not available", map[string]interface{}{"error": err.Error()})
	} else {
		internal.GlobalLogger.Info("MinIO storage initialized")
	}

	// Create gin router
	router := gin.New() // 使用gin.New()而不是gin.Default()，以便自定义中间件
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RateLimitMiddleware(100)) // 每秒100个请求
	router.Use(middleware.RequestTimeoutMiddleware(30 * time.Second))

	// Initialize services
	goodsService := miniprogram_services.NewGoodsService(db)
	adminGoodsService := admin_services.NewAdminGoodsService(db)
	userService := miniprogram_services.NewUserService(db)
	addressService := miniprogram_services.NewAddressService(db)
	cartService := miniprogram_services.NewCartService(db)
	orderService := miniprogram_services.NewOrderService(db)
	commentService := miniprogram_services.NewCommentService(db)
	wechatService := miniprogram_services.NewWechatService(db)
	adminCategoryService := admin_services.NewAdminCategoryService(db)

	// Initialize CRM services
	crmEventService := crm.NewCRMEventService(db)
	customerStatsService := crm.NewCustomerStatsService(db)
	productStatsService := crm.NewProductStatsService(db)

	// Initialize handlers
	mpHandler := miniprogram.NewHandler(goodsService, userService, cartService, orderService, commentService, wechatService, crmEventService, db)
	adminHandler := admin.NewHandler(adminGoodsService, adminCategoryService, crmEventService, customerStatsService, productStatsService, db)
	addressHandler := handlers.NewAddressHandler(addressService)

	// ====== 小程序端 API ======
	initMiniProgramRoutes(router, mpHandler, addressHandler)

	// ====== 管理后台 API ======
	initAdminRoutes(router, adminHandler)

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "message": "Server is running"})
	})

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		internal.GlobalLogger.Error("Failed to start server", err)
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initMiniProgramRoutes(router *gin.Engine, h *miniprogram.Handler, addressHandler *handlers.AddressHandler) {
	// 小程序端API都需要用户ID
	api := router.Group("/api")
	// 添加小程序端用户认证中间件
	api.Use(middleware.OptionalMiniProgramAuthMiddleware())

	// WeChat routes

	// WeChat routes
	wechat := api.Group("/wechat")
	{
		wechat.POST("/login", h.WxLogin)
		// TODO: 实现微信支付相关功能
		// wechat.POST("/updateUserInfo", h.UpdateWxUserInfo)
		// wechat.POST("/pay/create", h.CreateWxPayOrder)
		// wechat.POST("/pay/notify", h.WxPayNotifyHandler)
		// wechat.GET("/pay/query/:orderId", h.WxPayQuery)
		// wechat.POST("/pay/refund", h.WxRefund)
	}

	// User routes
	user := api.Group("/user")
	{
		user.GET("/info", h.GetUserInfo)
		user.PUT("/info", h.UpdateUserInfo)
	}

	// SKU routes
	sku := api.Group("/sku")
	{
		sku.GET("/list/:spuId", h.GetSKUsBySpuId)
		sku.GET("/:id", h.GetSKU)
	}

	// Goods routes
	goods := api.Group("/goods")
	{
		goods.GET("/list", h.GetGoodsList)
		goods.GET("/category/list", h.GetCategories)
		goods.GET("/search", h.SearchGoods)
		goods.GET("/:id", h.GetGood)
		goods.GET("/:id/comments", h.GetGoodsComments)
	}

	// Cart routes
	cart := api.Group("/cart")
	{
		cart.GET("/items", h.GetCartItems)
		cart.POST("/add", h.AddToCart)
		cart.PUT("/update/:id", h.UpdateCartItem)
		cart.DELETE("/remove/:id", h.RemoveFromCart)
		cart.POST("/clear", h.ClearCart)
	}

	// Address routes
	address := api.Group("/address")
	{
		address.GET("/list", addressHandler.GetAddressList)
		address.GET("/:id", addressHandler.GetAddress)
		address.POST("/create", addressHandler.CreateAddress)
		address.PUT("/update/:id", addressHandler.UpdateAddress)
		address.DELETE("/:id", addressHandler.DeleteAddress)
		address.POST("/setDefault/:id", addressHandler.SetDefaultAddress)
	}

	// Order routes
	order := api.Group("/order")
	{
		order.GET("/list", h.GetOrderList)
		order.GET("/:id", h.GetOrderDetail)
		order.POST("/create", h.CreateOrder)
		order.PUT("/cancel/:id", h.CancelOrder)
		order.POST("/confirm/:id", h.ConfirmReceipt)
		order.POST("/submit-comment/:id", h.SubmitComment)
	}

	// Comment routes
	comment := api.Group("/comment")
	{
		comment.GET("/list/:spuId", h.GetGoodsComments)
		comment.POST("/submit", h.SubmitComment)
	}

	// Home routes
	home := api.Group("/home")
	{
		home.GET("/swiper", h.GetHomeSwiper)
		home.GET("/content", h.GetHomeContent)
		home.GET("/categories", h.GetHomeCategories)
		home.GET("/promotions", h.GetPromotions)
	}
}

func initAdminRoutes(router *gin.Engine, h *admin.Handler) {
	adminGroup := router.Group("/admin/api")
	{
		// Auth routes (no auth required)
		adminGroup.POST("/login", h.AdminLogin)

		// Protected routes
		protected := adminGroup.Group("")
		protected.Use(internal.AdminAuthMiddleware())
		{
			// Auth
			protected.GET("/profile", h.AdminGetProfile)
			protected.PUT("/password", h.AdminChangePassword)

			// Dashboard & Stats
			protected.GET("/stats", h.AdminGetStats)
			protected.GET("/stats/sales", h.AdminGetSalesStats)
			protected.GET("/stats/top-products", h.AdminGetTopProducts)
			protected.GET("/stats/order-status", h.AdminGetOrderStatusStats)

			// Upload
			protected.POST("/upload/image", h.AdminUploadImage)
			protected.POST("/upload/images", h.AdminUploadImages)
			protected.DELETE("/upload/image", h.AdminDeleteImage)

			// Products (SPU)
			protected.GET("/products", h.AdminGetProducts)
			protected.POST("/products", h.AdminCreateProduct)
			protected.GET("/products/:id", h.AdminGetProduct)
			protected.PUT("/products/:id", h.AdminUpdateProduct)
			protected.DELETE("/products/:id", h.AdminDeleteProduct)
			protected.PUT("/products/:id/toggle-status", h.AdminToggleProductStatus)
			protected.GET("/products/:id/skus", h.AdminGetSKUs)

			// SKU
			protected.POST("/skus", h.AdminCreateSKU)
			protected.PUT("/skus/:id", h.AdminUpdateSKU)
			protected.DELETE("/skus/:id", h.AdminDeleteSKU)

			// Orders
			protected.GET("/orders", h.AdminGetOrders)
			protected.GET("/orders/:id", h.AdminGetOrder)
			protected.PUT("/orders/:id/ship", h.AdminShipOrder)
			protected.PUT("/orders/:id/refund", h.AdminRefundOrder)
			protected.PUT("/orders/:id/status", h.AdminUpdateOrderStatus)

			// Users
			protected.GET("/users", h.AdminGetUsers)
			protected.GET("/users/:id", h.AdminGetUser)
			protected.GET("/users/:id/orders", h.AdminGetUserOrders)
			protected.POST("/users/test", h.AdminCreateTestUser)

			// Categories
			protected.GET("/categories", h.AdminGetCategories)
			protected.GET("/categories/:id", h.AdminGetCategory)
			protected.POST("/categories", h.AdminCreateCategory)
			protected.PUT("/categories/:id", h.AdminUpdateCategory)
			protected.DELETE("/categories/:id", h.AdminDeleteCategory)

			// Tags
			protected.GET("/tags", h.AdminGetTags)
			protected.GET("/tags/:id", h.AdminGetTag)
			protected.POST("/tags", h.AdminCreateTag)
			protected.PUT("/tags/:id", h.AdminUpdateTag)
			protected.DELETE("/tags/:id", h.AdminDeleteTag)

			// Home Config - Banners
			protected.GET("/home-config/banners", h.AdminGetBanners)
			protected.POST("/home-config/banners", h.AdminCreateBanner)
			protected.PUT("/home-config/banners/:id", h.AdminUpdateBanner)
			protected.DELETE("/home-config/banners/:id", h.AdminDeleteBanner)
			protected.POST("/home-config/banners/reorder", h.AdminReorderBanners)

			// Home Config - Recommended Products
			protected.GET("/home-config/recommended", h.AdminGetRecommendedProducts)
			protected.POST("/home-config/recommended", h.AdminAddRecommendedProduct)
			protected.PUT("/home-config/recommended/:id", h.AdminUpdateRecommendedProduct)
			protected.DELETE("/home-config/recommended/:id", h.AdminRemoveRecommendedProduct)
			protected.POST("/home-config/recommended/reorder", h.AdminReorderRecommendedProducts)
			protected.GET("/home-config/products/search", h.AdminSearchProductsForRecommend)

			// Home Config - Rich Text Contents
			protected.GET("/home-config/contents", h.AdminGetHomeContents)
			protected.GET("/home-config/contents/:key", h.AdminGetHomeContent)
			protected.POST("/home-config/contents", h.AdminSaveHomeContent)
			protected.DELETE("/home-config/contents/:key", h.AdminDeleteHomeContent)

			// CRM - Dashboard
			protected.GET("/crm/dashboard", h.AdminGetCRMDashboard)

			// CRM - Product Analysis
			protected.GET("/crm/product/overview", h.AdminGetProductAnalysisOverview)
			protected.GET("/crm/product/list", h.AdminGetProductStatsList)
			protected.GET("/crm/product/top", h.AdminGetTopProductStats)
			protected.POST("/crm/product/:id/refresh", h.AdminRefreshProductStats)
			protected.POST("/crm/product/refresh-all", h.AdminRefreshAllProductStats)

			// CRM - Customer Analysis
			protected.GET("/crm/customer/overview", h.AdminGetCustomerAnalysisOverview)
			protected.GET("/crm/customer/list", h.AdminGetCustomerStatsList)
			protected.GET("/crm/customer/top", h.AdminGetTopCustomers)
			protected.GET("/crm/customer/level-distribution", h.AdminGetCustomerLevelDistribution)
			protected.POST("/crm/customer/:id/refresh", h.AdminRefreshCustomerStats)
			protected.POST("/crm/customer/refresh-all", h.AdminRefreshAllCustomerStats)

			// CRM - Events
			protected.GET("/crm/events/stats", h.AdminGetCRMEventStats)
			protected.GET("/crm/events/daily", h.AdminGetDailyEventStats)
			protected.GET("/crm/events/user/:userId", h.AdminGetCRMEventsByUser)
			protected.GET("/crm/events/product/:spuId", h.AdminGetCRMEventsBySPU)
		}
	}
}
