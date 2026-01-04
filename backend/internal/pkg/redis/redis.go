package redis

import (
	"context"
	"fmt"
	"github.com/ai-fitness-planner/backend/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

var Rdb *redis.Client
var ctx = context.Background()

func InitRedis() error {
	redisCfg := config.GlobalConfig.Database.Redis

	Rdb = redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:   redisCfg.Password,
		DB:         redisCfg.DB,
		PoolSize:   redisCfg.PoolSize,
		MaxRetries: redisCfg.MaxRetries,
	})

	// 测试连接
	_, err := Rdb.Ping(ctx).Result()
	if err != nil {
		return fmt.Errorf("Redis连接失败: %w", err)
	}

	return nil
}

func Close() error {
	if Rdb != nil {
		return Rdb.Close()
	}
	return nil
}

// Session操作
func SetSession(sessionID string, userID int64, ttl time.Duration) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return Rdb.Set(ctx, key, userID, ttl).Err()
}

func GetSession(sessionID string) (int64, error) {
	key := fmt.Sprintf("session:%s", sessionID)
	userID, err := Rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil // Session不存在
	}
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func DeleteSession(sessionID string) error {
	key := fmt.Sprintf("session:%s", sessionID)
	return Rdb.Del(ctx, key).Err()
}

// API限流
func CheckRateLimit(key string, limit int64, duration time.Duration) (bool, error) {
	current, err := Rdb.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	if current == 1 {
		Rdb.Expire(ctx, key, duration)
	}

	return current <= limit, nil
}

// Plan生成任务状态
func SetPlanTask(taskID string, data interface{}, ttl time.Duration) error {
	key := fmt.Sprintf("plan_task:%s", taskID)
	return Rdb.Set(ctx, key, data, ttl).Err()
}

func GetPlanTask(taskID string) (string, error) {
	key := fmt.Sprintf("plan_task:%s", taskID)
	return Rdb.Get(ctx, key).Result()
}

func DeletePlanTask(taskID string) error {
	key := fmt.Sprintf("plan_task:%s", taskID)
	return Rdb.Del(ctx, key).Err()
}

// 缓存操作
func SetCache(key string, value interface{}, ttl time.Duration) error {
	return Rdb.Set(ctx, key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	return Rdb.Get(ctx, key).Result()
}

func DeleteCache(key string) error {
	return Rdb.Del(ctx, key).Err()
}

// AI API调用次数统计
func IncrementAPICall(userID int64, apiID int64) error {
	now := time.Now()

	// 分钟级别
	minuteKey := fmt.Sprintf("api_calls:minute:%d:%d:%s", userID, apiID, now.Format("200601021504"))
	if err := Rdb.Incr(ctx, minuteKey).Err(); err != nil {
		return err
	}
	Rdb.Expire(ctx, minuteKey, time.Hour)

	// 小时级别
	hourKey := fmt.Sprintf("api_calls:hour:%d:%d:%s", userID, apiID, now.Format("2006010215"))
	if err := Rdb.Incr(ctx, hourKey).Err(); err != nil {
		return err
	}
	Rdb.Expire(ctx, hourKey, 24*time.Hour)

	// 天级别
	dayKey := fmt.Sprintf("api_calls:day:%d:%d:%s", userID, apiID, now.Format("20060102"))
	if err := Rdb.Incr(ctx, dayKey).Err(); err != nil {
		return err
	}
	Rdb.Expire(ctx, dayKey, 7*24*time.Hour)

	return nil
}

func GetAPICallCount(userID int64, apiID int64, period string) (int64, error) {
	now := time.Now()
	var key string

	switch period {
	case "minute":
		key = fmt.Sprintf("api_calls:minute:%d:%d:%s", userID, apiID, now.Format("200601021504"))
	case "hour":
		key = fmt.Sprintf("api_calls:hour:%d:%d:%s", userID, apiID, now.Format("2006010215"))
	case "day":
		key = fmt.Sprintf("api_calls:day:%d:%d:%s", userID, apiID, now.Format("20060102"))
	default:
		return 0, fmt.Errorf("不支持的时间段: %s", period)
	}

	count, err := Rdb.Get(ctx, key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}
