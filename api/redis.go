package main

import (
	"context"
	"strconv"
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


func UpdateTokenUsage(ctx context.Context, client *redis.Client, token string) (bool, string) {
    tokenData, err := client.HGetAll(ctx, token).Result()
    if err != nil {
        return false, "Error retrieving token data: " + err.Error()
    }

    if len(tokenData) == 0 {
        return false, "Token not found"
    }

    err = client.Expire(ctx, token, 14*24*time.Hour).Err()
    if err != nil {
        return false, "Error setting token TTL to 2 weeks: " + err.Error()
    }


    windowTime, err := strconv.Atoi(tokenData["window"])
    if err != nil {
        return false, "Error parsing window time: " + err.Error()
    }
    
    dailyUsage, err := strconv.Atoi(tokenData["daily_usage"])
    if err != nil {
        return false, "Error parsing daily usage: " + err.Error()
    }

    currentTime := time.Now().Unix()

    if currentTime-int64(windowTime) >= 24*60*60 {
        err = client.HSet(ctx, token, "daily_usage", 0, "window", currentTime).Err()
        if err != nil {
            return false, "Error resetting token window and usage: " + err.Error()
        }
        return true, "" 
    }

    if dailyUsage < 1000 {
        err = client.HIncrBy(ctx, token, "daily_usage", 1).Err()
        if err != nil {
            return false, "Error incrementing token usage: " + err.Error()
        }
        return true, ""
    }

    return false, "Token Limit Reached"
}


