package models

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"gopkg.in/ini.v1"
	"os"
	"time"
)

var ctx = context.Background()
var rdbClient *redis.Client
var redisEnable bool

func init() {
	config, iniErr := ini.Load("./conf/app.ini")
	if iniErr != nil {
		fmt.Printf("Fail to read file: %v", iniErr)
		os.Exit(1)
	}

	ip := config.Section("redis").Key("ip").String()
	port := config.Section("redis").Key("port").String()
	password := config.Section("redis").Key("password").String()
	database := config.Section("redis").Key("database").String()
	myDatabase, _ := Int(database)
	redisEnable, _ = config.Section("redis").Key("redisEnable").Bool()

	addr := fmt.Sprintf("%v:%v", ip, port)

	if redisEnable {
		rdbClient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       myDatabase,
		})

		_, err := rdbClient.Ping(ctx).Result()

		if err == nil {
			fmt.Println("链接成功！")
		} else {
			fmt.Println("链接失败！")
		}
	}
}

type cacheDb struct{}

func (c cacheDb) Set(key string, value interface{}, expiration int) {
	if redisEnable {
		v, err := json.Marshal(value)
		if err == nil {
			rdbClient.Set(ctx, key, string(v), time.Second*time.Duration(expiration))
		}
	}
}

func (c cacheDb) Get(key string, obj interface{}) bool {
	if redisEnable {
		valueStr, err1 := rdbClient.Get(ctx, key).Result()
		if err1 == nil && valueStr != "" {
			err2 := json.Unmarshal([]byte(valueStr), obj)
			return err2 == nil
		}
		return false
	}
	return false
}

func (c cacheDb) FlushAll() {
	if redisEnable {
		rdbClient.FlushAll(ctx)
	}
}

var CacheDb = &cacheDb{}
