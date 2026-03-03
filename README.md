# CyberMatchmaker 塞博红娘：根据算命来匹配心仪对象（后端部分）

一个将“算命/命理分析（LLM）”与“向量相似度检索（pgvector）”结合的后端服务。基于 Gin + GORM + PostgreSQL（pgvector）、Redis（会话校验）与 RabbitMQ（异步任务）。

## 功能
- 用户注册/登录，JWT 鉴权 + Redis 会话校验。
- 算命流程：异步 LLM 任务、结果落库、前端轮询查询。
- 匹配流程：加入广场（生成画像与向量）、pgvector 余弦相似度检索、退出广场（软删除）。
- 统一配置与 Prompt：`config/config.yaml` 与 `config/prompt.yaml`。

## 技术栈
- Go, Gin, GORM
- PostgreSQL + pgvector
- Redis
- RabbitMQ
- LLM（langchaingo / OpenAI 兼容接口）

## 项目结构（后端）
```
config/        配置与 Prompt 模板
controller/    HTTP 处理器
service/       业务逻辑
mapper/        数据库访问与向量检索
model/         GORM 模型
middleware/    JWT 鉴权 + LLM 封装
mq/            RabbitMQ 生产者/消费者
pkg/           基础设施、jwt、response、utils
route/         Gin 路由
```

## 核心流程（ASCII）

### 算命流程（异步）
```
+-----------+    +------------+    +-------------+    +-----------+
|  Client   | -> |  Gin API   | -> |  Service    | -> |  Postgres |
+-----------+    +------------+    +-------------+    +-----------+
                      |                   |
                      |                   v
                      |            +--------------+
                      |            |  RabbitMQ    |
                      |            +--------------+
                      |                   |
                      |                   v
                      |            +--------------+
                      |            |  Consumer    |
                      |            +--------------+
                      |                   |
                      |                   v
                      |            +--------------+
                      |            |  LLM API     |
                      |            +--------------+
                      |                   |
                      |                   v
                      |            +--------------+
                      |            |  Postgres    |
                      |            +--------------+
                      v
+-----------+    +------------+
|  Client   | <- |  Poll API  |
+-----------+    +------------+
```

### 匹配流程（加入/搜索/退出）
```
+-----------+    +------------+    +-------------+
|  Client   | -> |  Gin API   | -> |  Service    |
+-----------+    +------------+    +-------------+
                      |                   |
                      | (join)            v
                      |            +--------------+
                      |            |  RabbitMQ    |
                      |            +--------------+
                      |                   |
                      |                   v
                      |            +--------------+
                      |            |  Consumer    |
                      |            +--------------+
                      |                   |
                      |                   v
                      |            +--------------+
                      |            |  LLM Embed   |
                      |            +--------------+
                      |                   |
                      |                   v
                      |            +--------------+
                      |            |  Postgres    |
                      |            +--------------+
                      |
                      | (search) ---------> pgvector 相似度检索
                      |
                      | (leave) ----------> 软删除
```

## 接口文档（以代码为准）

### 全局
- Base URL: `/api/v1`
- Auth: `Authorization: Bearer <JWT>`（由 `middleware/auth.go` 校验）
- 统一返回：`{ "code": 0, "msg": "success", "data": {} }`

### 1) 用户
#### POST `/user/register`
Request:
```json
{ "username": "user123", "password": "123456" }
```
Response: `null`

#### POST `/user/login`
Request:
```json
{ "username": "user123", "password": "123456" }
```
Response:
```json
{ "token": "<jwt>" }
```

### 2) 算命
#### POST `/fortune/submit`（鉴权）
Request:
```json
{
  "real_name": "张三",
  "gender": "男",
  "birth_date": "1998-05-20",
  "birth_time": "23:00",
  "current_city": "杭州"
}
```
Response:
```json
{ "record_id": 1024 }
```

#### GET `/fortune/result`（鉴权）
说明：返回该用户最新一条已完成的算命记录。
Response:
```json
{
  "fortune_result": {
    "status": "completed",
    "bazi": "...",
    "five_elements": "...",
    "zodiac_sign": "...",
    "best_city": "...",
    "recent_fortune": "...",
    "description": "..."
  }
}
```

### 3) 匹配
#### POST `/match/join`（鉴权）
Request:
```json
{ "wechat_id": "wx_zhangsan_888" }
```
Response: 成功提示

#### GET `/match/search`（鉴权）
Request: 无
Response:
```json
{
  "message": "match success",
  "data": {
    "real_name": "李四",
    "wechat_id": "wx_lisi_999",
    "gender": "女",
    "birth_date": "1999-01-01",
    "current_city": "杭州",
    "bazi": "...",
    "five_elements": "...",
    "similarity": 0.92
  }
}
```

#### POST `/match/leave`（鉴权）
Request: 无
Response: 成功提示

## 配置说明
- `config/config.yaml`：DB/Redis/RabbitMQ/LLM/JWT 配置。
- `config/prompt.yaml`：算命任务 Prompt 模板。

## 入口说明
- `main.go`：加载配置、日志、基础设施并启动 HTTP 服务。
- `pkg/infra/infra.go`：初始化 DB/Redis/RabbitMQ/LLM，并启动 MQ 消费者。

## 快速启动（示例）
```bash
go run ./...
```

## 参考
- `接口文档.md` 可能与当前代码有差异；本 README 以代码为准。
