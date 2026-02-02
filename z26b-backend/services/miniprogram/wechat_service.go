package miniprogram

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strings"

	"z26b-backend/internal"

	"gorm.io/gorm"
)

type WechatService struct {
	db *gorm.DB
}

func NewWechatService(db *gorm.DB) WechatServiceInterface {
	return &WechatService{db: db}
}

// WechatConfig 微信配置
type WechatConfig struct {
	AppID, AppSecret, MchID, MchKey, NotifyURL string
	Enabled                                    bool
}

func (s *WechatService) getWechatConfig() WechatConfig {
	return WechatConfig{
		AppID:     os.Getenv("WECHAT_APP_ID"),
		AppSecret: os.Getenv("WECHAT_APP_SECRET"),
		MchID:     os.Getenv("WECHAT_MCH_ID"),
		MchKey:    os.Getenv("WECHAT_MCH_KEY"),
		NotifyURL: os.Getenv("WECHAT_NOTIFY_URL"),
		Enabled:   os.Getenv("WECHAT_ENABLED") == "true",
	}
}

// WxLogin 微信登录
func (s *WechatService) WxLogin(code string) (map[string]interface{}, error) {
	config := s.getWechatConfig()

	// 调用微信API获取openid
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code",
		config.AppID, config.AppSecret, code)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if errcode, ok := result["errcode"]; ok && errcode.(float64) != 0 {
		return nil, fmt.Errorf("wechat api error: %v", result["errmsg"])
	}

	openID, ok := result["openid"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid openid")
	}

	// 获取或创建用户
	user, err := s.dbGetOrCreateUser(openID)
	if err != nil {
		return nil, err
	}

	// 生成token (简化版)
	token := s.generateToken(user.ID, openID)

	return map[string]interface{}{
		"token": token,
		"user":  user,
	}, nil
}

// dbGetOrCreateUser 数据库操作：获取或创建用户
func (s *WechatService) dbGetOrCreateUser(openID string) (*internal.User, error) {
	var user internal.User
	err := s.db.Where("open_id = ?", openID).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 用户不存在，创建新用户
			user = internal.User{
				ID:     internal.GenerateUUID(),
				OpenID: openID,
			}
			if err := s.db.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return &user, nil
}

// GenerateSign 生成签名
func (s *WechatService) GenerateSign(params map[string]interface{}, key string) string {
	// 排序参数
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 拼接字符串
	var str strings.Builder
	for _, k := range keys {
		if params[k] != "" {
			str.WriteString(fmt.Sprintf("%s=%v&", k, params[k]))
		}
	}
	str.WriteString("key=" + key)

	// MD5加密
	hash := md5.Sum([]byte(str.String()))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
}

// generateToken 生成token (简化版)
func (s *WechatService) generateToken(userID, openID string) string {
	hasher := md5.New()
	hasher.Write([]byte(userID + openID + "simple_salt"))
	return hex.EncodeToString(hasher.Sum(nil))
}
