# 后端API设计文档

## 一、技术栈与环境

```
- 语言: Go 1.21+
- Web框架: Gin 1.9+
- 数据库: MySQL 8.0+ + GORM 1.25+
- 缓存: Redis 7.0+ + go-redis v9
- 配置管理: Viper
- 日志: Zap + Lumberjack (日志轮转)
- 依赖注入: Wire
- 文档生成: Swag (Swagger)
- 单元测试: Testify + Mockery
- API验证: Gin Validator + Custom Validators
- 加密: Crypto/AES (API Key加密)
- JWT: jwt-go v5
```

---

## 二、项目结构

```
backend/
├── cmd/
│   └── server/
│       ├── main.go              # 应用入口
│       └── wire.go              # 依赖注入
│
├── configs/                     # 配置文件
│   ├── config.yaml             # 主要配置
│   └── config.local.yaml       # 本地开发配置
│
├── docs/                       # 文档
│   ├── api.md                  # API说明
│   └── swagger/                # Swagger文档
│
├── internal/
│   ├── api/                    # API层
│   │   ├── handler/            # HTTP处理器
│   │   │   ├── auth.go         # 认证相关
│   │   │   ├── user.go         # 用户管理
│   │   │   ├── ai.go           # AI配置
│   │   │   ├── training.go     # 训练计划
│   │   │   ├── nutrition.go    # 饮食计划
│   │   │   └── record.go       # 记录相关
│   │   │
│   │   ├── middleware/         # 中间件
│   │   │   ├── auth.go         # JWT认证
│   │   │   ├── cors.go         # CORS处理
│   │   │   ├── logger.go       # 请求日志
│   │   │   ├── recovery.go     # 错误恢复
│   │   │   ├── rate_limit.go   # 限流
│   │   │   └── validator.go    # 数据验证
│   │   │
│   │   ├── request/            # 请求DTO
│   │   │   ├── auth_request.go
│   │   │   ├── user_request.go
│   │   │   ├── ai_request.go
│   │   │   └── plan_request.go
│   │   │
│   │   └── response/           # 响应DTO
│   │       ├── base_response.go
│   │       ├── auth_response.go
│   │       └── plan_response.go
│   │
│   ├── service/                # 服务层
│   │   ├── auth_service.go     # 认证服务
│   │   ├── user_service.go     # 用户服务
│   │   ├── ai_service.go       # AI配置服务
│   │   ├── ai_integration.go   # AI集成服务
│   │   ├── training_service.go # 训练计划服务
│   │   ├── nutrition_service.go# 饮食计划服务
│   │   └── record_service.go   # 记录服务
│   │
│   ├── repository/             # 数据访问层
│   │   ├── user_repo.go
│   │   ├── ai_repo.go
│   │   ├── training_repo.go
│   │   ├── nutrition_repo.go
│   │   └── record_repo.go
│   │
│   ├── model/                  # 数据模型
│   │   ├── user.go
│   │   ├── ai_api.go
│   │   ├── training.go
│   │   ├── nutrition.go
│   │   └── record.go
│   │
│   ├── config/                 # 配置管理
│   │   └── config.go
│   │
│   ├── pkg/                    # 内部工具包
│   │   ├── redis/              # Redis客户端封装
│   │   │   └── redis.go
│   │   ├── jwt/                # JWT工具
│   │   │   └── jwt.go
│   │   ├── crypto/             # 加密工具
│   │   │   └── encryptor.go
│   │   ├── logger/             # 日志工具
│   │   │   └── logger.go
│   │   ├── ai/                 # AI客户端
│   │   │   ├── client_manager.go
│   │   │   ├── openai_client.go
│   │   │   ├── wenxin_client.go
│   │   │   └── tongyi_client.go
│   │   └── validator/          # 验证器
│   │       └── validators.go
│   │
│   └── errors/                 # 错误定义
│       ├── errors.go
│       └── error_codes.go
│
├── pkg/                        # 公共包
│   ├── sms/                    # 短信服务
│   ├── email/                  # 邮件服务
│   └── scheduler/              # 定时任务
│       └── training_reminder.go # 训练提醒
│
├── tests/                      # 测试
│   ├── unit/                   # 单元测试
│   └── integration/            # 集成测试
│
├── scripts/                    # 脚本
│   └── generate_swagger.sh     # 生成Swagger文档
│
├── go.mod                      # Go模块
├── go.sum                      # 依赖校验
└── Makefile                    # 构建脚本
```

---

## 三、统一API响应格式

### 1. 基础响应结构

```go
// internal/api/response/base_response.go
package response

type BaseResponse struct {
    Code      int         `json:"code"`              // 业务状态码
    Message   string      `json:"message"`           // 响应消息
    Data      interface{} `json:"data,omitempty"`    // 响应数据
    Timestamp int64       `json:"timestamp"`         // 时间戳
}

func Success(data interface{}) *BaseResponse {
    return &BaseResponse{
        Code:      200,
        Message:   "success",
        Data:      data,
        Timestamp: time.Now().Unix(),
    }
}

func Error(code int, message string) *BaseResponse {
    return &BaseResponse{
        Code:      code,
        Message:   message,
        Timestamp: time.Now().Unix(),
    }
}
```

### 2. 业务状态码定义

```go
// internal/errors/error_codes.go
package errors

const (
    // 成功
    Success = 200

    // 客户端错误 (4000系列)
    ErrBadRequest          = 4000  // 请求错误
    ErrInvalidParam        = 4001  // 参数无效
    ErrUnauthorized        = 4010  // 未认证
    ErrForbidden           = 4030  // 无权限
    ErrNotFound            = 4040  // 资源不存在
    ErrMethodNotAllowed    = 4050  // 方法不允许
    ErrConflict            = 4090  // 冲突

    // 服务器错误 (5000系列)
    ErrInternalServer      = 5000  // 内部错误
    ErrExternalService     = 5001  // 外部服务错误
    ErrDatabase            = 5002  // 数据库错误
    ErrCache               = 5003  // 缓存错误

    // 业务错误 (6000系列)
    ErrUserExists          = 6001  // 用户已存在
    ErrUserNotFound        = 6002  // 用户不存在
    ErrWrongPassword       = 6003  // 密码错误
    ErrTokenExpired        = 6004  // Token过期
    ErrPlanNotFound        = 6005  // 计划不存在
    ErrAiApiNotConfigured  = 6006  // AI API未配置
    ErrApiLimitExceeded    = 6007  // API调用超限
)
```

---

## 四、API接口详细设计

### 1. 认证API

#### 1.1 用户注册
```
POST /api/v1/auth/register

Request:
{
  "username": "user123",
  "email": "user@example.com",
  "password": "Abc123!@#",
  "confirmPassword": "Abc123!@#"
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com",
      "created_at": "2024-01-01T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "timestamp": 1704067200
}

Validation Rules:
- username: 3-20字符，字母数字下划线
- email: 有效邮箱格式
- password: 8-20字符，至少1大写1小写1数字1特殊字符
- confirmPassword: 必须等于password
```

#### 1.2 用户登录
```
POST /api/v1/auth/login

Request:
{
  "username": "user123",    // 用户名或邮箱
  "password": "Abc123!@#"
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  },
  "timestamp": 1704067200
}

登录成功后在Header中返回:
Authorization: Bearer {access_token}
```

#### 1.3 Token刷新
```
POST /api/v1/auth/refresh

Headers:
Authorization: Bearer {refresh_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "access_token": "new_access_token...",
    "expires_in": 3600
  },
  "timestamp": 1704067200
}
```

#### 1.4 登出
```
POST /api/v1/auth/logout

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "timestamp": 1704067200
}
```

---

### 2. 用户管理API

#### 2.1 获取个人信息
```
GET /api/v1/user/profile

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "user123",
      "email": "user@example.com",
      "phone": "13800138000",
      "avatar": "https://example.com/avatar.jpg",
      "created_at": "2024-01-01T00:00:00Z"
    }
  },
  "timestamp": 1704067200
}
```

#### 2.2 更新个人信息
```
PUT /api/v1/user/profile

Headers:
Authorization: Bearer {access_token}

Request:
{
  "username": "new_username",
  "phone": "13900139000",
  "avatar": "https://example.com/new_avatar.jpg"
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "new_username",
      "phone": "13900139000",
      "avatar": "https://example.com/new_avatar.jpg"
    }
  },
  "timestamp": 1704067200
}
```

---

### 3. 身体数据API

#### 3.1 添加身体数据
```
POST /api/v1/body-data

Headers:
Authorization: Bearer {access_token}

Request:
{
  "age": 25,
  "gender": "male",
  "height": 175.50,
  "weight": 70.00,
  "body_fat_percentage": 15.50,
  "muscle_percentage": 38.20,
  "measurement_date": "2024-01-01"
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "user_id": 1,
    "age": 25,
    "gender": "male",
    "height": 175.50,
    "weight": 70.00,
    "measurement_date": "2024-01-01"
  },
  "timestamp": 1704067200
}
```

---

### 4. AI配置API

#### 4.1 添加AI API配置
```
POST /api/v1/ai-apis

Headers:
Authorization: Bearer {access_token}

Request:
{
  "provider": "openai",        // openai/wenxin/tongyi
  "name": "我的OpenAI",
  "api_endpoint": "https://api.openai.com/v1",
  "api_key": "sk-****************************************",
  "model": "gpt-4-turbo",
  "max_tokens": 2000,
  "temperature": 0.7,
  "is_default": false
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "api": {
      "id": 1,
      "provider": "openai",
      "name": "我的OpenAI",
      "api_endpoint": "https://api.openai.com/v1",
      "model": "gpt-4-turbo",
      "is_default": true,
      "status": 1,
      "created_at": "2024-01-01T00:00:00Z"
    }
  },
  "timestamp": 1704067200
}

注意: api_key会在服务端加密存储
```

#### 4.2 获取AI API列表
```
GET /api/v1/ai-apis

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "apis": [
      {
        "id": 1,
        "provider": "openai",
        "name": "我的OpenAI",
        "api_endpoint": "https://api.openai.com/v1",
        "model": "gpt-4-turbo",
        "is_default": true,
        "status": 1
      }
    ]
  },
  "timestamp": 1704067200
}
```

#### 4.3 测试AI API连通性
```
POST /api/v1/ai-apis/{id}/test

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "test_result": {
      "status": "success",
      "response_time": 850,      // 响应时间(毫秒)
      "model_info": {
        "name": "gpt-4-turbo",
        "max_tokens": 128000
      }
    }
  },
  "timestamp": 1704067200
}
```

---

### 5. 运动能力评估API

#### 5.1 创建评估
```
POST /api/v1/assessments

Headers:
Authorization: Bearer {access_token}

Request:
{
  "experience_level": "intermediate",
  "weekly_available_days": 4,
  "daily_available_minutes": 60,
  "activity_type": "strength_training",
  "injury_history": "曾有膝盖扭伤史，已康复",
  "health_conditions": "无严重疾病",
  "preferred_days": ["monday", "wednesday", "friday", "saturday"],
  "equipment_available": ["dumbbells", "barbell", "bench"]
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "assessment": {
      "id": 1,
      "user_id": 1,
      "experience_level": "intermediate",
      "weekly_available_days": 4,
      "daily_available_minutes": 60,
      "assessment_date": "2024-01-01"
    }
  },
  "timestamp": 1704067200
}
```

---

### 6. 训练计划API

#### 6.1 生成训练计划
```
POST /api/v1/training-plans/generate

Headers:
Authorization: Bearer {access_token}

Request:
{
  "name": "12周增肌计划",
  "duration_weeks": 12,
  "goal": "muscle_gain",      // muscle_gain/fat_loss/endurance
  "difficulty": "medium",     // easy/medium/hard/extreme
  "ai_api_id": 1
}

Response (异步):
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "tsk_abc123def456",
    "status": "generating",
    "estimated_time": 30        // 预计生成时间(秒)
  },
  "timestamp": 1704067200
}

该请求返回后立即返回task_id，客户端可通过WebSocket或轮询查询进度
```

#### 6.2 查询计划生成任务状态
```
GET /api/v1/training-plans/tasks/{task_id}

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "task": {
      "task_id": "tsk_abc123def456",
      "status": "completed",    // generating/completed/failed
      "progress": 100,          // 生成进度(百分比)
      "result": {
        "plan_id": 1001,
        "plan_name": "12周增肌计划"
      },
      "error_message": null
    }
  },
  "timestamp": 1704067200
}
```

#### 6.3 获取训练计划列表
```
GET /api/v1/training-plans?page=1&limit=10

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "plans": [
      {
        "id": 1001,
        "name": "12周增肌计划",
        "start_date": "2024-01-01",
        "end_date": "2024-03-31",
        "total_weeks": 12,
        "difficulty_level": "medium",
        "status": "active"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 5,
      "total_pages": 1
    }
  },
  "timestamp": 1704067200
}
```

#### 6.4 获取训练计划详情
```
GET /api/v1/training-plans/{id}

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "plan": {
      "id": 1001,
      "name": "12周增肌计划",
      "start_date": "2024-01-01",
      "end_date": "2024-03-31",
      "total_weeks": 12,
      "difficulty_level": "medium",
      "plan_data": {
        "weeks": [
          {
            "week": 1,
            "days": [
              {
                "day": 1,
                "date": "2024-01-01",
                "type": "strength",
                "focus_area": "upper_body",
                "exercises": [
                  {
                    "name": "杠铃卧推",
                    "sets": 4,
                    "reps": "8-10",
                    "weight": "70kg",
                    "rest": "90s",
                    "difficulty": "medium",
                    "safety_notes": "注意肩胛骨收紧"
                  }
                ],
                "duration": 60,
                "estimated_calories": 350
              }
            ]
          }
        ]
      },
      "status": "active"
    }
  },
  "timestamp": 1704067200
}
```

#### 6.5 获取今日训练
```
GET /api/v1/training-plans/today

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "schedule": {
      "date": "2024-01-15",
      "type": "strength",
      "focus_area": "lower_body",
      "exercises": [...],
      "duration": 60,
      "is_completed": false,
      "completed_exercises": 0,
      "total_exercises": 5
    }
  },
  "timestamp": 1704067200
}
```

---

### 7. 饮食计划API

#### 7.1 生成饮食计划
```
POST /api/v1/nutrition-plans/generate

Headers:
Authorization: Bearer {access_token}

Request:
{
  "name": "减脂饮食计划",
  "duration_days": 30,
  "daily_calories": 2000,
  "macro_ratio": {
    "protein": 0.35,
    "carbs": 0.40,
    "fat": 0.25
  },
  "dietary_restrictions": ["lactose_intolerant"],
  "ai_api_id": 1
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "task_id": "tsk_nut_abc123",
    "status": "generating",
    "estimated_time": 20
  },
  "timestamp": 1704067200
}
```

#### 7.2 获取今日饮食
```
GET /api/v1/nutrition-plans/today

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "plan": {
      "target_calories": 2000,
      "target_protein": 175,
      "target_carbs": 200,
      "target_fat": 56
    },
    "meals": {
      "breakfast": {
        "time": "07:00-08:00",
        "foods": [
          {
            "name": "燕麦片",
            "amount": "50g",
            "calories": 190,
            "protein": 5
          }
        ],
        "total_calories": 450
      },
      "lunch": {...},
      "dinner": {...},
      "snacks": {...}
    }
  },
  "timestamp": 1704067200
}
```

---

### 8. 训练记录API

#### 8.1 记录训练
```
POST /api/v1/training-records

Headers:
Authorization: Bearer {access_token}

Request:
{
  "plan_id": 1001,
  "workout_date": "2024-01-15",
  "workout_type": "strength",
  "duration_minutes": 65,
  "exercises": [
    {
      "exercise_name": "深蹲",
      "sets": 4,
      "reps_per_set": [12, 12, 10, 10],
      "weight_used": [70, 70, 70, 70],
      "notes": "最后一组较困难"
    }
  ],
  "performance_data": {
    "total_volume": 14400,
    "estimated_calories": 350
  },
  "notes": "状态不错",
  "rating": 4,
  "injury_report": null
}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "record": {
      "id": 5001,
      "user_id": 1,
      "workout_date": "2024-01-15",
      "duration_minutes": 65,
      "exercises": [...]
    }
  },
  "timestamp": 1704067200
}
```

---

### 9. 数据统计API

#### 9.1 获取训练统计
```
GET /api/v1/stats/training?start_date=2024-01-01&end_date=2024-01-31

Headers:
Authorization: Bearer {access_token}

Response:
{
  "code": 200,
  "message": "success",
  "data": {
    "summary": {
      "total_workouts": 16,
      "total_duration": 1040,      // 分钟
      "total_calories": 5200,
      "avg_rating": 4.2
    },
    "progress": {
      "weekly_workouts": [3, 4, 4, 5],  // 每周训练次数
      "weekly_volume": [12000, 14500, 15000, 16200]  // 每周训练量
    }
  },
  "timestamp": 1704067200
}
```

---

## 五、WebSocket接口 (可选)

### 1. 计划生成进度推送
```
连接: ws://api.example.com/ws/tasks?token={access_token}

服务端推送:
{
  "type": "plan_generation_progress",
  "data": {
    "task_id": "tsk_abc123def456",
    "progress": 50,
    "message": "正在生成训练计划..."
  }
}
```

---

## 六、Middleware设计

### 1. JWT认证中间件
```go
// internal/api/middleware/auth.go
func JWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            response.UnauthorizedError(c, "未提供认证信息")
            c.Abort()
            return
        }

        // 解析Token
        claims, err := jwt.ParseToken(tokenString)
        if err != nil {
            response.UnauthorizedError(c, "认证信息无效")
            c.Abort()
            return
        }

        // 检查Redis中的会话
        sessionKey := fmt.Sprintf("session:%s", claims.SessionID)
        if !redis.Exists(sessionKey) {
            response.UnauthorizedError(c, "会话已过期")
            c.Abort()
            return
        }

        // 将用户信息存入context
        c.Set("user_id", claims.UserID)
        c.Set("username", claims.Username)
        c.Next()
    }
}
```

### 2. 限流中间件
```go
// internal/api/middleware/rate_limit.go
func RateLimit(apiID string, limit int64, duration time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        userID := c.GetInt64("user_id")
        key := fmt.Sprintf("api_limit:%s:%d:%s", apiID, userID, duration)

        current, err := redis.Incr(key)
        if err != nil {
            c.JSON(500, gin.H{"error": "限流服务异常"})
            c.Abort()
            return
        }

        if current == 1 {
            redis.Expire(key, duration)
        }

        if current > limit {
            c.JSON(429, gin.H{
                "error": "API调用频率过高",
                "retry_after": duration / time.Second
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

### 3. 数据验证中间件
```go
// internal/api/middleware/validator.go
func ValidateRequest(req interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := c.ShouldBindJSON(req); err != nil {
            c.JSON(400, gin.H{
                "error": "参数验证失败",
                "details": err.Error()
            })
            c.Abort()
            return
        }

        if err := validator.Validate(req); err != nil {
            c.JSON(400, gin.H{
                "error": "参数验证失败",
                "details": err.Error()
            })
            c.Abort()
            return
        }

        c.Set("validated_request", req)
        c.Next()
    }
}
```

---

## 七、AI集成服务设计

### 1. AI客户端管理器
```go
// internal/pkg/ai/client_manager.go
type ClientManager struct {
    clients map[string]AIClient
}

type AIClient interface {
    GenerateTrainingPlan(ctx context.Context, params TrainingPlanParams) (*TrainingPlan, error)
    GenerateNutritionPlan(ctx context.Context, params NutritionPlanParams) (*NutritionPlan, error)
    TestConnection(ctx context.Context) error
}

func (m *ClientManager) GetClient(provider string) (AIClient, error) {
    client, exists := m.clients[provider]
    if !exists {
        return nil, fmt.Errorf("unsupported provider: %s", provider)
    }
    return client, nil
}
```

### 2. OpenAI实现
```go
// internal/pkg/ai/openai_client.go
type OpenAIClient struct {
    apiKey     string
    baseURL    string
    model      string
    maxTokens  int
    temperature float32
}

func (c *OpenAIClient) GenerateTrainingPlan(ctx context.Context, params TrainingPlanParams) (*TrainingPlan, error) {
    prompt := buildTrainingPlanPrompt(params)

    resp, err := c.doRequest(ctx, &OpenAIRequest{
        Model:       c.model,
        Messages:    []Message{{Role: "user", Content: prompt}},
        MaxTokens:   c.maxTokens,
        Temperature: c.temperature,
    })

    if err != nil {
        return nil, err
    }

    // 解析AI返回的JSON
    var plan TrainingPlan
    if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &plan); err != nil {
        return nil, fmt.Errorf("failed to parse AI response: %w", err)
    }

    return &plan, nil
}
```

---

## 八、错误处理设计

### 1. 自定义错误类型
```go
// internal/errors/errors.go
type AppError struct {
    Code    int
    Message string
    Err     error
}

func (e *AppError) Error() string {
    return fmt.Sprintf("code=%d, message=%s, error=%v", e.Code, e.Message, e.Err)
}

func New(code int, message string) *AppError {
    return &AppError{Code: code, Message: message}
}

func Wrap(err error, code int, message string) *AppError {
    return &AppError{Code: code, Message: message, Err: err}
}

// 常见错误
var (
    ErrUserNotFound    = New(ErrUserNotFound, "用户不存在")
    ErrUnauthorized    = New(ErrUnauthorized, "未授权")
    ErrInvalidAPIKey   = New(ErrBadRequest, "无效的API Key")
)
```

### 2. 全局错误处理
```go
// internal/api/handler/base.go
func HandleError(c *gin.Context, err error) {
    appErr, ok := err.(*errors.AppError)
    if !ok {
        // 未知错误
        c.JSON(500, response.Error(ErrInternalServer, "服务器内部错误"))
        return
    }

    switch appErr.Code {
    case ErrBadRequest:
        c.JSON(400, response.Error(appErr.Code, appErr.Message))
    case ErrUnauthorized:
        c.JSON(401, response.Error(appErr.Code, appErr.Message))
    case ErrForbidden:
        c.JSON(403, response.Error(appErr.Code, appErr.Message))
    case ErrNotFound:
        c.JSON(404, response.Error(appErr.Code, appErr.Message))
    default:
        c.JSON(500, response.Error(appErr.Code, appErr.Message))
    }
}
```

---

## 九、配置文件

```yaml
# configs/config.yaml
# 应用配置
app:
  name: "AI Fitness Planner"
  version: "1.0.0"
  port: 8080
  mode: "release"  # debug/release/test
  secret_key: "your-secret-key-here"

# 数据库配置
database:
  mysql:
    host: "localhost"
    port: 3306
    user: "fitness_user"
    password: "your-password"
    dbname: "fitness_db"
    max_open_conns: 25
    max_idle_conns: 5
    conn_max_lifetime: 300s

  redis:
    host: "localhost"
    port: 6379
    password: "your-redis-password"
    db: 0
    pool_size: 10
    max_retries: 3

# JWT配置
jwt:
  secret: "jwt-secret-key"
  access_token_expire: 3600      # 1小时
  refresh_token_expire: 604800    # 7天

# AI配置
ai:
  max_concurrent_requests: 10
  timeout: 60s
  retry_attempts: 3
  retry_delay: 5s

# 限流配置
rate_limit:
  api_calls_per_minute: 60
  api_calls_per_hour: 1000
  api_calls_per_day: 10000

# 日志配置
log:
  level: "info"  # debug/info/warn/error
  filename: "logs/app.log"
  max_size: 500  # MB
  max_backups: 10
  max_age: 30    # days
```

---

## 十、Makefile

```makefile
.PHONY: run build test clean swagger wire

# 变量
GO := go
BINARY_NAME := fitness-server

# 运行应用
run:
	go run cmd/server/main.go

# 构建应用
build:
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

# 运行测试
test:
	go test -v ./...

# 清理
 clean:
	rm -rf bin/
	go clean

# 生成Swagger文档
swagger:
	swag init -g cmd/server/main.go -o docs/swagger

# 生成依赖注入代码
wire:
	wire ./cmd/server

# 安装依赖
deps:
	go mod download
	go mod tidy

# 运行Docker
docker:
	docker build -t fitness-server .
	docker run -p 8080:8080 fitness-server
```
