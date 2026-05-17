package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/robfig/cron/v3"

	"interface-health-check/config"
	"interface-health-check/database"
	"interface-health-check/models"
	"interface-health-check/services"
)

// 定义一个全局锁
var mu sync.Mutex

// 定义运行状态
var isRunning bool

// StartCron 启动定时巡检（企业级防重叠版本）
func StartCron() *cron.Cron {

	spec := config.GetCronSpec()

	c := cron.New(
		cron.WithSeconds(),
	)

	_, err := c.AddFunc(spec, func() {

		// 1️⃣ 尝试获取锁
		mu.Lock()

		if isRunning {
			// 如果上一轮还没执行完，直接跳过
			log.Println("⚠ 上一次巡检仍在执行，本次跳过")
			mu.Unlock()
			return
		}

		// 标记为正在运行
		isRunning = true
		mu.Unlock()

		log.Println("🚀 开始执行定时巡检任务")

		// 2️⃣ 查询接口
		var apis []models.APIInfo
		database.DB.Find(&apis)

		if len(apis) == 0 {
			log.Println("⚠ 没有需要巡检的接口")

			mu.Lock()
			isRunning = false
			mu.Unlock()
			return
		}

		// 3️⃣ 创建带超时的 context
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// 4️⃣ 执行协程池
		services.RunWorkerPool(ctx, apis, config.GetMaxWorker())

		log.Println("✅ 本轮巡检完成")

		// 5️⃣ 执行结束，释放状态
		mu.Lock()
		isRunning = false
		mu.Unlock()
	})

	if err != nil {
		log.Fatal("添加 cron 任务失败:", err)
	}

	c.Start()

	log.Println("⏰ 定时任务已启动，规则:", spec)

	return c
}