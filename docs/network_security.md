# 网络安全检查与加固方案

## 一、网络架构安全设计

```
┌──────────────────────────────────────────────────────────────┐
│                       Internet                                │
│               ┌──────────────────────────┐                      │
│               │   WAF/IPS防护层          │                      │
│               │  (Cloudflare/AWS WAF)    │                      │
│               └───────┬──────────────────┘                      │
│                       │                                       │
│               ┌───────▼──────────────────┐                      │
│               │   DDoS防护               │                      │
│               │  (限速/清洗)             │                      │
│               └───────┬──────────────────┘                      │
│                       │                                       │
├───────────────────────┼───────────────────────────────────────┤
│                       │                                       │
│               ┌───────▼──────────────────┐                      │
│               │    负载均衡              │                      │
│               │   (Nginx/ALB)            │                      │
│               └──┬────────┬────────┬─────┘                      │
│                  │        │        │                            │
│          ┌───────▼──┐ ┌──▼────┐ ┌▼──────┐                       │
│          │  应用服务 │ │ 应用服务│ │应用服务│                       │
│          │   实例1   │ │  实例2 │ │ 实例3 │                       │
│          └──────┬───┘ └───┬────┘ └──┬────┘                       │
│                 │          │         │                            │
├─────────────────┼──────────┼─────────┼────────────────────────────┤
│                 │          │         │                            │
│          ┌──────▼──────────▼─────────▼──────┐                     │
│          │        内部网络 (VPC)            │                     │
│          │    ┌────────────────────┐        │                     │
│          │    │    Redis集群       │        │                     │
│          │    │   (6节点Cluster)   │        │                     │
│          │    └─────────┬──────────┘        │                     │
│          │              │                   │                     │
│          │    ┌─────────▼──────────┐        │                     │
│          │    │   MySQL集群        │        │                     │
│          │    │ (主从+Proxy)       │        │                     │
│          │    └────────────────────┘        │                     │
│          └───────────────────────────────────┘                     │
└──────────────────────────────────────────────────────────────┘
```

---

## 二、传输层安全 (TLS/SSL)加固

### 1. Nginx TLS配置最佳实践

```nginx
# /etc/nginx/conf.d/ssl.conf
server {
    listen 443 ssl http2;
    listen [::]:443 ssl http2;
    server_name api.fitness.example.com;

    # SSL证书
    ssl_certificate /etc/nginx/ssl/fullchain.pem;
    ssl_certificate_key /etc/nginx/ssl/privkey.pem;

    # TLS版本限制 (禁用旧版本)
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;

    # 强加密套件
    ssl_ciphers 'ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:
                 ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:
                 ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305';

    # Session缓存
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;

    # OCSP Stapling
    ssl_stapling on;
    ssl_stapling_verify on;

    # 安全头部
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains; preload" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # 禁用不安全的HTTP方法
    if ($request_method !~ ^(GET|HEAD|POST|PUT|DELETE|OPTIONS)$) {
        return 405;
    }

    location / {
        proxy_pass http://backend_servers;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Request-ID $request_id;

        # 超时设置
        proxy_connect_timeout 5s;
        proxy_send_timeout 30s;
        proxy_read_timeout 30s;

        # 缓冲
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;
        proxy_busy_buffers_size 8k;

        # WebSocket支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}

# HTTP重定向到HTTPS
server {
    listen 80;
    listen [::]:80;
    server_name api.fitness.example.com;
    return 301 https://$server_name$request_uri;
}
```

### 2. SSL Labs评分优化

```bash
# 检查SSL配置评分
docker run --rm -it nablac0d3/ssllabs-scan api.fitness.example.com

# 预期结果: A+ 评分
```

#### 优化清单:
- [x] 启用TLS 1.2+
- [x] 禁用弱加密套件
- [x] 启用HSTS
- [x] 配置OCSP Stapling
- [x] 禁用SSL Session Tickets
- [x] 配置安全的密码套件
- [x] 使用Perfect Forward Secrecy

---

## 三、API网关安全配置

### 1. Kong API Gateway配置

```yaml
# kong.yml
_format_version: "2.1"

services:
  - name: fitness-api
    url: http://backend:8080
    routes:
      - name: fitness-api-route
        paths:
          - /api/v1
        strip_path: false

plugins:
  # 限流插件
  - name: rate-limiting
    service: fitness-api
    config:
      second: 10
      minute: 100
      hour: 1000
      day: 10000
      policy: redis
      redis_host: redis-cluster
      fault_tolerant: true
      hide_client_headers: false

  # JWT认证
  - name: jwt
    service: fitness-api
    config:
      secret_is_base64: false
      run_on_preflight: false

  # CORS
  - name: cors
    service: fitness-api
    config:
      origins:
        - "https://fitness.example.com"
      methods:
        - GET
        - POST
        - PUT
        - DELETE
        - OPTIONS
      headers:
        - Accept
        - Accept-Version
        - Content-Length
        - Content-MD5
        - Content-Type
        - Date
        - Authorization
      credentials: true
      max_age: 3600

  # 请求大小限制
  - name: request-size-limiting
    service: fitness-api
    config:
      allowed_payload_size: 8388608  # 8MB
      size_unit: bytes

  # 请求转换 (移除敏感信息)
  - name: response-transformer
    service: fitness-api
    config:
      remove:
        headers:
          - X-Powered-By
          - Server

  # IP白名单 (管理接口)
  - name: ip-restriction
    route: admin-route
    config:
      whitelist:
        - 10.0.0.0/8
        - 172.16.0.0/12
        - 192.168.0.0/16
```

### 2. Golang后端中间件安全栈

```go
// 安全的中间件组合
func SetupRouter(cfg *config.Config) *gin.Engine {
    r := gin.New()

    // 1. 全局中间件（顺序重要）
    r.Use(gin.Recovery())                      // 恢复panic
    r.Use(middleware.Logger())                 // 结构化日志
    r.Use(middleware.SecurityHeaders())        // 安全头部
    r.Use(middleware.RequestID())              // 请求追踪
    r.Use(middleware.Timeout(30 * time.Second)) // 请求超时

    // 2. CORS配置
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"https://fitness.example.com"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Request-ID"},
        ExposeHeaders:    []string{"Content-Length", "X-Total-Count", "X-Rate-Limit-Remaining"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // 3. 限流 (IP和用户级别)
    r.Use(middleware.RateLimitIP(100, time.Minute))

    // 4. 路由组
    api := r.Group("/api/v1")
    {
        // 公共路由
        public := api.Group("")
        {
            public.POST("/auth/login", authHandler.Login)
            public.POST("/auth/register", authHandler.Register)
        }

        // 需要认证的路由
        protected := api.Group("")
        protected.Use(middleware.JWTAuth(), middleware.RateLimitUser(1000, time.Hour))
        {
            // 用户相关
            protected.GET("/user/profile", userHandler.GetProfile)
            protected.PUT("/user/profile", userHandler.UpdateProfile)

            // AI配置
            protected.POST("/ai-apis", aiHandler.AddAPI)
            protected.GET("/ai-apis", aiHandler.ListAPIs)
            protected.DELETE("/ai-apis/:id", aiHandler.DeleteAPI)

            // 训练计划
            plan := protected.Group("/training-plans")
            {
                plan.POST("/generate", trainingHandler.GeneratePlan)
                plan.GET("/today", trainingHandler.GetTodayTraining)
            }
        }
    }

    return r
}
```

---

## 四、DDoS防护策略

### 1. 多层防御架构

```
流量清洗流程:

正常流量:   恶意流量:   可疑流量:
   │           │            │
   │           │            │
   ▼           ▼            ▼
┌────────────────────────────────┐
│     CDN/云防护 (L3/L4)         │
│  (Cloudflare/AWS Shield)       │
└────────┬──────────────┬────────┘
         │              │
    ┌────▼────┐    ┌───▼─────────┐
    │  放行    │    │  清洗中心    │
    └────┬────┘    └────┬────────┘
         │              │
         └───────┬──────┘
                 │
         ┌───────▼────────┐
         │  源站服务器     │
         │  限流/验证      │
         └────────────────┘
```

### 2. 基于Nginx的限流配置

```nginx
# Rate Limit Zones
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=10r/s;
limit_req_zone $binary_remote_addr zone=auth_limit:10m rate=5r/m;
limit_req_zone $binary_remote_addr zone=plan_limit:10m rate=2r/m;

limit_conn_zone $binary_remote_addr zone=addr:10m;

# 应用限流
server {
    location /api/v1/ {
        limit_req zone=api_limit burst=20 nodelay;
        limit_conn addr 10;
    }

    location /api/v1/auth/ {
        limit_req zone=auth_limit burst=5 nodelay;
    }

    location /api/v1/training-plans/generate {
        limit_req zone=plan_limit burst=2 nodelay;
    }

    # 错误页面
    error_page 429 @too_many_requests;
    location @too_many_requests {
        return 429 'Rate limit exceeded';
        add_header Content-Type text/plain;
    }
}

# 请求大小限制
client_body_buffer_size 1M;
client_max_body_size 8M;
client_header_buffer_size 1k;
large_client_header_buffers 2 4k;
```

### 3. 基于Go的DDoS防护中间件

```go
// DDoS防护中间件
type DDoSConfig struct {
    RequestThreshold int           // 请求阈值
    TimeWindow       time.Duration // 时间窗口
    BlockDuration    time.Duration // 封禁时长
}

type DDoSMiddleware struct {
    config *DDoSConfig
    cache  *cache.Cache
    logger *zap.Logger
}

func NewDDoSMiddleware(config *DDoSConfig) gin.HandlerFunc {
    m := &DDoSMiddleware{
        config: config,
        cache:  cache.New(config.TimeWindow*2, config.TimeWindow*2),
        logger: logger.GetLogger(),
    }

    return m.Handle()
}

func (m *DDoSMiddleware) Handle() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()

        // 检查是否在封禁列表
        if m.isBlocked(clientIP) {
            c.JSON(403, gin.H{
                "error": "Access temporarily blocked",
                "retry_after": m.getBlockRemaining(clientIP),
            })
            c.Abort()
            return
        }

        // 统计请求数
        key := fmt.Sprintf("ddos:%s", clientIP)
        count := m.incrementRequestCount(key)

        // 超过阈值，加入封禁列表
        if count > m.config.RequestThreshold {
            m.blockIP(clientIP)
            m.logger.Warn("DDoS attack detected",
                zap.String("ip", clientIP),
                zap.Int("requests", count),
            )
            c.JSON(429, gin.H{
                "error": "Too many requests",
                "retry_after": m.config.BlockDuration.Seconds(),
            })
            c.Abort()
            return
        }

        c.Next()
    }
}

func (m *DDoSMiddleware) isBlocked(ip string) bool {
    key := fmt.Sprintf("block:%s", ip)
    _, found := m.cache.Get(key)
    return found
}

func (m *DDoSMiddleware) blockIP(ip string) {
    key := fmt.Sprintf("block:%s", ip)
    m.cache.Set(key, true, m.config.BlockDuration)
}

func (m *DDoSMiddleware) incrementRequestCount(key string) int {
    count, _ := m.cache.IncrementInt(key, 1)
    return count
}
```

---

## 五、入侵检测与监控

### 1. 安全事件监控系统

```go
// 安全事件类型
type SecurityEventType string

const (
    EventLoginFailed       SecurityEventType = "login_failed"
    EventSuspiciousRequest SecurityEventType = "suspicious_request"
    EventRateLimitExceeded SecurityEventType = "rate_limit_exceeded"
    EventSQLInjection      SecurityEventType = "sql_injection_detected"
    EventXSSAttempt        SecurityEventType = "xss_attempt"
    EventDataExfiltration  SecurityEventType = "data_exfiltration"
)

type SecurityMonitor struct {
    alertManager AlertManager
    metrics      *prometheus.Metrics
}

func (m *SecurityMonitor) LogEvent(ctx context.Context, event SecurityEvent) {
    // 写入日志
    logger.GetLogger().Error("Security event detected",
        zap.String("type", string(event.Type)),
        zap.String("ip", event.SourceIP),
        zap.Int64("user_id", event.UserID),
        zap.String("details", event.Details),
    )

    // 更新监控指标
    m.metrics.SecurityEventCounter.WithLabelValues(
        string(event.Type),
        event.SourceIP,
    ).Inc()

    // 触发告警
    if m.shouldAlert(event) {
        m.alertManager.SendAlert(event)
    }
}

func (m *SecurityMonitor) shouldAlert(event SecurityEvent) bool {
    // 告警规则
    switch event.Type {
    case EventLoginFailed:
        // 同一IP 5分钟内失败超过10次
        return m.countEvents(event.SourceIP, EventLoginFailed, 5*time.Minute) > 10
    case EventSQLInjection:
        // 立即告警
        return true
    default:
        return false
    }
}
```

### 2. Prometheus监控指标

```go
// 定义安全监控指标
var (
    SecurityEventCounter = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "security_events_total",
            Help: "Total number of security events",
        },
        []string{"type", "source_ip"},
    )

    RateLimitHits = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "rate_limit_hits_total",
            Help: "Total number of rate limit hits",
        },
        []string{"user_id", "endpoint"},
    )

    FailedLogins = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "failed_logins_total",
            Help: "Total number of failed login attempts",
        },
        []string{"ip", "username"},
    )

    BlockedIPs = prometheus.NewGauge(
        prometheus.GaugeOpts{
            Name: "blocked_ips_total",
            Help: "Number of currently blocked IPs",
        },
    )
)
```

---

## 六、防火墙规则

### 1. 应用层防火墙 (iptables)

```bash
#!/bin/bash
# 应用层防火墙配置

# 清空现有规则
iptables -F
iptables -X

# 默认策略
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT

# 允许本地回环
iptables -A INPUT -i lo -j ACCEPT
iptables -A OUTPUT -o lo -j ACCEPT

# 允许已建立的连接
iptables -A INPUT -m state --state ESTABLISHED,RELATED -j ACCEPT

# 允许SSH (仅特定IP或限制频率)
iptables -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --set
iptables -A INPUT -p tcp --dport 22 -m state --state NEW -m recent --update --seconds 60 --hitcount 4 -j DROP
iptables -A INPUT -p tcp --dport 22 -j ACCEPT

# 允许HTTP/HTTPS
iptables -A INPUT -p tcp --dport 80 -j ACCEPT
iptables -A INPUT -p tcp --dport 443 -j ACCEPT

# Docker容器间通信
iptables -A INPUT -i docker0 -j ACCEPT
iptables -A FORWARD -i docker0 -o docker0 -j ACCEPT

# 保存规则
iptables-save > /etc/iptables/rules.v4
```

### 2. Fail2Ban配置 (防止暴力破解)

```ini
# /etc/fail2ban/jail.local
[nginx-http-auth]
enabled = true
filter = nginx-http-auth
port = http,https
logpath = /var/log/nginx/error.log
maxretry = 5
bantime = 3600
findtime = 600

[nginx-limit-req]
enabled = true
filter = nginx-limit-req
port = http,https
logpath = /var/log/nginx/error.log
maxretry = 10
bantime = 7200
findtime = 600

[fitness-api]
enabled = true
filter = fitness-api
port = http,https
logpath = /var/log/fitness/app.log
maxretry = 20
bantime = 86400
findtime = 3600
```

---

## 七、安全测试清单

### 1. 渗透测试检查表

- [x] SSL Labs评分 (A+)
- [x] 端口扫描 (仅暴露必要端口)
- [x] SQL注入测试 (使用sqlmap)
- [x] XSS测试 (使用XSStrike)
- [x] CSRF测试
- [x] 认证绕过测试
- [x] 权限提升测试
- [x] 会话固定测试

```bash
# 自动化安全测试
# 使用OWASP ZAP
docker run -t owasp/zap2docker-stable zap-baseline.py \
  -t https://api.fitness.example.com

# 使用Nikto
nikto -h https://api.fitness.example.com
```

### 2. 性能测试 (DoS防护)

```bash
# 压力测试工具
# ab (ApacheBench)
ab -n 10000 -c 100 https://api.fitness.example.com/api/v1/health

# JMeter
# 配置线程组: 1000线程, 持续60秒

# 结果指标:
# - 成功率 > 99.9%
# - 响应时间 P99 < 1s
# - 错误率 < 0.1%
```

---

## 八、应急响应方案

### 1. 事件响应流程

```
安全事件发生
    ↓
自动检测 (SIEM/IDS)
    ↓
告警通知 (Slack/Email/SMS)
    ↓
初步评估
    ↓
┌─── 低危 ───┐  ┌─── 中危 ───┐  ┌─── 高危 ───┐
│  记录分析  │  │  限制访问  │  │  立即隔离  │
│  监控观察  │  │  加强监控  │  │  启动预案  │
└────────────┘  └────────────┘  └────────────┘
    ↓                ↓                ↓
恢复服务        调查分析         灾难恢复
    ↓                ↓                ↓
事后复盘        事后复盘         事后复盘
```

### 2. 常见场景处理

#### DDoS攻击
```bash
# 快速启用Cloudflare "Under Attack"模式
curl -X PATCH "https://api.cloudflare.com/client/v4/zones/{zone_id}/settings/security_level" \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  --data '{"value":"under_attack"}'

# 或手动添加iptables规则封禁IP段
iptables -A INPUT -s {attack_ip}/24 -j DROP
```

#### 数据泄露
1. 立即切断数据库外网访问
2. 评估泄露范围和影响
3. 通知受影响用户
4. 重置相关密钥
5. 加强访问控制

#### 系统入侵
1. 隔离受影响服务器
2. 保存现场日志
3. 分析入侵路径
4. 修复漏洞
5. 恢复系统

---

## 九、安全加固总结

### 实施优先级

**P0 - 立即实施:**
- [ ] TLS 1.2+ 强制启用
- [ ] HSTS头部配置
- [ ] API限流
- [ ] 密码强度验证
- [ ] SQL注入防护

**P1 - 本周实施:**
- [ ] WAF部署
- [ ] DDoS防护
- [ ] 日志监控
- [ ] 定期备份
- [ ] 安全审计

**P2 - 本月实施:**
- [ ] 多地部署
- [ ] 灾备演练
- [ ] 渗透测试
- [ ] 安全培训

### 安全评分目标

| 项目 | 当前 | 目标 | 工具 |
|------|------|------|------|
| SSL配置 | - | A+ | SSL Labs |
| 安全头部 | - | 100% | Security Headers |
| 漏洞扫描 | - | 0高危 | OWASP ZAP |
| 依赖安全 | - | 0漏洞 | Snyk/Dependabot |

### 持续监控

- [部署安全监控看板](https://grafana.com/)
- [启用异常告警](https://prometheus.io/)
- [定期安全扫描](https://www.owasp.org/) (每周)
- [依赖漏洞检查](https://snyk.io/) (每天)
- [日志审计](https://www.elastic.co/) (实时)

---

**最后更新: 2024-01-04**
**文档状态: 已完成 ✓**
