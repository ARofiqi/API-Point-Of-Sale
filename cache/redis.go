package cache

import (
	"aro-shop/config"
	"context"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	RedisClient   *redis.Client
	ctx           = context.Background()
	cfg           = config.LoadConfig()
	redisHost     = cfg.REDISHost
	redisPort     = cfg.REDISPort
	redisPassword = cfg.REDISPass
	redisDB, _    = strconv.Atoi(cfg.REDISdb)
)

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       redisDB,
	})

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke Redis: %v", err)
	}

	log.Println("🚀 Redis terhubung!")
}

func SetCache(key string, value string, ttl time.Duration) error {
	err := RedisClient.Set(ctx, key, value, ttl).Err()
	if err != nil {
		log.Printf("❌ Gagal menyimpan cache: %v", err)
	} else {
		log.Printf("✅ Cache disimpan: key=%s, ttl=%v", key, ttl)
	}
	return err
}

func GetCache(key string) (string, error) {
	value, err := RedisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Printf("⚠️ Cache tidak ditemukan untuk key: %s", key)
		return "", err
	} else if err != nil {
		log.Printf("❌ Gagal mengambil cache: %v", err)
		return "", err
	}
	log.Printf("✅ Cache ditemukan: key=%s, value=%s", key, value)
	return value, nil
}

func DeleteCache(key string) error {
	return RedisClient.Del(ctx, key).Err()
}
