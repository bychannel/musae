package service

import (
	"errors"
)

type RedisDbType string

const (
	RedisLock    RedisDbType = "redis-lock"
	RedisGlobal  RedisDbType = "redis-global"
	RedisCache   RedisDbType = "redis-cache"
	ConfigCenter RedisDbType = "config-center"
)

const (
	REDIS_TTL_NAME = "ttlInSeconds"
)

type MongoDbType string

const (
	MongoDbType_MongoNil     MongoDbType = ""
	MongoDbType_MongoAccount MongoDbType = "mongo-account"
	MongoDbType_MongoGame    MongoDbType = "mongo-game"
	MongoDbType_MongoMail    MongoDbType = "mongo-mail"
	MongoDbType_MongoGmt     MongoDbType = "mongo-gmt"
)

var (
	DB_ERROR_TIMEOUT   = errors.New("DB_ERROR_TIMEOUT")   // 连接db超时
	DB_ERROR_NOT_EXIST = errors.New("DB_ERROR_NOT_EXIST") // 数据不存在
	DB_ERROR_MARSHAL   = errors.New("DB_ERROR_MARSHAL")   // 序列化错误
	DB_ERROR_UNMARSHAL = errors.New("DB_ERROR_UNMARSHAL") // 反序列化错误
)
