package internal

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"gorm.io/datatypes"
)

// ============================================
// 管理员
// ============================================

type Admin struct {
	ID          string     `gorm:"primaryKey" json:"id"`
	Email       string     `gorm:"uniqueIndex" json:"email"`
	Password    string     `json:"-"` // 密码不输出到JSON
	Username    string     `json:"username"`
	Role        string     `gorm:"default:admin" json:"role"`
	Status      string     `gorm:"default:active" json:"status"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (Admin) TableName() string { return "admin" }

// ============================================
// 用户
// ============================================

type User struct {
	ID        string    `gorm:"primaryKey" json:"_id"`
	OpenID    string    `gorm:"uniqueIndex" json:"openid"`
	NickName  string    `json:"nickName"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string { return "user" }

// ============================================
// 商品
// ============================================

type SPU struct {
	ID          string         `gorm:"primaryKey" json:"_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Detail      string         `gorm:"column:detail" json:"detail"`
	CoverImage  string         `gorm:"column:cover_image" json:"cover_image"`
	SwipeImages datatypes.JSON `gorm:"column:swipe_images;type:json" json:"swiper_images"`
	CategoryID  string         `gorm:"column:category_id" json:"categoryId"`
	Category    *Category      `gorm:"foreignKey:CategoryID;references:ID" json:"category,omitempty"`
	Tags        []Tag          `gorm:"-" json:"tags,omitempty"`
	MinPrice    float64        `gorm:"column:min_price" json:"minPrice"`
	MaxPrice    float64        `gorm:"column:max_price" json:"maxPrice"`
	Status      string         `gorm:"column:status" json:"status"`
	Priority    int            `gorm:"column:priority" json:"priority"`
	Owner       string         `gorm:"column:owner" json:"owner"`
	CreatedAt   int64          `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   int64          `gorm:"column:updated_at" json:"updatedAt"`
	CreatedBy   string         `gorm:"column:created_by" json:"createBy"`
	UpdatedBy   string         `gorm:"column:updated_by" json:"updateBy"`
	OpenID      string         `gorm:"column:_openid" json:"_openid"`
}

func (SPU) TableName() string { return "spu" }

type SKU struct {
	ID          string  `gorm:"primaryKey" json:"_id"`
	SPUID       string  `gorm:"column:SPUID" json:"spuId"`
	SPU         *SPU    `gorm:"foreignKey:SPUID;references:ID" json:"spu,omitempty"`
	Description string  `gorm:"column:description" json:"description"`
	Image       string  `gorm:"column:image" json:"image"`
	Price       float64 `gorm:"column:price" json:"price"`
	Count       int     `gorm:"column:count" json:"count"`
	Owner       string  `gorm:"column:owner" json:"owner"`
	CreatedAt   int64   `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt   int64   `gorm:"column:updated_at" json:"updatedAt"`
	CreatedBy   string  `gorm:"column:created_by" json:"createBy"`
	UpdatedBy   string  `gorm:"column:updated_by" json:"updateBy"`
	OpenID      string  `gorm:"column:_openid" json:"_openid"`
}

func (SKU) TableName() string { return "sku" }

type Category struct {
	ID        string    `gorm:"primaryKey" json:"_id"`
	Name      string    `json:"name"`
	Icon      string    `json:"icon"`
	Image     string    `json:"image"`
	ParentID  string    `json:"parentId"`
	Sort      int       `json:"sort"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Category) TableName() string { return "category" }

type Tag struct {
	ID          string    `gorm:"primaryKey" json:"_id"`
	Name        string    `gorm:"uniqueIndex" json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	SortOrder   int       `json:"sortOrder"`
	Status      string    `gorm:"default:active" json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Tag) TableName() string { return "tag" }

type SPUTag struct {
	ID        string    `gorm:"primaryKey" json:"_id"`
	SPUID     string    `gorm:"column:spu_id;index" json:"spuId"`
	TagID     string    `gorm:"column:tag_id;index" json:"tagId"`
	Tag       *Tag      `gorm:"foreignKey:TagID;references:ID" json:"tag,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

func (SPUTag) TableName() string { return "spu_tag" }

// ============================================
// 购物车
// ============================================

type CartItem struct {
	ID         string    `gorm:"primaryKey" json:"_id"`
	UserID     string    `gorm:"column:user_id" json:"userId"`
	SKUID      string    `gorm:"column:sku_id" json:"skuId"`
	SKU        *SKU      `gorm:"foreignKey:SKUID;references:ID" json:"sku,omitempty"`
	Quantity   int       `json:"quantity"`
	IsSelected bool      `gorm:"column:is_selected" json:"isSelected"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (CartItem) TableName() string { return "cart_item" }

// ============================================
// 地址
// ============================================

type Address struct {
	ID            string    `gorm:"primaryKey" json:"_id"`
	UserID        string    `json:"userId"`
	Name          string    `json:"name"`
	Phone         string    `json:"phone"`
	CountryName   string    `json:"countryName"`
	CountryCode   string    `json:"countryCode"`
	ProvinceName  string    `json:"provinceName"`
	ProvinceCode  string    `json:"provinceCode"`
	CityName      string    `json:"cityName"`
	CityCode      string    `json:"cityCode"`
	DistrictName  string    `json:"districtName"`
	DistrictCode  string    `json:"districtCode"`
	DetailAddress string    `json:"detailAddress"`
	IsDefault     int       `json:"isDefault"`
	AddressTag    string    `json:"addressTag"`
	Latitude      string    `json:"latitude"`
	Longitude     string    `json:"longitude"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (Address) TableName() string { return "address" }

// ============================================
// 订单
// ============================================

const (
	OrderStatusToPay         = "TO_PAY"
	OrderStatusToSend        = "TO_SEND"
	OrderStatusToReceive     = "TO_RECEIVE"
	OrderStatusFinished      = "FINISHED"
	OrderStatusCanceled      = "CANCELED"
	OrderStatusReturnApplied = "RETURN_APPLIED"
	OrderStatusReturnRefused = "RETURN_REFUSED"
	OrderStatusReturnFinish  = "RETURN_FINISH"
)

type Order struct {
	ID            string         `gorm:"primaryKey" json:"_id"`
	UserID        string         `json:"userId"`
	Status        string         `json:"status"`
	DeliveryInfo  datatypes.JSON `gorm:"type:json" json:"delivery_info"`
	Items         []OrderItem    `gorm:"foreignKey:OrderID;references:ID" json:"items,omitempty"`
	TotalPrice    float64        `json:"totalPrice"`
	DiscountPrice float64        `json:"discountPrice"`
	FinalPrice    float64        `json:"finalPrice"`
	Remarks       string         `json:"remarks"`
	CreatedAt     int64          `json:"createdAt"`
	UpdatedAt     int64          `json:"updatedAt"`
}

func (Order) TableName() string { return "order" }

type OrderItem struct {
	ID        string    `gorm:"primaryKey" json:"_id"`
	OrderID   string    `gorm:"column:order_id" json:"orderId"`
	SKUID     string    `gorm:"column:sku_id" json:"skuId"`
	SKU       *SKU      `gorm:"foreignKey:SKUID;references:ID" json:"sku,omitempty"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (OrderItem) TableName() string { return "order_item" }

// ============================================
// 评论
// ============================================

type Comment struct {
	ID               string         `gorm:"primaryKey" json:"_id"`
	SPUID            string         `gorm:"column:spu_id" json:"spuId"`
	SKUID            string         `gorm:"column:sku_id" json:"skuId"`
	UserID           string         `gorm:"column:user_id" json:"userId"`
	UserName         string         `json:"userName"`
	UserHeadURL      string         `json:"userHeadUrl"`
	CommentContent   string         `json:"commentContent"`
	CommentScore     int            `json:"commentScore"`
	CommentResources datatypes.JSON `gorm:"type:json" json:"commentResources"`
	IsAnonymity      bool           `json:"isAnonymity"`
	SellerReply      string         `json:"sellerReply"`
	CreatedAt        int64          `json:"commentTime"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

func (Comment) TableName() string { return "comment" }

// ============================================
// 优惠券 & 促销
// ============================================

type Coupon struct {
	ID            string    `gorm:"primaryKey" json:"_id"`
	Code          string    `json:"code"`
	DiscountType  string    `json:"discountType"`
	DiscountValue float64   `json:"discountValue"`
	MinAmount     float64   `json:"minAmount"`
	MaxAmount     float64   `json:"maxAmount"`
	UsageLimit    int       `json:"usageLimit"`
	UsageCount    int       `json:"usageCount"`
	Status        string    `json:"status"`
	ValidFrom     int64     `json:"validFrom"`
	ValidUntil    int64     `json:"validUntil"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

func (Coupon) TableName() string { return "coupon" }

type Promotion struct {
	ID               string         `gorm:"primaryKey" json:"_id"`
	Title            string         `json:"title"`
	PromotionCode    string         `json:"promotionCode"`
	PromotionSubCode string         `json:"promotionSubCode"`
	Tag              string         `json:"tag"`
	Description      string         `json:"description"`
	TagText          datatypes.JSON `gorm:"type:json" json:"tagText"`
	PromotionStatus  int            `json:"promotionStatus"`
	MinAmount        float64        `json:"doorSillRemain"`
	ValidFrom        int64          `json:"validFrom"`
	ValidUntil       int64          `json:"validUntil"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
}

func (Promotion) TableName() string { return "promotion" }

// ============================================
// 轮播图
// ============================================

type Swiper struct {
	ID        string         `gorm:"primaryKey" json:"_id"`
	Images    datatypes.JSON `gorm:"type:json" json:"images"`
	Title     string         `json:"title"`
	Link      string         `json:"link"`
	Priority  int            `json:"priority"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (Swiper) TableName() string { return "swiper" }

// ============================================
// 推荐商品（独立关联表）
// ============================================

type RecommendedProduct struct {
	ID        string         `gorm:"primaryKey" json:"_id"`
	SPUID     string         `gorm:"column:spu_id;index" json:"spuId"`
	SPU       *SPU           `gorm:"foreignKey:SPUID;references:ID" json:"spu,omitempty"`
	Tags      datatypes.JSON `gorm:"type:json" json:"tags"`
	Priority  int            `gorm:"column:priority;default:0" json:"priority"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
}

func (RecommendedProduct) TableName() string { return "recommended_product" }

// ============================================
// 首页内容（富文本）
// ============================================

type HomeContent struct {
	ID        string    `gorm:"primaryKey" json:"_id"`
	Key       string    `gorm:"uniqueIndex" json:"key"`   // 内容标识，如 "main", "promotion", "notice"
	Title     string    `json:"title"`                    // 标题
	Content   string    `gorm:"type:text" json:"content"` // 富文本内容
	Enabled   bool      `gorm:"default:true" json:"enabled"`
	Priority  int       `gorm:"default:0" json:"priority"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (HomeContent) TableName() string { return "home_content" }

// ============================================
// 工具类型
// ============================================

type JSONMap map[string]interface{}

func (jm JSONMap) Value() (driver.Value, error) {
	return json.Marshal(jm)
}

func (jm *JSONMap) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &jm)
}
