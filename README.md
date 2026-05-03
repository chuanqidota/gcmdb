# GCMDB

Golang CMDB（配置管理数据库）—— 通用资产管理系统，支持自定义模型、字段、关系，灵活管理各类 IT 资产。

## 功能特性

- **自定义模型**：动态创建模型（如主机、数据库、服务），定义字段类型（string/number/bool/date/datetime/json/enum）、分组、唯一约束
- **实例管理**：基于模型创建实例，支持字段校验、索引加速查询、批量操作
- **模型关系**：定义模型间的关系类型（如「部署在」「依赖」），自动同步实例级关系
- **字段关联**：通过字段值匹配自动建立实例关系（如 A.host_ip 关联 B.ip）
- **全文检索**：基于 MySQL JSON_SEARCH 的全文搜索，支持模型过滤和分页
- **结构化搜索**：支持 eq/ne/gt/lt/in/contains/startswith 等操作符，可组合条件
- **实例拓扑**：可视化实例上下游关系，支持力导向图展示
- **SQL 查询**：保存并执行自定义 SQL 查询（只读 SELECT）
- **审计日志**：记录所有写操作，支持按资源类型和 ID 追踪变更历史
- **用户管理**：登录认证、Token API、角色权限（管理员/普通用户）
- **OpenAPI**：对外暴露 RESTful API，支持 Token 认证

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.21 + Gin + GORM + MySQL |
| 前端 | Vue 3 + Element Plus + Vite + Pinia + AntV G6 |
| 部署 | Docker + Docker Compose + Nginx |

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

# 开发服务器
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
  session_max_age: 604800    # Session 有效期（秒）
database:
  username: root
  password: "123456"
  host: 127.0.0.1
  port: 3306
  name: gcmdb
cors:
  allowed_origins: []
token_cache_ttl: 5           # Token 缓存 TTL（分钟）
```

Docker 部署时，环境变量会覆盖配置文件（前缀 `GCMDB_`）：

| 环境变量 | 对应配置 |
|----------|----------|
| `GCMDB_DATABASE_HOST` | `database.host` |
| `GCMDB_DATABASE_PORT` | `database.port` |
| `GCMDB_DATABASE_USERNAME` | `database.username` |
| `GCMDB_DATABASE_PASSWORD` | `database.password` |
| `GCMDB_DATABASE_NAME` | `database.name` |

## 项目结构

```
gcmdb/
├── backend/
│   ├── app/
│   │   ├── auth/          # 认证（登录、用户管理）
│   │   ├── audit/         # 审计日志
│   │   ├── cmdb/          # 核心 CMDB 逻辑
│   │   │   ├── api/       # HTTP handlers
│   │   │   ├── models/    # 数据模型定义
│   │   │   ├── params/    # 请求参数结构
│   │   │   └── utils/     # 工具函数（校验、关系同步）
│   │   ├── openapi/       # 对外 OpenAPI 接口
│   │   └── tasks/         # 后台任务（关系同步）
│   ├── cmd/               # CLI 命令（cobra）
│   ├── config/            # 配置加载（viper）
│   ├── pkg/               # 公共包（database、middleware、logger）
│   ├── router/            # 路由定义
│   ├── Dockerfile
│   └── main.go
├── frontend/
│   ├── src/
│   │   ├── api/           # API 请求封装
│   │   ├── layout/        # 布局组件（侧边栏、头部）
│   │   ├── router/        # 路由配置
│   │   ├── stores/        # Pinia 状态管理
│   │   ├── utils/         # 公共工具函数
│   │   └── views/         # 页面组件
│   │       ├── home/      # 首页
│   │       ├── search/    # 综合检索（全文/实例/模型/接口调试）
│   │       ├── instance/  # 实例管理
│   │       ├── model-manage/   # 模型管理
│   │       ├── model-topology/ # 模型拓扑
│   │       ├── audit/     # 审计日志
│   │       └── user-manage/    # 用户管理
│   ├── Dockerfile
│   └── nginx.conf
├── docker-compose.yaml
└── scripts/               # 测试脚本
```

## API 概览

### 认证接口 `/v1/auth`

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/login` | 登录 |
| POST | `/logout` | 退出 |
| GET | `/me` | 当前用户信息 |
| POST | `/change-password` | 修改密码 |

### CMDB 接口 `/v1/cmdb`（需登录）

| 模块 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 模型分组 | GET/POST | `/models-group` | 查询/创建 |
| 模型 | GET/POST/PUT/DELETE | `/models` | CRUD |
| 模型字段 | GET/POST/PUT/DELETE | `/models-field` | CRUD |
| 模型关系 | GET/POST/DELETE | `/models-relation` | 管理模型间关系 |
| 实例 | GET/POST/PUT/DELETE | `/instance` | CRUD |
| 实例关系 | GET/POST/DELETE | `/instance-relation` | 管理实例间关系 |
| SQL 查询 | GET/POST/PUT/DELETE | `/search-direct-sql` | 保存/执行自定义 SQL |
| 审计日志 | GET | `/audit-log` | 查询操作记录 |

### OpenAPI `/openapi`（Token 认证）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/model/:range` | 模型查询（all/group/single/relation/relation-type） |
| GET | `/instance/:action` | 实例查询（fulltext/search/detail/topology/source/target/direct） |
| POST | `/instance/:action` | 实例写入（create/update/delete） |
| POST | `/instance-relation/:action` | 实例关系（create/delete） |

## License

MIT
