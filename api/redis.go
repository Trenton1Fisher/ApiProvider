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

func AddNewToken(client *redis.Client, token string) (bool, error) {
    ctx := context.Background()

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

func CheckIfTokenExists(client *redis.Client, token string) (bool, error) {
    ctx := context.Background()
    exists, err := client.HExists(ctx, token, "daily_usage").Result()
    if err != nil {
        return false, err
    }
    return exists, nil
}