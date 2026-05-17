# Interface Health Check — 接口巡检服务

定时巡检 HTTP 接口，记录响应时间与状态，通过 Prometheus + Grafana 可视化监控。

## 项目结构

```
interface-health-check/
├── config/         # 环境变量读取（DSN、MaxWorker、CronSpec）
├── database/       # GORM + MySQL 初始化
├── models/         # APIInfo / APICheck 数据模型
├── controllers/    # HTTP 控制层（Gin）
├── services/       # 巡检逻辑 + 协程池
├── scheduler/      # 定时任务（防重叠锁）
├── metrics/        # Prometheus 指标定义
├── templates/      # Dashboard HTML 模板
├── monitoring/
│   ├── prometheus.yml                          # Prometheus 采集配置
│   └── grafana/
│       ├── provisioning/datasources/           # 数据源自动注入
│       ├── provisioning/dashboards/            # 面板目录自动注入
│       └── dashboards/healthcheck.json         # 预置监控面板
├── Dockerfile
├── docker-compose.yml
├── .env.example
└── main.go
```

## 快速启动

```bash
# 1. 复制环境变量
cp .env.example .env

# 2. 一键启动所有服务
docker compose up -d --build

# 3. 访问各服务
#   应用 Dashboard : http://localhost:8080/dashboard
#   Prometheus     : http://localhost:9090
#   Grafana        : http://localhost:3000  (admin / admin)
```

## API 接口

| 方法   | 路径              | 说明             |
|--------|-------------------|------------------|
| POST   | /api/apis         | 添加巡检接口     |
| GET    | /api/apis         | 查看接口列表     |
| DELETE | /api/apis/:id     | 删除接口         |
| GET    | /api/checks       | 查看最近巡检结果 |
| GET    | /dashboard        | 可视化 Dashboard |
| GET    | /metrics          | Prometheus 指标  |

### 添加接口示例

```bash
curl -X POST http://localhost:8080/api/apis \
  -H "Content-Type: application/json" \
  -d '{"name":"百度","url":"https://www.baidu.com","method":"GET"}'
```

## Prometheus 指标说明

| 指标名                          | 类型      | 说明                   |
|---------------------------------|-----------|------------------------|
| `healthcheck_total`             | Counter   | 巡检总次数（按 URL）   |
| `healthcheck_errors_total`      | Counter   | 错误总次数（按 URL）   |
| `healthcheck_duration_ms`       | Histogram | 响应时间分布（毫秒）   |
| `healthcheck_status_code_total` | Counter   | HTTP 状态码分布        |

### 常用 PromQL 查询

```promql
# 各接口过去 5 分钟错误率
rate(healthcheck_errors_total[5m]) / rate(healthcheck_total[5m]) * 100

# P95 响应时间
histogram_quantile(0.95, rate(healthcheck_duration_ms_bucket[5m]))

# 每分钟巡检频率
rate(healthcheck_total[1m]) * 60
```

## 环境变量

| 变量名            | 默认值        | 说明                              |
|-------------------|---------------|-----------------------------------|
| `DB_USER`         | root          | MySQL 用户名                     |
| `DB_PASSWORD`     | （空）        | MySQL 密码                       |
| `DB_HOST`         | 127.0.0.1     | MySQL 地址                       |
| `DB_PORT`         | 3306          | MySQL 端口                       |
| `DB_NAME`         | healthcheck   | 数据库名                         |
| `CRON_SPEC`       | @every 1m     | 巡检 cron 表达式（支持秒级）     |
| `MAX_WORKER`      | 5             | 协程池并发数                     |
| `GRAFANA_USER`    | admin         | Grafana 管理员账号               |
| `GRAFANA_PASSWORD`| admin         | Grafana 管理员密码               |

## 修改 checker.go

将原 `services/checker.go` 替换为本仓库的 `checker_updated.go`，在每次巡检结束时自动上报指标到 Prometheus。

## 修改 main.go

将原 `main.go` 替换为 `main_updated.go`，新增 `/metrics` 路由暴露 Prometheus 采集端点。

## 新增依赖

```bash
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```
