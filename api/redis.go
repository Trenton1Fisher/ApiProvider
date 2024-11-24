package main

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func NewRedisClient(redisAddr string, redisDB int , redisPass string) (*redis.Client, error){
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