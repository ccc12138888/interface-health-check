好的 👌
下面是一份 **完整、专业、可直接放到 GitHub 的 README.md**，完全基于你现在的项目能力（Gin + Gorm + Cron + 协程池 + Context + 防重叠执行 + Dashboard），**没有乱写不存在的功能**。

你可以直接复制为 `README.md` 使用。

---

# 🚀 Interface Health Check

> 基于 Go 构建的接口健康巡检微服务
> 支持定时探活、协程池并发控制、异常记录与 Web 可视化展示。

---

## 📌 项目简介

Interface Health Check 是一个轻量级接口健康监控系统。

系统可以：

* 动态添加 / 删除需要巡检的接口
* 定时并发检测接口可用性
* 记录状态码与响应耗时
* 自动标记异常请求
* 提供 Web 仪表盘展示巡检统计
* 防止定时任务重叠执行
* 使用 Context 控制超时与取消

该项目模拟企业内部接口监控系统的核心能力，重点体现：

* 并发控制能力
* 定时任务调度能力
* RESTful API 设计能力
* 分层架构设计能力
* 数据持久化能力
* 微服务工程实践能力

---

## 🏗 系统架构

```
Cron 定时任务
        ↓
查询数据库接口列表
        ↓
Worker Pool 并发执行 HTTP 探活
        ↓
Context 控制超时取消
        ↓
写入巡检日志表
        ↓
Dashboard 统计展示
```

---

## 🧠 核心特性

### ✅ 动态接口管理

* 新增巡检接口
* 删除巡检接口
* 查询接口列表

---

### ✅ 定时巡检

* 基于 robfig/cron 实现
* 支持秒级表达式
* 从 .env 读取 CRON_SPEC
* 支持防止任务重叠执行（企业级版本）

---

### ✅ 协程池并发控制

* 使用 Worker Pool 模式
* 控制最大并发数
* 从 .env 读取 MAX_WORKER
* 防止 goroutine 无限制增长

---

### ✅ Context 超时控制

* 每轮巡检设置统一超时
* 超时自动取消所有 HTTP 请求
* 避免资源泄漏

---

### ✅ 巡检日志记录

每次巡检记录：

* URL
* HTTP 状态码
* 响应耗时
* 是否异常
* 巡检时间

---

### ✅ Web 可视化 Dashboard

* 成功 / 失败饼图
* 总巡检次数统计
* 最近巡检记录

访问：

```
http://localhost:8080/dashboard
```

---

## 🛠 技术栈

| 技术              | 说明      |
| --------------- | ------- |
| Go              | 核心语言    |
| Gin             | HTTP 框架 |
| Gorm            | ORM 框架  |
| MySQL           | 数据库     |
| robfig/cron     | 定时任务    |
| Goroutine Pool  | 并发控制    |
| Context         | 超时与取消   |
| HTML + Chart.js | 可视化展示   |

---

## 📂 项目结构

```
interface-health-check/
│
├── config/         # 环境变量读取
├── database/       # 数据库初始化
├── models/         # 数据模型
├── controllers/    # HTTP 控制层
├── services/       # 业务逻辑 + 协程池
├── scheduler/      # 定时任务（防重叠）
├── templates/      # Dashboard 页面
│
├── .env.example
├── main.go
└── README.md
```

采用分层架构：

* Controller：处理 HTTP 请求
* Service：核心业务逻辑
* Model：数据结构
* Scheduler：任务调度
* Database：数据连接

---

## ⚙️ 使用方式

### 1️⃣ 克隆项目

```bash
git clone https://github.com/yourname/interface-health-check.git
cd interface-health-check
```

---

### 2️⃣ 创建数据库

在 MySQL 中执行：

```sql
CREATE DATABASE healthcheck DEFAULT CHARSET=utf8mb4;
```

---

### 3️⃣ 配置环境变量

复制示例文件：

```bash
cp .env.example .env
```

编辑 `.env`：

```env
DB_USER=root
DB_PASSWORD=你的密码
DB_HOST=127.0.0.1
DB_PORT=3306
DB_NAME=healthcheck

MAX_WORKER=5
CRON_SPEC=@every 30s
```

---

### 4️⃣ 安装依赖

```bash
go mod tidy
```

---

### 5️⃣ 运行项目

```bash
go run .
```

启动成功后：

```
http://localhost:8080/dashboard
```

---

## 📡 API 示例

### 添加巡检接口

Windows CMD：

```cmd
curl -X POST http://localhost:8080/api/apis ^
 -H "Content-Type: application/json" ^
 -d "{\"name\":\"百度\",\"url\":\"https://www.baidu.com\",\"method\":\"GET\"}"
```


---

### 查询接口列表

```bash
curl http://localhost:8080/api/apis
```

---

### 删除接口

```bash
curl -X DELETE http://localhost:8080/api/apis/1
```

---

### 查询巡检记录

```bash
curl http://localhost:8080/api/checks
```

---

