package handlers

import (
	"net/http"

	"z26b-backend/internal"
	miniprogram_services "z26b-backend/services/miniprogram"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AddressHandler struct {
	addressService miniprogram_services.AddressServiceInterface
}

// NewAddressHandler creates a new address handler
func NewAddressHandler(addressService miniprogram_services.AddressServiceInterface) *AddressHandler {
	return &AddressHandler{addressService: addressService}
}

// GetAddressList retrieves all addresses for the user
func (h *AddressHandler) GetAddressList(c *gin.Context) {
	userID := c.GetString("userID")
	if userID == "" {
		userID = "USER_MOCK"
	}

	addresses, err := h.addressService.GetAddressList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch addresses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": addresses})
}

// GetAddress retrieves a single address by ID
func (h *AddressHandler) GetAddress(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "USER_MOCK"
	}

	address, err := h.addressService.GetAddress(id, userID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": address})
}

// CreateAddress creates a new address
func (h *AddressHandler) CreateAddress(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Phone         string `json:"phone" binding:"required"`
		CountryName   string `json:"countryName"`
		CountryCode   string `json:"countryCode"`
		ProvinceName  string `json:"provinceName"`
		ProvinceCode  string `json:"provinceCode"`
		CityName      string `json:"cityName"`
		CityCode      string `json:"cityCode"`
		DistrictName  string `json:"districtName"`
		DistrictCode  string `json:"districtCode"`
		DetailAddress string `json:"detailAddress"`
		Address       string `json:"address"` // 简单地址字段，兼容前端
		AddressTag    string `json:"addressTag"`
		Latitude      string `json:"latitude"`
		Longitude     string `json:"longitude"`
		IsDefault     int    `json:"isDefault"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请填写收货人姓名和手机号", "detail": err.Error()})
		return
	}

	userID := c.GetString("userID")
	if userID == "" {
		userID = "USER_MOCK"
	}

	// 兼容简单地址格式
	detailAddress := req.DetailAddress
	if detailAddress == "" && req.Address != "" {
		detailAddress = req.Address
	}

	address := internal.Address{
		ID:            internal.GenerateID(),
		UserID:        userID,
		Name:          req.Name,
		Phone:         req.Phone,
		CountryName:   req.CountryName,
		CountryCode:   req.CountryCode,
		ProvinceName:  req.ProvinceName,
		ProvinceCode:  req.ProvinceCode,
		CityName:      req.CityName,
		CityCode:      req.CityCode,
		DistrictName:  req.DistrictName,
		DistrictCode:  req.DistrictCode,
		DetailAddress: detailAddress,
		AddressTag:    req.AddressTag,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		IsDefault:     req.IsDefault,
	}

	if err := h.addressService.CreateAddress(&address); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": address})
}

// UpdateAddress updates an address
func (h *AddressHandler) UpdateAddress(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "USER_MOCK"
	}

	var req map[string]interface{}
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.addressService.UpdateAddress(id, userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address updated"})
}

// DeleteAddress deletes an address
func (h *AddressHandler) DeleteAddress(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "USER_MOCK"
	}

	err := h.addressService.DeleteAddress(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Address deleted"})
}

// SetDefaultAddress sets an address as default
func (h *AddressHandler) SetDefaultAddress(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("userID")
	if userID == "" {
		userID = "USER_MOCK"
	}

	err := h.addressService.SetDefaultAddress(id, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set default address"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Default address set"})
}
