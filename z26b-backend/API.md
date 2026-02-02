# Z26B Backend API Documentation

## Base URL
```
http://localhost:8080/api
```

## Response Format

All responses are in JSON format:

### Success Response (2xx)
```json
{
  "data": {
    // Response data
  }
}
```

### Error Response (4xx, 5xx)
```json
{
  "error": "Error description"
}
```

---

## Goods API

### List Products
Get paginated list of products with optional filtering and search.

**Request:**
```
GET /goods/list?page=1&pageSize=10&categoryId=cat1&search=product
```

**Parameters:**
- `page` (int, optional): Page number, default 1
- `pageSize` (int, optional): Items per page, default 10, max 100
- `categoryId` (string, optional): Filter by category
- `search` (string, optional): Search by name or description

**Response:**
```json
{
  "data": {
    "records": [
      {
        "_id": "P1_prod",
        "name": "Product Name",
        "cover_image": "https://...",
        "price": 99.99,
        "status": "ENABLED"
      }
    ],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```

---

### Get Product Details
Get full details of a specific product including SKUs.

**Request:**
```
GET /goods/:id
```

**Parameters:**
- `id` (string, required): Product ID (SPU ID)

**Response:**
```json
{
  "data": {
    "_id": "P1_prod",
    "name": "Product Name",
    "detail": "<html>Product description</html>",
    "cover_image": "https://...",
    "swiper_images": ["https://...", "https://..."],
    "status": "ENABLED",
    "skus": [
      {
        "_id": "K1_prod",
        "price": 100,
        "count": 95,
        "description": "SKU variant"
      }
    ]
  }
}
```

---

### Get Product Comments
Get reviews and comments for a specific product.

**Request:**
```
GET /goods/:id/comments?page=1&pageSize=10&score=4
```

**Parameters:**
- `id` (string, required): Product ID
- `page` (int, optional): Page number
- `pageSize` (int, optional): Items per page
- `score` (int, optional): Filter by rating (1-5)

**Response:**
```json
{
  "data": {
    "pageNum": 1,
    "pageSize": 10,
    "totalCount": 50,
    "pageList": [
      {
        "_id": "comment_1",
        "userName": "John Doe",
        "commentScore": 4,
        "commentContent": "Great product!",
        "commentTime": 1234567890
      }
    ]
  }
}
```

---

## Shopping Cart API

### Get Cart Items
Get all items in user's shopping cart.

**Request:**
```
GET /cart/items
```

**Response:**
```json
{
  "data": {
    "isNotEmpty": true,
    "storeGoods": [
      {
        "storeId": "1000",
        "storeName": "Store Name",
        "goodsList": [
          {
            "_id": "cartitem_1",
            "quantity": 2,
            "isSelected": true,
            "sku": {
              "price": 100,
              "description": "SKU description"
            }
          }
        ]
      }
    ]
  }
}
```

---

### Add to Cart
Add a product variant to shopping cart.

**Request:**
```
POST /cart/add
Content-Type: application/json

{
  "skuId": "K1_prod",
  "quantity": 2
}
```

**Parameters:**
- `skuId` (string, required): SKU ID
- `quantity` (int, required): Quantity to add, min 1

**Response:**
```json
{
  "data": {
    "_id": "cartitem_1",
    "userId": "USER_MOCK",
    "skuId": "K1_prod",
    "quantity": 2,
    "isSelected": true
  }
}
```

---

### Update Cart Item
Update quantity or selection status of a cart item.

**Request:**
```
PUT /cart/update/:id
Content-Type: application/json

{
  "quantity": 3,
  "isSelected": true
}
```

**Parameters:**
- `id` (string, required): Cart item ID
- `quantity` (int, optional): New quantity
- `isSelected` (bool, optional): Selection status

**Response:**
```json
{
  "data": {
    "_id": "cartitem_1",
    "quantity": 3,
    "isSelected": true
  }
}
```

---

### Remove from Cart
Remove an item from shopping cart.

**Request:**
```
DELETE /cart/remove/:id
```

**Parameters:**
- `id` (string, required): Cart item ID

**Response:**
```json
{
  "message": "Item removed"
}
```

---

### Clear Cart
Remove all items from shopping cart.

**Request:**
```
POST /cart/clear
```

**Response:**
```json
{
  "message": "Cart cleared"
}
```

---

## Order API

### List Orders
Get user's orders with optional status filter.

**Request:**
```
GET /order/list?page=1&pageSize=10&status=TO_PAY
```

**Parameters:**
- `page` (int, optional): Page number
- `pageSize` (int, optional): Items per page
- `status` (string, optional): Filter by status

**Status Values:**
- `TO_PAY`: Awaiting payment
- `TO_SEND`: Awaiting shipment
- `TO_RECEIVE`: Awaiting receipt
- `FINISHED`: Completed
- `CANCELED`: Canceled
- `RETURN_APPLIED`: Return requested
- `RETURN_FINISHED`: Return completed

**Response:**
```json
{
  "data": {
    "records": [
      {
        "_id": "order_1",
        "userId": "USER_MOCK",
        "status": "TO_PAY",
        "totalPrice": 200,
        "finalPrice": 200,
        "createdAt": 1234567890,
        "items": []
      }
    ],
    "total": 15,
    "page": 1,
    "pageSize": 10
  }
}
```

---

### Get Order Details
Get full details of a specific order.

**Request:**
```
GET /order/:id
```

**Parameters:**
- `id` (string, required): Order ID

**Response:**
```json
{
  "data": {
    "_id": "order_1",
    "userId": "USER_MOCK",
    "status": "TO_PAY",
    "totalPrice": 200,
    "finalPrice": 200,
    "remarks": "Please handle with care",
    "items": [
      {
        "_id": "orderitem_1",
        "skuId": "K1_prod",
        "quantity": 2,
        "price": 100
      }
    ],
    "delivery_info": {}
  }
}
```

---

### Create Order
Create a new order from selected cart items.

**Request:**
```
POST /order/create
Content-Type: application/json

{
  "addressId": "addr_1",
  "remarks": "Please handle with care"
}
```

**Parameters:**
- `addressId` (string, required): Delivery address ID
- `remarks` (string, optional): Order remarks

**Response:**
```json
{
  "data": {
    "_id": "order_new_1",
    "userId": "USER_MOCK",
    "status": "TO_PAY",
    "totalPrice": 300,
    "finalPrice": 300,
    "createdAt": 1234567890
  }
}
```

---

### Cancel Order
Cancel a pending order.

**Request:**
```
PUT /order/cancel/:id
```

**Parameters:**
- `id` (string, required): Order ID

**Valid Status:** Only orders with status `TO_PAY` or `TO_SEND` can be canceled

**Response:**
```json
{
  "data": {
    "_id": "order_1",
    "status": "CANCELED"
  }
}
```

---

### Confirm Receipt
Confirm order receipt and change status to completed.

**Request:**
```
POST /order/confirm/:id
```

**Parameters:**
- `id` (string, required): Order ID

**Valid Status:** Only orders with status `TO_RECEIVE` can be confirmed

**Response:**
```json
{
  "data": {
    "_id": "order_1",
    "status": "FINISHED"
  }
}
```

---

## Address API

### List Addresses
Get all addresses for the user.

**Request:**
```
GET /address/list
```

**Response:**
```json
{
  "data": [
    {
      "_id": "addr_1",
      "name": "John Doe",
      "phone": "17612345678",
      "provinceName": "甘肃省",
      "cityName": "甘南藏族自治州",
      "detailAddress": "松日鼎盛大厦1层1号",
      "isDefault": 1
    }
  ]
}
```

---

### Get Address Details
Get details of a specific address.

**Request:**
```
GET /address/:id
```

**Parameters:**
- `id` (string, required): Address ID

**Response:**
```json
{
  "data": {
    "_id": "addr_1",
    "name": "John Doe",
    "phone": "17612345678",
    "countryName": "中国",
    "provinceName": "甘肃省",
    "cityName": "甘南藏族自治州",
    "districtName": "碌曲县",
    "detailAddress": "松日鼎盛大厦1层1号",
    "latitude": "34.59103",
    "longitude": "102.48699",
    "isDefault": 1
  }
}
```

---

### Create Address
Create a new address.

**Request:**
```
POST /address/create
Content-Type: application/json

{
  "name": "John Doe",
  "phone": "17612345678",
  "provinceName": "甘肃省",
  "cityName": "甘南藏族自治州",
  "detailAddress": "松日鼎盛大厦1层1号",
  "addressTag": "Home",
  "latitude": "34.59103",
  "longitude": "102.48699",
  "isDefault": 0
}
```

**Required Parameters:**
- `name`: Name
- `phone`: Phone number
- `provinceName`: Province name
- `cityName`: City name
- `detailAddress`: Detailed address

**Optional Parameters:**
- `addressTag`: Tag (Home, Office, etc.)
- `latitude`: Latitude
- `longitude`: Longitude
- `isDefault`: Set as default (0 or 1)

**Response:**
```json
{
  "data": {
    "_id": "addr_new_1",
    "name": "John Doe",
    "phone": "17612345678"
  }
}
```

---

### Update Address
Update an existing address.

**Request:**
```
PUT /address/update/:id
Content-Type: application/json

{
  "phone": "17612345679",
  "detailAddress": "New address"
}
```

**Parameters:**
- `id` (string, required): Address ID
- Other fields as needed

**Response:**
```json
{
  "data": {
    "_id": "addr_1",
    "phone": "17612345679"
  }
}
```

---

### Delete Address
Delete an address.

**Request:**
```
DELETE /address/:id
```

**Parameters:**
- `id` (string, required): Address ID

**Response:**
```json
{
  "message": "Address deleted"
}
```

---

### Set Default Address
Set an address as the default delivery address.

**Request:**
```
POST /address/setDefault/:id
```

**Parameters:**
- `id` (string, required): Address ID

**Response:**
```json
{
  "message": "Default address set"
}
```

---

## Comments/Reviews API

### Get Product Comments
Get all comments for a product (duplicate of goods/:id/comments).

**Request:**
```
GET /comment/list/:spuId?page=1&pageSize=10&score=4
```

**Parameters:**
- `spuId` (string, required): Product ID
- `page` (int, optional): Page number
- `pageSize` (int, optional): Items per page
- `score` (int, optional): Filter by rating (1-5)

**Response:**
```json
{
  "data": {
    "pageNum": 1,
    "pageSize": 10,
    "totalCount": 30,
    "pageList": []
  }
}
```

---

### Submit Comment
Submit a review for a product.

**Request:**
```
POST /comment/submit
Content-Type: application/json

{
  "spuId": "P1_prod",
  "skuId": "K1_prod",
  "userName": "John Doe",
  "commentContent": "Excellent product!",
  "commentScore": 5,
  "isAnonymity": false
}
```

**Parameters:**
- `spuId` (string, required): Product ID
- `commentContent` (string, required): Comment text
- `commentScore` (int, required): Rating 1-5
- `skuId` (string, optional): SKU ID
- `userName` (string, optional): User name
- `isAnonymity` (bool, optional): Anonymous submission

**Response:**
```json
{
  "data": {
    "_id": "comment_new_1",
    "spuId": "P1_prod",
    "commentScore": 5,
    "commentContent": "Excellent product!"
  }
}
```

---

## Home API

### Get Homepage Carousel
Get homepage banner/carousel images.

**Request:**
```
GET /home/swiper
```

**Response:**
```json
{
  "data": {
    "_id": "swiper_1",
    "images": ["https://...", "https://..."],
    "title": "Promotion Title",
    "priority": 1
  }
}
```

---

### Get Categories
Get product categories for homepage.

**Request:**
```
GET /home/categories
```

**Response:**
```json
{
  "data": [
    {
      "_id": "cat_1",
      "name": "Electronics",
      "icon": "https://...",
      "sort": 1
    },
    {
      "_id": "cat_2",
      "name": "Clothing",
      "icon": "https://...",
      "sort": 2
    }
  ]
}
```

---

### Get Promotions
Get active promotions.

**Request:**
```
GET /home/promotions
```

**Response:**
```json
{
  "data": [
    {
      "_id": "promo_1",
      "title": "Summer Sale",
      "promotionCode": "SUMMER2024",
      "description": "Save up to 50%",
      "promotionStatus": 3
    }
  ]
}
```

---

## SKU API

### Get SKU Details
Get details of a specific SKU.

**Request:**
```
GET /sku/:id
```

**Parameters:**
- `id` (string, required): SKU ID

**Response:**
```json
{
  "data": {
    "_id": "K1_prod",
    "spuId": "P1_prod",
    "price": 100,
    "count": 95,
    "description": "SKU variant description",
    "image": "https://..."
  }
}
```

---

### Get SKUs by Product
Get all SKU variants for a specific product.

**Request:**
```
GET /sku/list/:spuId
```

**Parameters:**
- `spuId` (string, required): Product ID

**Response:**
```json
{
  "data": [
    {
      "_id": "K1_prod",
      "price": 100,
      "count": 95,
      "description": "Size M"
    },
    {
      "_id": "K2_prod",
      "price": 100,
      "count": 50,
      "description": "Size L"
    }
  ]
}
```

---

## Error Codes

| Code | Status | Message |
|------|--------|---------|
| 200 | OK | Request successful |
| 400 | Bad Request | Invalid parameters |
| 404 | Not Found | Resource not found |
| 500 | Server Error | Internal server error |

---

## Rate Limiting

Currently no rate limiting. Implement based on your requirements.

---

## Authentication

Currently using mock user ID. Implement JWT authentication for production.

**Default User ID:** `USER_MOCK`

---

## Pagination

List endpoints support standard pagination:
- `page`: Page number (default 1)
- `pageSize`: Items per page (default 10, max 100)

Response format:
```json
{
  "data": {
    "records": [],
    "total": 100,
    "page": 1,
    "pageSize": 10
  }
}
```
