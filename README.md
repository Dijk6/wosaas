# wosaas

`wosaas` 是一个基于 **Gin + GORM** 的可运行 Web 服务框架，提供配置、日志、路由、中间件、统一响应、数据访问、事件总线、定时任务和异步导出等基础能力。

当前范围：只提供框架能力与最小可运行示例，业务模块可以按目录分层继续扩展。

## 快速开始

```bash
cd wosaas
go mod tidy
go run . -conf ./conf/app.toml
```

启动后：

```bash
curl http://127.0.0.1:8600/health
curl http://127.0.0.1:8600/ping
curl http://127.0.0.1:8600/driver/profile
curl -X POST http://127.0.0.1:8600/openapi/partner/price/estimate \
  -H 'Content-Type: application/json' \
  -d '{"start_lng":116.397,"start_lat":39.908,"end_lng":116.407,"end_lat":39.918}'
```

管理端口：`http://127.0.0.1:8081/health`、`http://127.0.0.1:8081/metrics`。

本地没有 MySQL/Redis 时服务仍可启动（初始化阶段只打印警告），便于先验证框架；生产环境建议改为 fail-fast。

## 目录结构

```text
wosaas/
├── main.go                     程序入口，加载配置、启动 bootstrap 与 server
├── config/                     TOML 配置加载、端口占位符解析
├── conf/                       本地/开发配置
├── bootstrap/                  日志、MySQL、Redis、缓存、Trace、cron、eventbus 初始化
├── servers/
│   ├── start.go                启动与优雅退出
│   └── httpserver/
│       ├── server.go           HTTP Server（主端口 + 管理端口）
│       ├── router.go           路由总入口
│       ├── middleware/         Recovery/Trace/AccessLog/CORS/Timeout/Auth
│       ├── routers/            openapi、driver 两个路由示例
│       └── controller/         控制器
├── bizdata/
│   ├── model/                  GORM 模型基类
│   ├── request/                请求 DTO 示例
│   └── response/               响应 DTO 示例
├── library/
│   ├── constant/               常量、上下文 key
│   ├── errno/                  通用错误码
│   ├── eventbus/               进程内事件总线
│   ├── log/                    请求级扩展日志
│   ├── resource/               全局单例资源
│   ├── sessionlib/             登录态签发/校验
│   ├── trace/                  TraceID 生成
│   ├── types/                  UserInfo、统一响应
│   └── utils/paging/           分页工具
└── model/
    ├── dao/db/                 GORM 数据源入口
    ├── service/                业务服务层示例
    └── cronjob/                定时任务注册
```

## 框架特性

- Gin 路由与中间件
- TOML 配置与端口占位符
- GORM 多数据源
- Redis 与本地缓存
- 统一错误码与 JSON 响应
- 进程内事件总线
- 定时任务
- 异步导出

## 新增一个业务模块

以 `openapi/price/estimate` 为例：

1. 在 `bizdata/request`、`bizdata/response` 定义 DTO。
2. 在 `model/service` 编写业务逻辑。
3. 需要访问数据库时，通过 `model/dao/db.DB(name)` 获取数据源并编写 DAO。
4. 在 `servers/httpserver/controller` 写 handler，使用 `response.JSON*` 返回。
5. 在 `servers/httpserver/routers` 注册路由。
6. 需要异步解耦时，通过 `resource.GlobalEventBus.Publish` 发布事件。

## 导出模块

导出采用 **任务表 + 定时 worker + 文件存储** 的异步架构：

```text
前端提交导出任务
  -> export_task 表 status=1
定时 worker 扫描任务
  -> 条件更新抢占 status=2
  -> 执行 SELECT SQL
  -> excelize 流式写入 xlsx
  -> 保存到本地文件存储
  -> 更新 status=3 + file_url
前端轮询任务列表，成功后下载文件
```

### 目录

```text
bizdata/
├── model/export_task.go                 导出任务 GORM 模型
├── request/manage/export.go             创建/列表请求
└── response/manage/export.go            列表响应

library/
├── constant/export.go                   任务状态常量
└── filestore/local.go                   文件存储抽象（当前为本地磁盘）

model/
├── dao/db/saas/export_task.go           导出任务 DAO
├── data/exporttask/exporttask.go        数据层
├── service/manage/exporttask/           创建、列表、SQL 校验、状态更新
└── cronjob/export.go                    导出 worker
servers/httpserver/
├── controller/manage/exporttask/        新增、列表、下载、手动触发
└── routers/manage.go                    路由注册
```

### 接口

```text
POST /manage/export/task/add
GET  /manage/export/task/list
POST /manage/export/task/run      手动触发 worker（调试/补偿用）
GET  /manage/export/download/:date/:file
```

创建任务示例：

```bash
curl -X POST http://127.0.0.1:8600/manage/export/task/add \
  -H 'Content-Type: application/json' \
  -d '{
    "src_page": "example_list",
    "database_name": "saas",
    "columns_name": "ID,名称,创建时间",
    "sql_content": "SELECT id, name, created_at FROM example_table",
    "sql_count": 1000
  }'
```

调试模式不会启动 cron，可以手动触发：

```bash
curl -X POST http://127.0.0.1:8600/manage/export/task/run
```

### 建表

```bash
mysql -uroot -p wosaas < data/sql/export_task.sql
```

### 安全限制

- 只允许单条 `SELECT` / `WITH` SQL
- 禁止分号、顶层 `LIMIT/OFFSET`、`INTO OUTFILE`、`FOR UPDATE`
- 生产环境建议使用只读数据库账号
- 生产环境建议改为服务端生成 SQL，而不是直接接收外部 SQL
- 后续可把 `library/filestore` 替换为 S3/OSS 等对象存储，并接入短链与过期清理

## 配置说明

- `app.toml`：应用名、运行模式、端口。`From = "env"` 时优先读取 `Value` 指定的环境变量。
- `server.toml`：`{port.Main}`、`{port.Admin}` 会被解析为实际端口。
- `mysql.toml`：支持多数据源，`name` 即 `resource.DB(name)` 的名称。
- `redis.toml`：Redis 连接参数。
- `global.toml`：业务全局配置，由 `resource.GlobalConfig` 解析。

常用环境变量：

| 变量 | 说明 |
| --- | --- |
| `LISTEN_PORT` | 主 HTTP 端口，默认 8600 |
| `LISTEN_PORT_1` | 管理端口，默认 8081 |
| `LOG_LEVEL` | 日志级别，debug/info/warn/error |
| `LOG_FILE` | 日志文件路径，设置后同时输出到文件并滚动 |
| `SESSION_SECRET` | 登录 Token 签名密钥 |
| `SESSION_TTL` | Token 有效期，如 `24h` |
| `EXPORT_DIR` | 导出文件本地目录，默认 `./data/export` |
| `EXPORT_URL_PREFIX` | 导出文件下载前缀，默认 `/manage/export/download` |

## 后续建议

1. 接入 Prometheus：在 `/metrics` 输出 `promhttp`，并为 Gin 增加指标中间件。
2. 接入 OpenTelemetry / Jaeger：替换 `library/trace`，在 HTTP/GORM/Redis 客户端统一注入 span。
3. 用 Redis 实现 Session 与分布式锁。
4. 引入 migration 工具（golang-migrate / atlas）管理表结构。
5. 补齐单元测试与基于 `httptest` 的路由测试。
