package config

import (
	"fmt"
	"os"
	"strconv"
)

// 读取环境变量
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// 生成 MySQL DSN
func GetDSN() string {

	user := GetEnv("DB_USER", "root")
	password := GetEnv("DB_PASSWORD", "")
	host := GetEnv("DB_HOST", "127.0.0.1")
	port := GetEnv("DB_PORT", "3306")
	dbname := GetEnv("DB_NAME", "healthcheck")

	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		password,
		host,
		port,
		dbname,
	)
}


 //新增：支持从 .env 读取 MaxWorker
func GetMaxWorker() int {
	valueStr := GetEnv("MAX_WORKER", "5")

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 5
	}
	return value
}

//新增：支持从 .env 读取 CronSpec
func GetCronSpec() string {
	return GetEnv("CRON_SPEC", "@every 1m")
}