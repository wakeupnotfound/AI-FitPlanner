# Redis数据结构设计方案

## 1. 用户Session管理

### Key格式
```
session:{session_id} -> Hash
user:{user_id}:session -> Set (存储用户的所有session)
```

### 数据结构示例
```
session:abc123def456 -> {
    user_id: 1001,
    username: "user1",
    login_time: 1672502400,
    last_activity: 1672506000,
    ip_address: "192.168.1.1",
    user_agent: "Mozilla/5.0..."
}

user:1001:session -> ["abc123def456", "xyz789stu012"] (所有活跃session)
```

### TTL设置
- Session: 7天
- 用户session集合: 永久 (通过定时任务清理过期session)

---

## 2. AI API限流控制

### Key格式
```
api_limit:{api_id}:minute -> String (API调用计数-分钟)
api_limit:{api_id}:day -> String (API调用计数-天)
api_limit:{api_id}:month -> String (API调用计数-月)
```

### Redis Lua脚本 (原子操作)
```lua
-- 检查并增加计数
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local current = tonumber(redis.call('GET', key) or 0)

if current >= limit then
    return 0  -- 超过限制
end

redis.call('INCR', key)
if current == 0 then
    redis.call('EXPIRE', key, ARGV[2])  -- 设置过期时间
end

return 1  -- 调用成功
```

---

## 3. 训练计划缓存

### Key格式
```
training_plan:{plan_id} -> Hash (完整计划数据)
training_plan:{user_id}:active -> String (当前激活的计划ID)
training_plan:{plan_id}:schedule:{date} -> Hash (某天的训练安排)
```

### 数据结构示例
```
training_plan:1001 -> {
    id: 1001,
    name: "12周增肌计划",
    start_date: "2024-01-01",
    end_date: "2024-03-31",
    difficulty: "medium",
    total_weeks: 12,
    plan_data: "...JSON数据..."
}

training_plan:1001:schedule:2024-01-15 -> {
    day_type: "strength",
    focus_area: "lower_body",
    exercises: [
        {
            name: "深蹲",
            sets: 4,
            reps: 12,
            weight: "70kg",
            rest: "90s",
            difficulty: "medium"
        }
    ],
    duration: 60,
    estimated_calories: 350
}
```

### TTL设置
- 计划详情：永久（通过版本更新机制清理）
- 每日计划：7天
- 用户当前计划：永久

---

## 4. 饮食计划缓存

### Key格式
```
nutrition_plan:{plan_id} -> Hash (完整饮食计划)
nutrition_plan:{user_id}:{date}:meals -> Hash (某天的饮食安排)
```

### 数据结构示例
```
nutrition_plan:2001 -> {
    id: 2001,
    name: "减脂饮食计划",
    daily_calories: 2000,
    macro_ratio: {
        protein: 0.35,
        carbs: 0.40,
        fat: 0.25
    }
}

daily_nutrition:1001:2024-01-15 -> {
    target_calories: 2000,
    target_protein: 175,
    target_cabs: 200,
    target_fat: 56,
    meals: {
        breakfast: [...],
        lunch: [...],
        dinner: [...],
        snacks: [...]
    }
}
```

### TTL设置
- 饮食计划：30天
- 每日营养安排：7天

---

## 5. AI提示词模板缓存

### Key格式
```
template:{category}:{template_id} -> Hash (模板详情)
template:categories -> Set (所有模板分类)
template:{category}:default -> String (默认模板ID)
```

### 数据结构示例
```
template:training:1001 -> {
    id: 1001,
    name: "基础训练计划模板",
    category: "training",
    subcategory: "strength",
    template: "根据以下信息为用户生成训练计划：...",
    variables: ["fitness_level", "goals", "equipment", "schedule"],
    is_default: 1
}

template:categories -> ["training", "nutrition", "injury_prevention"]
```

### TTL设置
- 模板详情：30天 (配置更新频率较低)

---

## 6. 计划生成进度缓存

### Key格式 (用于长时间AI生成任务)
```
generation_task:{task_id} -> Hash (生成任务状态)
user:{user_id}:pending_tasks -> Set (用户待完成任务)
```

### 数据结构示例
```
generation_task:task_abc123 -> {
    task_id: "task_abc123",
    user_id: 1001,
    type: "training_plan",
    status: "generating",  // generating/completed/failed
    progress: 60,  // 百分比
    result_data: "",  // 生成结果(完成后)
    error_message: "",  // 错误信息(if failed)
    created_at: 1672502400
}

user:1001:pending_tasks -> ["task_abc123", "task_def456"]
```

### TTL设置
- 生成任务：24小时 (超时自动清理)
- 用户任务集合：永久 (通过定时任务清理)

---

## 7. 安全与恢复机制

### 缓存降级策略
当Redis不可用时，系统应能直接访问数据库，但会记录告警。

### 数据一致性
- 更新MySQL数据时同步更新Redis
- 使用Cache-Aside模式
- 关键数据设置版本号用于冲突检测

### 缓存预热
系统启动时从数据库加载：
1. 默认提示词模板
2. 活跃用户的计划和安排
3. API限流配置

---

## 8. Redis配置建议

```conf
# redis.conf关键配置
maxmemory 1gb
maxmemory-policy allkeys-lru
save 3600 1 300 100 60 1000
appendonly yes
appendfsync everysec
```

### Key前缀规范
| 模块 | 前缀 | 示例 |
|------|------|------|
| Session | session:* | session:abc123 |
| API | api_limit:* | api_limit:1001:minute |
| 计划 | training_plan:* / nutrition_plan:* | training_plan:1001 |
| 模板 | template:* | template:training:1001 |
| 生成任务 | generation_task:* | generation_task:task_abc123 |

---

## 9. 监控指标
通过Redis命令监控：
```bash
INFO keyspace   # 查看各DB的key数量
INFO memory     # 内存使用
INFO stats      # 统计信息
SLOWLOG GET 10  # 慢查询
```

### 重要业务指标
- 缓存命中率 (query_cache_hits / query_cache_misses)
- API限流触发次数
- 活跃Session数
