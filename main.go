package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"interface-health-check/controllers"
	"interface-health-check/database"
	"interface-health-check/models"
	"interface-health-check/scheduler"
)

func main() {

	// 加载 .env
	godotenv.Load()

	// 初始化数据库
	database.InitDB()

	// 自动建表
	database.DB.AutoMigrate(
		&models.APIInfo{},
		&models.APICheck{},
	)

	// 启动定时巡检任务
	scheduler.StartCron()

	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	// ── Prometheus 指标采集端点 ──────────────────────────────
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// ── 业务 API ────────────────────────────────────────────
	api := r.Group("/api")
	{
		api.POST("/apis", controllers.AddAPI)
		api.GET("/apis", controllers.ListAPI)
		api.DELETE("/apis/:id", controllers.DeleteAPI)
		api.GET("/checks", controllers.ListChecks)
	}

	r.GET("/dashboard", controllers.Dashboard)

	r.Run(":8080")
}
