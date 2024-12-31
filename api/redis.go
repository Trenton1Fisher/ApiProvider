package main

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(redisAddr string, redisDB int , redisPass string) (*redis.Client, error){
    ctx := context.Background()
	client := redis.NewClient(&redis.Options{
        Addr:	  redisAddr,
        Password: redisPass,
        DB:		  redisDB,
        Protocol: 2,
    })

    _ , err := client.Ping(ctx).Result()
    if err != nil{
        return nil, err
    }

    return client, nil
}

func AddNewToken(ctx context.Context, client *redis.Client, token string) (bool, error) {
    err := client.HSet(ctx, token, "daily_usage", 0, "window", 0).Err()
    if err != nil {
        return false, err
    }

    err = client.Expire(ctx, token, 14*24*time.Hour).Err()
    if err != nil {
        return false, err
    }

    return true, nil
}

func CheckIfTokenExists(ctx context.Context, client *redis.Client, token string) (bool, error) {
    exists, err := client.Exists(ctx, token).Result()
    if err != nil {
        return false, err
    }
    return exists > 0, nil
}

func UpdateTokenWindow(ctx context.Context, client *redis.Client, token string) {
}

func UpdateTokenUsageCount(ctx context.Context, client *redis.Client, token string){
}


