package services

import (
	"context"
	"net/http"
	"sync"
	"time"

	"interface-health-check/database"
	"interface-health-check/models"
)

// worker 协程
func worker(ctx context.Context, jobs <-chan models.APIInfo, wg *sync.WaitGroup) {

	for {
		select {
		case <-ctx.Done():
			// 收到取消信号
			return

		case api, ok := <-jobs:
			if !ok {
				return
			}

			CheckAPI(ctx, api)
			wg.Done()
		}
	}
}

// 真正执行巡检
func CheckAPI(ctx context.Context, api models.APIInfo) {

	start := time.Now()

	// 创建带 context 的请求
	req, err := http.NewRequestWithContext(ctx, "GET", api.URL, nil)
	if err != nil {
		return
	}

	client := &http.Client{}

	resp, err := client.Do(req)

	cost := time.Since(start).Milliseconds()

	check := models.APICheck{
		URL:       api.URL,
		CostTime:  cost,
		CheckedAt: time.Now(),
	}

	if err != nil {
		check.IsError = true
		check.StatusCode = 0
	} else {
		check.StatusCode = resp.StatusCode
		check.IsError = resp.StatusCode != 200
		resp.Body.Close()
	}

	database.DB.Create(&check)
}

// 协程池控制并发
func RunWorkerPool(parentCtx context.Context, apiList []models.APIInfo, workerCount int) {

	// 设置本轮巡检最大执行时间（比如 10 秒）
	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	jobs := make(chan models.APIInfo, len(apiList))
	var wg sync.WaitGroup

	// 启动 worker
	for i := 0; i < workerCount; i++ {
		go worker(ctx, jobs, &wg)
	}

	// 投递任务
	for _, api := range apiList {
		wg.Add(1)
		jobs <- api
	}

	close(jobs)

	wg.Wait()
}