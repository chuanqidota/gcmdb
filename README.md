# GCMDB

Golang CMDB（配置管理数据库）—— 通用 IT 资产管理系统，支持动态模型定义、实例管理、关系拓扑、全文检索、审计日志，开箱即用。

## 功能特性

- **动态模型**：自定义资产类型（主机、数据库、服务等），支持 7 种字段类型（string/number/bool/date/datetime/json/enum）、字段分组、唯一约束
- **实例管理**：基于模型创建资产实例，JSON 存储 + 索引字段加速查询，支持字段校验
- **模型关系**：定义模型间的关系类型（如「部署在」「依赖」），双向描述
- **字段关联**：通过字段值匹配自动建立实例关系（如 A.host_ip = B.ip），自动同步
- **全文检索**：基于 MySQL JSON_SEARCH 的全文搜索，支持模型过滤和分页
- **结构化搜索**：支持 eq/ne/gt/lt/in/contains/startswith 等操作符，可组合多条件
- **实例拓扑**：AntV G6 力导向图可视化实例上下游关系
- **SQL 查询**：保存并执行自定义只读 SQL 查询
- **审计日志**：记录所有写操作的变更前后数据，支持按资源追踪
- **用户管理**：Session + Token 双认证，管理员/普通用户角色
- **OpenAPI**：对外 RESTful API，支持 Token 认证，便于第三方集成

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.25 + Gin + GORM + MySQL 8 |
| 前端 | Vue 3 + Element Plus + Vite 6 + Pinia + AntV G6 |
| 部署 | Docker + Docker Compose + Nginx |

## 系统架构

```
┌──────────────┐       ┌──────────────┐       ┌──────────────┐
│   Browser    │──────▶│   Nginx:80   │──────▶│  Go API:8080 │──────▶ MySQL:3306
│              │       │  静态文件 +   │       │   Gin REST   │
│              │       │  反向代理     │       │              │
└──────────────┘       └──────────────┘       └──────────────┘
```

- **Nginx**：提供前端 SPA 静态文件，反向代理 `/v1/` 和 `/openapi/` 到后端
- **Go 后端**：Gin 框架提供 REST API，GORM 管理数据库
- **MySQL 8**：数据存储，utf8mb4 字符集

## 数据库表结构

### 表说明

| 表名 | 说明 | 核心字段 |
|------|------|----------|
| `user` | 用户表 | username, password_hash, token, is_active, is_admin |
| `model_group` | 模型分组 | alias, name, description, order |
| `model` | 模型定义 | group_id, alias, name, is_usable, icon, order |
| `model_field_group` | 字段分组 | model_id, name, order |
| `model_field` | 字段定义 | model_id, field_group_id, alias, name, type, is_required, is_indexed, options |
| `model_field_unique` | 唯一约束 | model_id, fields（逗号分隔的字段名） |
| `model_relation_type` | 关系类型 | name, s2t（源→目标描述）, t2s（目标→源描述） |
| `model_relation` | 模型关系 | source_id, target_id, type_id |
| `model_field_relation` | 字段关联 | source_model_id, target_model_id, source_field_id, target_field_id |
| `instance` | 实例数据 | model_id, data（JSON）, indexed_values（索引值） |
| `instance_relation` | 实例关系 | source_model_id, target_model_id, source_instance_id, target_instance_id |
| `search_direct_sql` | SQL 查询 | uuid, name, sql |
| `audit_log` | 审计日志 | resource_type, resource_id, action, method, path, username, before_data, after_data |

### 关系说明

```
user                          # 用户（独立）
audit_log                     # 审计日志（独立）
search_direct_sql             # 保存的 SQL 查询（独立）

model_group                   # 模型分组
└── model                     # 模型（属于分组）
    ├── model_field_group     # 字段分组
    │   └── model_field       # 字段（属于字段分组）
    ├── model_field           # 字段（直接属于模型）
    ├── model_field_unique    # 唯一约束
    └── instance              # 实例（数据以 JSON 存储）

# 模型关系层
model ──model_relation── model      # 模型间关系（由 model_relation_type 定义类型）
model_field ──model_field_relation── model_field  # 字段间关联

# 实例关系层（通过字段关联自动同步）
instance ──instance_relation── instance
```

- **模型定义层**：`model_group` 包含多个 `model`，每个 `model` 有多个 `model_field`（可分组）和 `model_field_unique`
- **模型关系层**：`model` 之间通过 `model_relation` 建立关系（由 `model_relation_type` 定义类型），`model_field` 之间通过 `model_field_relation` 建立字段级关联
- **实例数据层**：`model` 下有多个 `instance`（数据以 JSON 存储），`instance` 之间通过 `instance_relation` 建立关系
- **字段关联自动同步**：当定义了 `model_field_relation` 后，系统自动根据字段值匹配创建/更新 `instance_relation`

## 快速开始

### Docker Compose（推荐）

```bash
# 克隆项目
git clone <repo-url> gcmdb && cd gcmdb

# 启动所有服务（MySQL + 后端 + 前端）
docker compose up -d

# 查看日志
docker compose logs -f
```

启动后访问：
- 前端：http://localhost
- 后端 API：http://localhost:8080
- 默认管理员：`admin` / `admin123`

首次启动会自动执行数据库迁移并创建默认管理员账号。

### 本地开发

**后端：**

```bash
cd backend

# 安装依赖
go mod tidy

# 启动（需先准备好 MySQL）
go run main.go

# 数据库迁移
go run main.go migrate
```

**前端：**

```bash
cd frontend

# 安装依赖
npm install

# 开发服务器（端口 3000，自动代理到后端 8080）
npm run dev

# 构建生产包
npm run build
```

## 配置

后端配置文件位于 `backend/config/config.yaml`：

```yaml
server:
  host: 0.0.0.0
  port: 8080
  auto_relation: true        # 创建实例时自动同步关系
  session_max_age: 604800    # Session 有效期（秒，默认 7 天）
database:
  username: root
  password: "123456"
  host: 127.0.0.1
  port: 3306
  name: gcmdb
cors:
  allowed_origins: []        # 为空则允许所有来源
token_cache_ttl: 5           # Token 缓存 TTL（分钟）
```

### 环境变量

环境变量优先于配置文件，前缀 `GCMDB_`，点号替换为下划线：

| 环境变量 | 说明 | 配置文件字段 | Docker 默认值 |
|----------|------|-------------|---------------|
| `GCMDB_SERVER_HOST` | 监听地址 | `server.host` | `0.0.0.0` |
| `GCMDB_SERVER_PORT` | 监听端口 | `server.port` | `8080` |
| `GCMDB_DATABASE_HOST` | 数据库地址 | `database.host` | `mysql` |
| `GCMDB_DATABASE_PORT` | 数据库端口 | `database.port` | `3306` |
| `GCMDB_DATABASE_USERNAME` | 数据库用户 | `database.username` | `root` |
| `GCMDB_DATABASE_PASSWORD` | 数据库密码 | `database.password` | `root` |
| `GCMDB_DATABASE_NAME` | 数据库名称 | `database.name` | `gcmdb` |
| `GCMDB_CONFIG` | 自定义配置文件路径 | — | — |

本地开发直接使用 `config.yaml`，Docker 部署通过环境变量覆盖，互不干扰。

## Docker 部署说明

### 架构

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  gcmdb-web  │────▶│    gcmdb    │────▶│    mysql     │
│  (nginx)    │     │  (Go/Gin)   │     │  (MySQL 8)  │
│  :80        │     │  :8080      │     │  :3306       │
└─────────────┘     └─────────────┘     └─────────────┘
```

- **gcmdb-web**：Nginx 提供前端静态文件，反向代理 `/v1/` 和 `/openapi/` 到后端
- **gcmdb**：Go 后端 API 服务，启动时自动执行 migrate 建表
- **mysql**：MySQL 8 数据库，utf8mb4 字符集，健康检查确保就绪后后端才启动

### 数据持久化

MySQL 数据通过 Docker volume `mysql_data` 持久化。

```bash
# 查看数据卷
docker volume inspect gcmdb_mysql_data

# 重置所有数据（危险操作）
docker compose down -v
docker compose up -d
```

### Nginx 代理

前端 Nginx 使用 Docker 内置 DNS（`resolver 127.0.0.11`）动态解析后端容器地址，后端容器重启后自动更新 IP，无需手动配置。

### 自定义端口

修改 `docker-compose.yaml` 中的 ports 映射：

```yaml
services:
  gcmdb-web:
    ports:
      - "8080:80"    # 前端改为 8080 端口
  gcmdb:
    ports:
      - "9090:8080"  # 后端改为 9090 端口
  mysql:
    ports:
      - "3306:3306"  # MySQL 改为 3306 端口
```

## Kubernetes 部署

推荐使用 Ingress 路由方案，前后端独立部署：

```
Ingress
  ├── /v1/, /openapi/  →  gcmdb-backend Service  →  后端 Pod
  └── /*               →  gcmdb-frontend Service  →  前端 Pod
```

### 前端 Nginx 配置（K8s 环境）

K8s 中前端 Nginx 只负责静态文件和 SPA 路由，不需要反向代理（由 Ingress 处理）：

```nginx
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

### Ingress 配置示例

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: gcmdb
spec:
  rules:
  - http:
      paths:
      - path: /v1/
        pathType: Prefix
        backend:
          service:
            name: gcmdb-backend
            port:
              number: 8080
      - path: /openapi/
        pathType: Prefix
        backend:
          service:
            name: gcmdb-backend
            port:
              number: 8080
      - path: /
        pathType: Prefix
        backend:
          service:
            name: gcmdb-frontend
            port:
              number: 80
```

## 项目结构

```
gcmdb/
├── backend/
│   ├── app/
│   │   ├── auth/              # 认证（登录、用户管理）
│   │   │   ├── api/           #   HTTP handlers
│   │   │   └── models/        #   User 模型
│   │   ├── audit/             # 审计日志
│   │   │   ├── api/           #   HTTP handlers
│   │   │   └── models/        #   AuditLog 模型
│   │   ├── cmdb/              # 核心 CMDB 逻辑
│   │   │   ├── api/           #   HTTP handlers
│   │   │   ├── models/        #   数据模型定义
│   │   │   ├── params/        #   请求参数结构
│   │   │   └── utils/         #   工具函数（校验、关系同步）
│   │   ├── openapi/           # 对外 OpenAPI 接口
│   │   └── tasks/             # 后台任务（关系同步）
│   ├── cmd/                   # CLI 命令（cobra）
│   ├── config/                # 配置加载（viper）
│   ├── pkg/                   # 公共包
│   │   ├── database/          #   数据库初始化
│   │   ├── middleware/        #   中间件（认证、CORS、审计）
│   │   ├── logger/            #   日志
│   │   └── utils/             #   工具函数
│   ├── router/                # 路由定义
│   ├── Dockerfile
│   └── main.go
├── frontend/
│   ├── src/
│   │   ├── api/               # API 请求封装（17 个模块）
│   │   ├── layout/            # 布局组件（侧边栏、头部）
│   │   ├── router/            # 路由配置
│   │   ├── stores/            # Pinia 状态管理
│   │   ├── utils/             # 公共工具函数
│   │   └── views/             # 页面组件
│   │       ├── home/          #   首页仪表盘
│   │       ├── search/        #   综合检索
│   │       ├── instance/      #   实例管理
│   │       ├── model-manage/  #   模型管理
│   │       ├── model-topology/#   模型拓扑图
│   │       ├── search-direct-sql/ # SQL 查询
│   │       ├── audit/         #   审计日志
│   │       └── user-manage/   #   用户管理
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yaml
└── scripts/                   # 测试脚本（.gitignore）
```

## API 概览

### 认证接口 `/v1/auth`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/login` | 登录 |
| POST | `/logout` | 退出 |
| GET | `/me` | 当前用户信息 |
| POST | `/change-password` | 修改密码 |
| POST | `/reset-password` | 重置密码（管理员） |

### CMDB 接口 `/v1/cmdb`（需 Session 登录）

| 模块 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 模型分组 | GET/POST | `/models-group` | 查询/创建 |
| 模型 | GET/POST/PUT/DELETE | `/models` | CRUD |
| 模型字段分组 | GET/POST/PUT/DELETE | `/models-field-group` | CRUD |
| 模型字段 | GET/POST/PUT/DELETE | `/models-field` | CRUD |
| 模型关系类型 | GET/POST/PUT/DELETE | `/models-relation-type` | CRUD |
| 模型关系 | GET/POST/DELETE | `/models-relation` | 管理模型间关系 |
| 字段关联 | GET/POST/DELETE | `/models-field-relation` | 管理字段级关联 |
| 唯一约束 | GET/POST/DELETE | `/models-field-unique` | 管理唯一约束 |
| 实例 | GET/POST/PUT/DELETE | `/instance` | CRUD |
| 实例关系 | GET/POST/DELETE | `/instance-relation` | 管理实例间关系 |
| SQL 查询 | GET/POST/PUT/DELETE | `/search-direct-sql` | 保存/执行自定义 SQL |
| 审计日志 | GET | `/audit-log` | 查询操作记录 |
| 统计 | GET | `/stats` | 首页统计数据 |
| 任务 | POST | `/tasks/sync-relations` | 手动触发关系同步 |

### OpenAPI `/openapi`（Session 或 Token 认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/model/:range` | 模型查询（all/group/single/relation/relation-type） |
| GET | `/instance/:action` | 实例查询（fulltext/search/detail/topology/source/target/direct） |
| POST | `/instance/:action` | 实例写入（create/update/delete） |
| POST | `/instance-relation/:action` | 实例关系（create/delete） |

### 认证方式

| 方式 | 适用场景 | 说明 |
|------|---------|------|
| Session | 前端 Web | 登录后自动管理 Cookie（`gcmdb_session`） |
| Token | API 集成 | 每个用户一个 UUID Token，通过 `Authorization: Bearer <token>` 传递 |

## License

MIT
