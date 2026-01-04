# 安全性和权限机制设计文档

## 一、系统安全架构概述

```
┌─────────────────────────────────────────────────────────┐
│                    应用层安全                            │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   身份认证   │  │    授权     │  │   加密存储    │  │
│  │  (JWT+Session)│  │  (RBAC)    │  │  (AES-256)   │  │
│  └─────────────┘  └──────────────┘  └──────────────┘  │
├─────────────────────────────────────────────────────────┤
│                    传输层安全                            │
│  ┌───────────────────────────────────────────────────┐ │
│  │     TLS 1.3 + 证书管理 + 强制HTTPS重定向          │ │
│  └───────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────┤
│                    存储层安全                            │
│  ┌────────────┐  ┌────────────┐  ┌─────────────────┐   │
│  │数据库访问控制│  │ 数据加密   │  │ 备份与恢复机制   │ │
│  │            │  │(动态/静态) │  │                 │  │
│  └────────────┘  └────────────┘  └─────────────────┘  │
├─────────────────────────────────────────────────────────┤
│                    防御机制                              │
│  ┌──────────┐  ┌──────────┐  ┌─────────┐  ┌─────────┐  │
│  │ API限流  │  │   WAF    │  │ SQL注入 │  │   XSS   │  │
│  │ (Redis)  │  │ (规则)   │  │ 防护    │  │  防护   │  │
│  └──────────┘  └──────────┘  └─────────┘  └─────────┘  │
└─────────────────────────────────────────────────────────┘
```

---

## 二、身份认证与授权

### 1. JWT Token 体系设计

#### 1.1 Token结构
```go
// internal/pkg/jwt/claims.go
type Claims struct {
    UserID    int64  `json:"user_id"`
    Username  string `json:"username"`
    SessionID string `json:"session_id"`  // 用于会话管理
    Type      string `json:"type"`        // access/refresh
    jwt.RegisteredClaims
}
```

#### 1.2 双Token机制

**Access Token（访问令牌）**
- 有效期：1小时
- 存储位置：HTTP Header (Authorization: Bearer xxx)
- 使用场景：常规API请求

**Refresh Token（刷新令牌）**
- 有效期：7天
- 存储位置：HTTP Only Cookie
- 使用场景：Token续期

```go
// internal/api/handler/auth.go
func (h *AuthHandler) Login(c *gin.Context) {
    // ...验证用户凭证...

    // 生成Token对
    accessClaims := &jwt.Claims{
        UserID:    user.ID,
        Username:  user.Username,
        SessionID: uuid.New().String(),
        Type:      "access",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
        },
    }

    refreshClaims := &jwt.Claims{
        UserID:    user.ID,
        Username:  user.Username,
        SessionID: accessClaims.SessionID,
        Type:      "refresh",
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
        },
    }

    accessToken, _ := jwt.GenerateToken(accessClaims, h.config.JWT.Secret)
    refreshToken, _ := jwt.GenerateToken(refreshClaims, h.config.JWT.Secret)

    // 存储Session到Redis
    sessionKey := fmt.Sprintf("session:%s", accessClaims.SessionID)
    redis.HSet(c.Request.Context(), sessionKey, map[string]interface{}{
        "user_id":    user.ID,
        "created_at": time.Now().Unix(),
    })
    redis.Expire(c.Request.Context(), sessionKey, 7*24*time.Hour)

    // 设置HTTP Only Cookie
    c.SetCookie(
        "refresh_token",
        refreshToken,
        7*24*3600,
        "/",
        "",
        true,    // Secure
        true,    // HttpOnly
    )

    c.JSON(200, response.Success(gin.H{
        "access_token": accessToken,
        "expires_in":   3600,
        "user":         user,
    }))
}
```

#### 1.3 Token刷新流程

```
用户          前端          后端           Redis
 │             │             │              │
 │  Token过期   │             │              │
 │────────────>│             │              │
 │             │ 检测401     │              │
 │             │────────────>│              │
 │             │             │ 验证Refresh  │
 │             │             │─────────────>│
 │             │             │              │
 │             │             │ 生成新Token  │
 │             │ 返回新Token │<─────────────│
 │             │<────────────│              │
 │ 继续请求    │             │              │
 │────────────>│             │              │
 │             │ 新Token     │              │
 │             │────────────>│              │
```

---

### 2. 权限控制模型 (RBAC)

#### 2.1 角色定义
```go
const (
    RoleUser      = "user"      // 普通用户
    RoleVIP       = "vip"       // VIP用户
    RoleAdmin     = "admin"     // 管理员
    RoleSuper     = "super"     // 超级管理员
)
```

#### 2.2 权限矩阵

| 功能模块 | 用户(user) | VIP(vip) | 管理(admin) | 超管(super) |
|---------|----------|----------|------------|------------|
| **个人档案** | | | | |
| 查看/修改自己信息 | ✓ | ✓ | ✓ | ✓ |
| 查看他人信息 | ✗ | ✗ | ✓ | ✓ |
| **AI配置** | | | | |
| 管理自己的API | ✓ | ✓ | ✓ | ✓ |
| 查看所有API配置 | ✗ | ✗ | ✓ | ✓ |
| **训练计划** | | | | |
| 生成个人计划 | ✓ | ✓ | ✓ | ✓ |
| 高级计划模板 | ✗ | ✓ | ✓ | ✓ |
| **数据统计** | | | | |
| 查看个人统计 | ✓ | ✓ | ✓ | ✓ |
| 查看全局统计 | ✗ | ✗ | ✓ | ✓ |

---

### 3. 跨域资源共享 (CORS)

```go
// internal/api/middleware/cors.go
func CORS() gin.HandlerFunc {
    corsConfig := cors.Config{
        AllowOrigins:     []string{"https://fitness.example.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
        AllowOriginFunc: func(origin string) bool {
            // 动态验证
            return strings.HasSuffix(origin, ".example.com")
        },
    }

    return cors.New(corsConfig)
}
```

#### 3.1 CORS请求流程
```
客户端          服务端
  │                │
  │─OPTIONS───────>│ 预检请求
  │                │
  │<─200───────────│ 返回允许的Methods/Headers
  │                │
  │─实际请求───────>│
  │                │
  │<─响应──────────│
```

---

## 三、API密钥安全管理

### 1. AES-256加密存储

```go
// internal/pkg/crypto/encryptor.go
type AESEncryptor struct {
    key []byte
}

func NewAESEncryptor(key string) *AESEncryptor {
    // 密钥派生 (PBKDF2)
    derivedKey := pbkdf2.Key(
        []byte(key),
        []byte("fitness-salt"),
        10000,
        32,
        sha256.New,
    )

    return &AESEncryptor{key: derivedKey}
}

func (e *AESEncryptor) Encrypt(plaintext string) (string, error) {
    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", err
    }

    // 生成随机IV
    iv := make([]byte, aes.BlockSize)
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return "", err
    }

    // GCM模式加密
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    ciphertext := gcm.Seal(nil, iv, []byte(plaintext), nil)

    // 存储格式: iv:nonce:ciphertext
    encrypted := append(iv, ciphertext...)
    return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (e *AESEncryptor) Decrypt(ciphertext string) (string, error) {
    data, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    if len(data) < aes.BlockSize {
        return "", errors.New("ciphertext too short")
    }

    iv := data[:aes.BlockSize]
    ciphertextData := data[aes.BlockSize:]

    block, err := aes.NewCipher(e.key)
    if err != nil {
        return "", err
    }

    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    plaintext, err := gcm.Open(nil, iv, ciphertextData, nil)
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}
```

### 2. API密钥管理流程

```
用户添加API Key:

用户输入             前端            后端              数据库
  │                  │              │                  │
  │─明文API Key────>│              │                  │
  │                  │─HTTPS───────>│                  │
  │                  │              │ 1. 验证格式      │
  │                  │              │ 2. 测试连通性    │
  │                  │              │ 3. AES加密       │
  │                  │              │─加密后数据───────>│
  │                  │              │                  │
  │                  │              │                  │ 存储:
  │                  │              │                  │ provider, name,
  │                  │              │                  │ endpoint, model,
  │                  │              │                  │ encrypted_key
```

### 3. 使用环境变量管理密钥

```bash
# .env.example
APP_SECRET=your-app-secret-key-32-char-long
JWT_SECRET=your-jwt-secret-key-at-least-32-char
ENCRYPTION_KEY=your-aes-256-encryption-key

# Redis
REDIS_PASSWORD=your-redis-password

# Database
DB_PASSWORD=your-db-password

# AI API Keys (可选)
OPENAI_API_KEY=sk-********
WENXIN_API_KEY=********
```

---

## 四、数据安全与隐私保护

### 1. 数据分类与保护级别

#### Level 1 - 公开数据
- 类型：普通用户信息、训练计划模板
- 保护：基础访问控制

#### Level 2 - 敏感数据
- 类型：身体数据、健身目标、训练记录
- 保护：用户隔离加密存储，AES-256

#### Level 3 - 机密数据
- 类型：API密钥、支付信息、个人身份信息
- 保护：专用加密、访问审计、定期轮换

### 2. 数据隔离策略

```go
// 数据访问控制 base.go
type DataAccessControl struct {
    UserID   int64
    Role     string
}

// 所有数据访问必须验证所属用户
func (r *Repository) GetTrainingPlan(ctx context.Context, planID int64, userID int64) (*model.TrainingPlan, error) {
    query := `
        SELECT * FROM training_plans
        WHERE id = ? AND user_id = ?
    `

    if err := r.db.Raw(query, planID, userID).Scan(&plan).Error; err != nil {
        return nil, err
    }

    return plan, nil
}
```

### 3. 数据备份与恢复

```go
// 自动备份策略
type BackupConfig struct {
    Enabled          bool   `yaml:"enabled"`
    BackupDir        string `yaml:"backup_dir"`
    BackupSchedule   string `yaml:"backup_schedule"`   // 备份计划 (cron)
    RetentionDays    int    `yaml:"retention_days"`    // 保留天数
    EncryptionKey    string `yaml:"encryption_key"`    // 备份加密密钥
}

// 备份策略
// - 完全备份: 每周日 02:00
// - 增量备份: 每天 02:00
// - 备份验证: 备份后自动验证
// - 异地备份: 同步到其他数据中心
```

---

## 五、攻击防护机制

### 1. API限流与防刷

#### 1.1 基于Redis的令牌桶算法

```go
// internal/api/middleware/rate_limit.go
func TokenBucketRateLimit(limit int, window time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt64("user_id")
        key := fmt.Sprintf("rate_limit:%d:%s", userID, window)

        // 令牌桶算法实现
        current, err := redis.Get(c.Request.Context(), key).Int()
        if err == redis.Nil {
            // 桶不存在，创建并填充令牌
            redis.Set(c.Request.Context(), key, limit-1, window)
            c.Next()
            return
        }

        if err != nil {
            c.JSON(500, gin.H{"error": "限流服务异常"})
            c.Abort()
            return
        }

        if current <= 0 {
            // 无可用令牌
            retryAfter := redis.TTL(c.Request.Context(), key).Val()
            c.JSON(429, gin.H{
                "error": "请求频率过高",
                "retry_after": retryAfter,
                "limit":       limit,
                "window":      window.String(),
            })
            c.Abort()
            return
        }

        // 消耗一个令牌
        redis.Decr(c.Request.Context(), key)
        c.Next()
    }
}
```

#### 1.2 等级限流策略

| 用户等级 | 每分钟 | 每小时 | 每日 |
|---------|-------|-------|------|
| User    | 60    | 1000  | 10000 |
| VIP     | 120   | 2000  | 20000 |
| Admin   | 200   | 5000  | 50000 |

### 2. SQL注入防护

```go
// 使用GORM预编译参数化查询
// 错误示例 (SQL注入风险)
db.Raw(fmt.Sprintf("SELECT * FROM users WHERE username='%s' AND password='%s'", username, password))

// 正确示例
var user model.User
db.Where("username = ? AND password = ?", username, encryptedPassword).First(&user)

// 对特殊字符进行参数化
db.Where("email LIKE ?", "%"+searchEmail+"%").Find(&users)
```

#### 2.1 输入验证过滤

```go
// SQL注入关键词检测
var sqlInjectionPattern = regexp.MustCompile(`(?i)(union|select|insert|update|delete|drop|create|--)`)

func SanitizeInput(input string) (string, error) {
    if sqlInjectionPattern.MatchString(input) {
        return "", errors.ErrInvalidInput
    }

    // HTML转义
    input = html.EscapeString(input)

    return input, nil
}
```

### 3. XSS跨站脚本防护

#### 3.1 输出编码

```go
// 对所有用户输入进行HTML转义
func EncodeHTML(s string) string {
    return html.EscapeString(s)
}

// 特殊场景: JSON响应时自动转义
func JSONSafeResponse(data interface{}) ([]byte, error) {
    return json.Marshal(data, json.EscapeHTML)
}
```

#### 3.2 Content Security Policy

```http
Content-Security-Policy: default-src 'self';
                         script-src 'self' 'unsafe-inline' https://trusted-cdn.example.com;
                         style-src 'self' 'unsafe-inline';
                         img-src 'self' data: https:;
                         font-src 'self' data:;
                         object-src 'none';
                         base-uri 'self';
```

### 4. CSRF防护

#### 4.1 双重Cookie提交

```go
// 生成CSRF Token
func GenerateCSRFToken() string {
    return uuid.New().String()
}

// 中间件验证
func CSRFFilter() gin.HandlerFunc {
    return func(c *gin.Context) {
        // GET请求不验证
        if c.Request.Method == http.MethodGet {
            c.Next()
            return
        }

        csrfToken := c.GetHeader("X-CSRF-Token")
        cookieToken, err := c.Cookie("csrf_token")

        if err != nil || csrfToken != cookieToken {
            c.JSON(403, gin.H{"error": "CSRF验证失败"})
            c.Abort()
            return
        }

        c.Next()
    }
}
```

---

## 六、日志与审计

### 1. 日志分级策略

```go
type LogLevel string

const (
    LogDebug     LogLevel = "DEBUG"
    LogInfo      LogLevel = "INFO"
    LogWarning   LogLevel = "WARNING"
    LogError     LogLevel = "ERROR"
    LogCritical  LogLevel = "CRITICAL"
)

// 安全相关操作日志
type SecurityLog struct {
    Timestamp   time.Time
    UserID      int64
    IP          string
    Action      string
    Resource    string
    Status      string
    Description string
}
```

#### 1.1 审计日志事件

| 事件类型 | 级别 | 记录内容 |
|---------|-----|---------|
| 用户登录/登出 | INFO | 时间、用户ID、IP、结果 |
| API密钥增删改 | WARNING | 用户ID、API ID、操作类型 |
| 训练计划生成 | INFO | 用户ID、API ID、耗时 |
| 权限变更 | WARNING | 管理员ID、目标用户、新权限 |
| 数据导出 | WARNING | 用户ID、导出数据量 |
| 敏感操作失败 | ERROR | 用户ID、失败原因 |