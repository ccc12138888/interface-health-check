package models

import "time"

// 需要巡检的接口
type APIInfo struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string
	URL       string
	Method    string
	CreatedAt time.Time
}

// 巡检结果记录
type APICheck struct {
	ID         uint      `gorm:"primaryKey"`
	URL        string
	StatusCode int
	CostTime   int64
	IsError    bool
	CheckedAt  time.Time
}