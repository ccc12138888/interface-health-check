package services

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"interface-health-check/database"
	"interface-health-check/metrics"
	"interface-health-check/models"
)

// worker 协程：从 jobs 通道取任务执行
func worker(ctx context.Context, jobs <-chan models.APIInfo, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
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

// CheckAPI 执行单次巡检，并上报 Prometheus 指标
func CheckAPI(ctx context.Context, api models.APIInfo) {

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", api.URL, nil)
	if err != nil {
		// 请求构建失败，直接记录错误
		metrics.CheckErrors.WithLabelValues(api.URL).Inc()
		metrics.CheckTotal.WithLabelValues(api.URL).Inc()
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)

	costMs := float64(time.Since(start).Milliseconds())

	check := models.APICheck{
		URL:       api.URL,
		CostTime:  int64(costMs),
		CheckedAt: time.Now(),
	}

	if err != nil {
		check.IsError = true
		check.StatusCode = 0
		// ── Prometheus ──────────────────────────────────────
		metrics.CheckErrors.WithLabelValues(api.URL).Inc()
		metrics.CheckStatusCode.WithLabelValues(api.URL, "0").Inc()
	} else {
		check.StatusCode = resp.StatusCode
		check.IsError = resp.StatusCode != 200
		resp.Body.Close()
		// ── Prometheus ──────────────────────────────────────
		if check.IsError {
			metrics.CheckErrors.WithLabelValues(api.URL).Inc()
		}
		metrics.CheckStatusCode.WithLabelValues(
			api.URL,
			fmt.Sprintf("%d", resp.StatusCode),
		).Inc()
	}

	// 响应时间 & 总次数（无论成败都记录）
	metrics.CheckDuration.WithLabelValues(api.URL).Observe(costMs)
	metrics.CheckTotal.WithLabelValues(api.URL).Inc()

	database.DB.Create(&check)
}

// RunWorkerPool 协程池：并发执行本轮所有巡检
func RunWorkerPool(parentCtx context.Context, apiList []models.APIInfo, workerCount int) {

	ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
	defer cancel()

	jobs := make(chan models.APIInfo, len(apiList))
	var wg sync.WaitGroup

	for i := 0; i < workerCount; i++ {
		go worker(ctx, jobs, &wg)
	}

	for _, api := range apiList {
		wg.Add(1)
		jobs <- api
	}

	close(jobs)
	wg.Wait()
}
